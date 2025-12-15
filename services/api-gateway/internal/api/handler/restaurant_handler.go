package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/api/dto"
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/client"
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/domain"
	"github.com/tamirat-dejene/ha-soranu/services/api-gateway/internal/errs"
	"github.com/tamirat-dejene/ha-soranu/shared/protos/restaurantpb"
)

type RestaurantHandler struct {
	client *client.RestaurantServiceClient
}

func NewRestaurantHandler(client *client.RestaurantServiceClient) *RestaurantHandler {
	return &RestaurantHandler{
		client: client,
	}
}

/**
type RestaurantServiceClient interface {
    Login(ctx context.Context, in *RestaurantLoginRequest, opts ...grpc.CallOption) (*RestaurantLoginResponse, error)
    RegisterRestaurant(ctx context.Context, in *RegisterRestaurantRequest, opts ...grpc.CallOption) (*RegisterRestaurantResponse, error)
    GetRestaurant(ctx context.Context, in *GetRestaurantRequest, opts ...grpc.CallOption) (*GetRestaurantResponse, error)
    ListRestaurants(ctx context.Context, in *ListRestaurantsRequest, opts ...grpc.CallOption) (*ListRestaurantsResponse, error)
    AddMenuItem(ctx context.Context, in *AddMenuItemRequest, opts ...grpc.CallOption) (*MenuItem, error)
    RemoveMenuItem(ctx context.Context, in *RemoveMenuItemRequest, opts ...grpc.CallOption) (*MenuItem, error)
    UpdateMenuItem(ctx context.Context, in *UpdateMenuItemRequest, opts ...grpc.CallOption) (*MenuItem, error)
}
*/

func (h *RestaurantHandler) Login(c *gin.Context) {
	var req dto.RestaurantLoginDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errs.MsgInvalidRequest)
		return
	}

	resp, err := h.client.RestaurantClient.Login(c.Request.Context(), req.ToProto())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	c.JSON(http.StatusOK, dto.RestaurantResponseFromProto(resp))
}

func (h *RestaurantHandler) RegisterRestaurant(c *gin.Context) {
	var req dto.RegisterRestaurantDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errs.MsgInvalidRequest)
		return
	}

	resp, err := h.client.RestaurantClient.RegisterRestaurant(c.Request.Context(), req.ToProto())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	c.JSON(http.StatusOK, dto.RestaurantResponseFromProto(resp))
}

func (h *RestaurantHandler) GetRestaurant(c *gin.Context) {
	restaurantID := c.Query("restaurant_id")
	if restaurantID == "" {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse("restaurant id is required"))
		return
	}

	req := &restaurantpb.GetRestaurantRequest{RestaurantId: restaurantID}

	resp, err := h.client.RestaurantClient.GetRestaurant(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	c.JSON(http.StatusOK, dto.RestaurantResponseFromProto(resp))
}

func (h *RestaurantHandler) ListRestaurants(c *gin.Context) {
	var req dto.ListRestaurantsDTO

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errs.MsgInvalidRequest)
		return
	}

	if req.Latitude == nil || req.Longitude == nil || req.RadiusKm == nil {
		c.JSON(http.StatusBadRequest, errs.MsgInvalidRequest)
		return
	}

	if *req.RadiusKm <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "radius_km must be greater than zero",
		})
		return
	}

	if *req.Latitude < -90 || *req.Latitude > 90 ||
		*req.Longitude < -180 || *req.Longitude > 180 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid coordinates",
		})
		return
	}

	ctx := c.Request.Context()

	stream, err := h.client.RestaurantClient.
		ListRestaurants(ctx, req.ToProto())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	var restaurants []*domain.Restaurant

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
			return
		}

		restaurants = append(
			restaurants,
			dto.RestaurantResponseFromProto(res),
		)
	}

	c.JSON(http.StatusOK, gin.H{
		"restaurants": restaurants,
	})
}


func (h *RestaurantHandler) AddMenuItem(c *gin.Context) {
	var menuItem dto.AddMenuItemDTO
	if err := c.ShouldBindJSON(&menuItem); err != nil {
		c.JSON(http.StatusBadRequest, errs.MsgInvalidRequest)
		return
	}

	resp, err := h.client.RestaurantClient.AddMenuItem(c.Request.Context(), menuItem.ToProto())
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	c.JSON(http.StatusOK, dto.MenuItemResponseFromProto(resp))
}

func (h *RestaurantHandler) RemoveMenuItem(c *gin.Context) {
	restaurantID := c.Query("restaurant_id")
	itemID := c.Query("item_id")

	if restaurantID == "" || itemID == "" {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse("restaurant_id and item_id are required"))
		return
	}

	req := &restaurantpb.RemoveMenuItemRequest{
		RestaurantId: restaurantID,
		ItemId:       itemID,
	}

	resp, err := h.client.RestaurantClient.RemoveMenuItem(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	c.JSON(http.StatusOK, dto.MenuItemResponseFromProto(resp))
}

func (h *RestaurantHandler) UpdateMenuItem(c *gin.Context) {
	restaurantID := c.Query("restaurant_id")
	itemID := c.Query("item_id")

	if restaurantID == "" || itemID == "" {
		c.JSON(http.StatusBadRequest, errs.NewErrorResponse("restaurant_id and item_id are required"))
		return
	}

	var updateItem dto.UpdateMenuItemDTO
	if err := c.ShouldBindJSON(&updateItem); err != nil {
		c.JSON(http.StatusBadRequest, errs.MsgInvalidRequest)
		return
	}

	updateProto := &restaurantpb.UpdateMenuItemRequest{
		RestaurantId: restaurantID,
		ItemId: 	itemID,
		Name:         updateItem.Name,
		Description:  updateItem.Description,
		Price:        updateItem.Price,
	}

	resp, err := h.client.RestaurantClient.UpdateMenuItem(c.Request.Context(), updateProto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponseFromGRPCError(err))
		return
	}

	c.JSON(http.StatusOK, dto.MenuItemResponseFromProto(resp))
}
