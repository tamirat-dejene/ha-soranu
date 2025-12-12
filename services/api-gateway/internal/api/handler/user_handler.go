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

type UserHandler struct {
	client *client.UAServiceClient
}

func NewUserHandler(client *client.UAServiceClient) *UserHandler {
	return &UserHandler{client: client}
}

/**
	GetUser(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*GetUserResponse, error)
    GetPhoneNumber(ctx context.Context, in *GetPhoneNumberRequest, opts ...grpc.CallOption) (*GetPhoneNumberResponse, error)
    AddPhoneNumber(ctx context.Context, in *AddPhoneNumberRequest, opts ...grpc.CallOption) (*MessageResponse, error)
    UpdatePhoneNumber(ctx context.Context, in *UpdatePhoneNumberRequest, opts ...grpc.CallOption) (*MessageResponse, error)
    RemovePhoneNumber(ctx context.Context, in *RemovePhoneNumberRequest, opts ...grpc.CallOption) (*MessageResponse, error)
    GetAddresses(ctx context.Context, in *GetAddressesRequest, opts ...grpc.CallOption) (*GetAddressesResponse, error)
    AddAddress(ctx context.Context, in *AddAddressRequest, opts ...grpc.CallOption) (*AddAddressResponse, error)
    RemoveAddress(ctx context.Context, in *RemoveAddressRequest, opts ...grpc.CallOption) (*MessageResponse, error)
*/


func (h *UserHandler) GetUser(c *gin.Context) {
	userId := c.Query("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse("user_id is required"))
		return
	}

	req := &dto.GetUserRequestDTO{UserId: userId}

	resp, err := h.client.UserClient.GetUser(c.Request.Context(), req.ToProto())
	if err != nil {
		logger.Error("Failed to get user", zap.String("user_id", req.UserId), zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	logger.Info("Successfully retrieved user", zap.Any("User", resp))
	c.JSON(200, dto.GetUserResponseFromProto(resp))
}

func (h *UserHandler) GetPhoneNumber(c *gin.Context) {
	userId := c.Query("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse("user_id is required"))
		return
	}

	req := &dto.GetPhoneNumberRequestDTO{UserId: userId}

	resp, err := h.client.UserClient.GetPhoneNumber(c.Request.Context(), req.ToProto())
	if err != nil {
		logger.Error("Failed to get phone number", zap.String("user_id", req.UserId), zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	c.JSON(200, dto.GetPhoneNumberResponseFromProto(resp))
}

func (h *UserHandler) AddPhoneNumber(c *gin.Context) {
	var req dto.AddPhoneNumberRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.MsgInvalidRequest))
		return
	}

	resp, err := h.client.UserClient.AddPhoneNumber(c.Request.Context(), req.ToProto())
	if err != nil {
		logger.Error("AddPhoneNumber failed", zap.String("user_id", req.UserId), zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	c.JSON(200, gin.H{"message": resp.Message})
}

func (h *UserHandler) UpdatePhoneNumber(c *gin.Context) {
	var req dto.UpdatePhoneNumberRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.MsgInvalidRequest))
		return
	}

	resp, err := h.client.UserClient.UpdatePhoneNumber(c.Request.Context(), req.ToProto())
	if err != nil {
		logger.Error("UpdatePhoneNumber failed", zap.String("user_id", req.UserId), zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	c.JSON(200, gin.H{"message": resp.Message})
}

func (h *UserHandler) RemovePhoneNumber(c *gin.Context) {
	var req dto.RemovePhoneNumberRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.MsgInvalidRequest))
		return
	}

	resp, err := h.client.UserClient.RemovePhoneNumber(c.Request.Context(), req.ToProto())
	if err != nil {
		logger.Error("RemovePhoneNumber failed", zap.String("user_id", req.UserId), zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	c.JSON(200, gin.H{"message": resp.Message})
}

func (h *UserHandler) GetAddresses(c *gin.Context) {
	userId := c.Query("user_id")
	if userId == "" {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse("user_id is required"))
		return
	}

	req := &dto.GetAddressesRequestDTO{UserId: userId}

	resp, err := h.client.UserClient.GetAddresses(c.Request.Context(), req.ToProto())
	if err != nil {
		logger.Error("GetAddresses failed", zap.String("user_id", userId), zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	c.JSON(200, dto.GetAddressesResponseFromProto(resp))
}

func (h *UserHandler) AddAddress(c *gin.Context) {
	var req dto.AddAddressRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.MsgInvalidRequest))
		return
	}

	resp, err := h.client.UserClient.AddAddress(c.Request.Context(), req.ToProto())
	if err != nil {
		logger.Error("AddAddress failed", zap.String("user_id", req.UserId), zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	c.JSON(200, dto.AddAddressResponseFromProto(resp))
}

func (h *UserHandler) RemoveAddress(c *gin.Context) {
	var req dto.RemoveAddressRequestDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.MsgInvalidRequest))
		return
	}

	resp, err := h.client.UserClient.RemoveAddress(c.Request.Context(), req.ToProto())
	if err != nil {
		logger.Error("RemoveAddress failed", zap.String("address_id", req.AddressId), zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	c.JSON(200, gin.H{"message": resp.Message})
}
