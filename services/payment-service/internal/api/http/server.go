package apihttp

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/webhook"
	"github.com/tamirat-dejene/ha-soranu/services/payment-service/internal/usecase"
)

type Server struct {
	uc               usecase.Service
	stripeWebhookKey string
}

func NewServer(uc usecase.Service, webhookSecret string) *Server {
	return &Server{uc: uc, stripeWebhookKey: webhookSecret}
}

func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/payments/intent", s.handleCreateIntent)
	mux.HandleFunc("/payments/", s.handleGetPayment)
	mux.HandleFunc("/payments/webhook", s.handleWebhook)
	return mux
}

type createIntentRequest struct {
	OrderID  string `json:"order_id"`
	Amount   int64  `json:"amount"` // cents
	Currency string `json:"currency"`
}

type createIntentResponse struct {
	PaymentID    string `json:"payment_id"`
	ClientSecret string `json:"client_secret"`
}

func (s *Server) handleCreateIntent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req createIntentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	payment, err := s.uc.CreatePaymentIntent(ctx, req.OrderID, req.Amount, strings.ToLower(req.Currency))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_ = json.NewEncoder(w).Encode(createIntentResponse{PaymentID: payment.ID, ClientSecret: payment.ClientSecret})
}

func (s *Server) handleGetPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// Path is /payments/{id}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 || parts[0] != "payments" {
		http.NotFound(w, r)
		return
	}
	id := parts[1]
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	payment, err := s.uc.GetPayment(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	_ = json.NewEncoder(w).Encode(payment)
}

func (s *Server) handleWebhook(w http.ResponseWriter, r *http.Request) {
	const maxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read error", http.StatusServiceUnavailable)
		return
	}

	sig := r.Header.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, sig, s.stripeWebhookKey)
	if err != nil {
		http.Error(w, "signature verification failed", http.StatusBadRequest)
		return
	}

	switch event.Type {
	case "payment_intent.succeeded":
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err == nil {
			_ = s.uc.HandleSucceeded(r.Context(), &pi)
		}
	case "payment_intent.payment_failed":
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err == nil {
			_ = s.uc.HandleFailed(r.Context(), &pi)
		}
	default:
		// ignore others
	}
	w.WriteHeader(http.StatusOK)
}

// Optional helper for computing HMAC if needed elsewhere
func hmac256(secret, payload string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}
