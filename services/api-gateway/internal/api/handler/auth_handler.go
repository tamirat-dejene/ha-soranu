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

type AuthHandler struct {
	client *client.UAServiceClient
}

func NewAuthHandler(client *client.UAServiceClient) *AuthHandler {
	return &AuthHandler{
		client: client,
	}
}

/**
	Register(ctx context.Context, in *UserRegisterRequest, opts ...grpc.CallOption) (*UserRegisterResponse, error)
    LoginWithEmailAndPassword(ctx context.Context, in *EPLoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
    LoginWithGoogle(ctx context.Context, in *GLoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
    Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*userpb.MessageResponse, error)
    Refresh(ctx context.Context, in *RefreshRequest, opts ...grpc.CallOption) (*RefreshResponse, error)
*/

func (h *AuthHandler) Register(c *gin.Context) {
	logger.Info("Register request received")
	var req *dto.UserRegisterRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.MsgInvalidRequest))
		return
	}

	resp, err := h.client.AuthClient.Register(c.Request.Context(), req.ToProto())
	if err != nil {
		logger.Error("Failed to register user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	c.JSON(http.StatusOK, dto.UserRegisterResponseFromProto(resp))
}

func (h *AuthHandler) LoginWithEmailAndPassword(c *gin.Context) {
	logger.Info("Email-password login request received")
	var req *dto.EPLoginRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.MsgInvalidRequest))
		return
	}
	resp, err := h.client.AuthClient.LoginWithEmailAndPassword(c.Request.Context(), req.ToProto())
	if err != nil {
		logger.Error("Failed to login with email and password", zap.String("email", req.Email), zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponseFromProto(resp))
}

func (h *AuthHandler) LoginWithGoogle(c *gin.Context) {
	logger.Info("Google login request received")
	var req *dto.GoogleLoginRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.MsgInvalidRequest))
		return
	}

	resp, err := h.client.AuthClient.LoginWithGoogle(c.Request.Context(), req.ToProto())
	if err != nil {
		logger.Error("Failed to login with Google", zap.String("id_token", req.IdToken), zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponseFromProto(resp))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	logger.Info("Logout request received")
	var req *dto.LogoutRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.MsgInvalidRequest))
		return
	}

	_, err := h.client.AuthClient.Logout(c.Request.Context(), req.ToProto())
	if err != nil {
		logger.Error("Failed to logout user", zap.String("refresh_token", req.RefreshToken), zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	logger.Info("Token refresh request received")
	var req *dto.RefreshRequestDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse(errs.MsgInvalidRequest))
		return
	}

	resp, err := h.client.AuthClient.Refresh(c.Request.Context(), req.ToProto())
	if err != nil {
		logger.Error("Failed to refresh tokens", zap.String("refresh_token", req.RefreshToken), zap.Error(err))
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	c.JSON(http.StatusOK, dto.RefreshResponseFromProto(resp))
}
