package handler

import (
	"net/http"
	"strconv"

	"github.com/Dang-Hai-Tran/newfeed-go/internal/delivery/http/dto"
	"github.com/Dang-Hai-Tran/newfeed-go/internal/delivery/http/middleware"
	"github.com/Dang-Hai-Tran/newfeed-go/internal/domain"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(router *gin.RouterGroup, userUsecase domain.UserUsecase, authMiddleware *middleware.AuthMiddleware) {
	handler := &UserHandler{
		userUsecase: userUsecase,
	}

	// Public routes
	router.POST("/sessions", handler.Login)
	router.POST("/users", handler.Register)

	// Protected routes
	protected := router.Group("")
	protected.Use(authMiddleware.AuthRequired())
	{
		protected.GET("/users/:user_id", handler.GetProfile)
		protected.PUT("/users", handler.UpdateProfile)
		protected.DELETE("/users", handler.DeleteProfile)
		protected.GET("/friends/:user_id", handler.GetFollowers)
		protected.POST("/friends/:user_id", handler.Follow)
		protected.DELETE("/friends/:user_id", handler.Unfollow)
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: err.Error()})
		return
	}

	user := &domain.User{
		Username:  req.Username,
		Password:  req.Password,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Birthday:  req.Birthday,
	}

	if err := h.userUsecase.Register(user); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Message: "user registered successfully",
		Data:    dto.ToUserResponse(user),
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: err.Error()})
		return
	}

	token, err := h.userUsecase.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Response{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    dto.LoginResponse{Token: token},
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: "invalid user id"})
		return
	}

	user, err := h.userUsecase.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    dto.ToUserResponse(user),
	})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.Response{Success: false, Message: "unauthorized"})
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: err.Error()})
		return
	}

	user := &domain.User{
		ID:        userID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Birthday:  req.Birthday,
		Password:  req.Password,
	}

	if err := h.userUsecase.UpdateProfile(user); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "profile updated successfully",
	})
}

func (h *UserHandler) DeleteProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.Response{Success: false, Message: "unauthorized"})
		return
	}

	if err := h.userUsecase.DeleteProfile(userID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "profile deleted successfully",
	})
}

func (h *UserHandler) Follow(c *gin.Context) {
	followerID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.Response{Success: false, Message: "unauthorized"})
		return
	}

	followingID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: "invalid user id"})
		return
	}

	if err := h.userUsecase.Follow(followerID, followingID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "followed successfully",
	})
}

func (h *UserHandler) Unfollow(c *gin.Context) {
	followerID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.Response{Success: false, Message: "unauthorized"})
		return
	}

	followingID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: "invalid user id"})
		return
	}

	if err := h.userUsecase.Unfollow(followerID, followingID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "unfollowed successfully",
	})
}

func (h *UserHandler) GetFollowers(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: "invalid user id"})
		return
	}

	followers, err := h.userUsecase.GetFollowers(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	// Convert to response DTOs
	followerResponses := make([]*dto.UserResponse, len(followers))
	for i, follower := range followers {
		followerResponses[i] = dto.ToUserResponse(&follower)
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    followerResponses,
	})
}
