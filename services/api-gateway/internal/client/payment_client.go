package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type PaymentClient struct {
	baseURL string
	http    *http.Client
}

func NewPaymentClient(host, port string) *PaymentClient {
	return &PaymentClient{
		baseURL: fmt.Sprintf("http://%s:%s", host, port),
		http:    &http.Client{},
	}
}

type CreateIntentRequest struct {
	OrderID  string `json:"order_id"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type CreateIntentResponse struct {
	PaymentID    string `json:"payment_id"`
	ClientSecret string `json:"client_secret"`
}

func (c *PaymentClient) CreateIntent(req CreateIntentRequest) (*CreateIntentResponse, error) {
	b, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest(http.MethodPost, c.baseURL+"/payments/intent", bytes.NewReader(b))
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("payment-service error: %s", resp.Status)
	}
	var out CreateIntentResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &out, nil
}
