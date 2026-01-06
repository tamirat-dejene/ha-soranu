package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/client"
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/errs"
)

type PaymentHandler struct {
	client *client.PaymentClient
}

func NewPaymentHandler(c *client.PaymentClient) *PaymentHandler {
	return &PaymentHandler{client: c}
}

type createPaymentDTO struct {
	OrderID  string `json:"order_id"` // required
	Amount   int64  `json:"amount"`   // required, cents
	Currency string `json:"currency"` // required, e.g., "usd"
}

func (h *PaymentHandler) CreatePaymentIntent(c *gin.Context) {
	var req createPaymentDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.MsgInvalidRequest))
		return
	}
	if req.OrderID == "" || req.Amount <= 0 || req.Currency == "" {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse("order_id, amount and currency are required"))
		return
	}
	res, err := h.client.CreateIntent(client.CreateIntentRequest{OrderID: req.OrderID, Amount: req.Amount, Currency: req.Currency})
	if err != nil {
		c.JSON(http.StatusBadGateway, errs.NewErrorResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"payment_id": res.PaymentID, "client_secret": res.ClientSecret})
}
