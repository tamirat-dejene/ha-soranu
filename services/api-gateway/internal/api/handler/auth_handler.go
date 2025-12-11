package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/api/dto"
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/client"
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/errs"
	"github.com/tamirat-dejene/ha-soranu/shared/pkg/logger"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/authpb"
	"go.uber.org/zap"
)

type AuthHandler struct {
	client *client.AuthClient
}

func NewAuthHandler(client *client.AuthClient) *AuthHandler {
	return &AuthHandler{
		client: client,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	logger.Info("Register request received")
	var req authpb.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.MsgInvalidRequest))
		return
	}

	res, err := h.client.Client.Register(c.Request.Context(), &req)
	if err != nil {
		logger.Error("Failed to register user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	c.JSON(http.StatusOK, dto.RegisterResponseFromProto(res))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req authpb.EPLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.MsgInvalidRequest))
		return
	}

	res, err := h.client.Client.LoginWithEmailAndPassword(c.Request.Context(), &req)
	if err != nil {
		logger.Error("Failed to login user", zap.Error(err))
		c.JSON(http.StatusUnauthorized, dto.ErrorResponseFromGRPCError(err))
		return
	}

	c.JSON(http.StatusOK, dto.EPLoginResponseFromProto(res))
}

func (h *AuthHandler) LoginWithGoogle(c *gin.Context) {
	var req authpb.GLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.MsgInvalidRequest))
		return
	}

	res, err := h.client.Client.LoginWithGoogle(c.Request.Context(), &req)
	if err != nil {
		logger.Error("Failed to login with google", zap.Error(err))
		c.JSON(http.StatusUnauthorized, errs.NewErrorResponse(errs.MsgInvalidCredentials))
		return
	}

	c.JSON(http.StatusOK, dto.GLoginResponseFromProto(res))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req authpb.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.MsgInvalidRequest))
		return
	}

	res, err := h.client.Client.Logout(c.Request.Context(), &req)
	if err != nil {
		logger.Error("Failed to logout user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, errs.NewErrorResponse(errs.MsgInternalError))
		return
	}

	c.JSON(http.StatusOK, dto.LogoutResponseFromProto(res))
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req authpb.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.MsgInvalidRequest))
		return
	}

	res, err := h.client.Client.Refresh(c.Request.Context(), &req)
	if err != nil {
		logger.Error("Failed to refresh token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, errs.NewErrorResponse(errs.MsgInvalidCredentials))
		return
	}

	c.JSON(http.StatusOK, dto.EPLoginResponseFromProto(&authpb.EPLoginResponse{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}))
}
