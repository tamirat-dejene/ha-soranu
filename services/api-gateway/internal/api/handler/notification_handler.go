package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/api/dto"
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/client"
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/errs"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"go.uber.org/zap"
)

type NotificationHandler struct {
	client *client.NotificationServiceClient
}

func NewNotificationHandler(client *client.NotificationServiceClient) *NotificationHandler {
	return &NotificationHandler{
		client: client,
	}
}

func (h *NotificationHandler) GetUserNotifications(c *gin.Context) {
	// TODO: Extract userID from JWT/session
	userID := c.Param("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse("user_id is required"))
		return
	}

	notifications, err := h.client.GetNotifications(c.Request.Context(), userID, "USER")
	if err != nil {
		logger.Error("failed to get user notifications", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	type NotificationResponse struct {
		ID            string `json:"id"`
		RecipientID   string `json:"recipient_id"`
		RecipientType string `json:"recipient_type"`
		OrderID       string `json:"order_id"`
		Title         string `json:"title"`
		Message       string `json:"message"`
		IsRead        bool   `json:"is_read"`
		Type          string `json:"type"`
		CreatedAt     string `json:"created_at"`
	}

	var response []NotificationResponse
	for _, n := range notifications {
		response = append(response, NotificationResponse{
			ID:            n.Id,
			RecipientID:   n.RecipientId,
			RecipientType: n.RecipientType,
			OrderID:       n.OrderId,
			Title:         n.Title,
			Message:       n.Message,
			IsRead:        n.IsRead,
			Type:          n.Type,
			CreatedAt:     n.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"notifications": response})
}

func (h *NotificationHandler) GetRestaurantNotifications(c *gin.Context) {
	restaurantID := c.Param("restaurant_id")
	if restaurantID == "" {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse("restaurant_id is required"))
		return
	}

	notifications, err := h.client.GetNotifications(c.Request.Context(), restaurantID, "RESTAURANT")
	if err != nil {
		logger.Error("failed to get restaurant notifications", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	type NotificationResponse struct {
		ID            string `json:"id"`
		RecipientID   string `json:"recipient_id"`
		RecipientType string `json:"recipient_type"`
		OrderID       string `json:"order_id"`
		Title         string `json:"title"`
		Message       string `json:"message"`
		IsRead        bool   `json:"is_read"`
		Type          string `json:"type"`
		CreatedAt     string `json:"created_at"`
	}

	var response []NotificationResponse
	for _, n := range notifications {
		response = append(response, NotificationResponse{
			ID:            n.Id,
			RecipientID:   n.RecipientId,
			RecipientType: n.RecipientType,
			OrderID:       n.OrderId,
			Title:         n.Title,
			Message:       n.Message,
			IsRead:        n.IsRead,
			Type:          n.Type,
			CreatedAt:     n.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, gin.H{"notifications": response})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	notificationID := c.Param("notification_id")
	if notificationID == "" {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse("notification_id is required"))
		return
	}

	success, err := h.client.MarkAsRead(c.Request.Context(), notificationID)
	if err != nil {
		logger.Error("failed to mark notification as read", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": success})
}
