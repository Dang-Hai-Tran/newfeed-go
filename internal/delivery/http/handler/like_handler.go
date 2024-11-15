package handler

import (
	"net/http"
	"strconv"

	"github.com/datran42/newfeed-go/internal/delivery/http/dto"
	"github.com/datran42/newfeed-go/internal/delivery/http/middleware"
	"github.com/datran42/newfeed-go/internal/domain"
	"github.com/gin-gonic/gin"
)

type LikeHandler struct {
	likeUsecase domain.LikeUsecase
}

func NewLikeHandler(router *gin.RouterGroup, likeUsecase domain.LikeUsecase, authMiddleware *middleware.AuthMiddleware) {
	handler := &LikeHandler{
		likeUsecase: likeUsecase,
	}

	// All routes require authentication
	protected := router.Group("")
	protected.Use(authMiddleware.AuthRequired())
	{
		protected.GET("/posts/:post_id/likes", handler.GetPostLikes)
		protected.POST("/posts/:post_id/likes", handler.LikePost)
		protected.DELETE("/posts/:post_id/likes", handler.UnlikePost)
	}
}

func (h *LikeHandler) LikePost(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.Response{Success: false, Message: "unauthorized"})
		return
	}

	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: "invalid post id"})
		return
	}

	if err := h.likeUsecase.LikePost(postID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "post liked successfully",
	})
}

func (h *LikeHandler) UnlikePost(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.Response{Success: false, Message: "unauthorized"})
		return
	}

	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: "invalid post id"})
		return
	}

	if err := h.likeUsecase.UnlikePost(postID, userID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "post unliked successfully",
	})
}

func (h *LikeHandler) GetPostLikes(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: "invalid post id"})
		return
	}

	var pagination dto.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: err.Error()})
		return
	}

	likes, err := h.likeUsecase.GetPostLikes(postID, pagination.Page, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	// Convert to response DTOs
	likeResponses := make([]*dto.LikeResponse, len(likes))
	for i, like := range likes {
		likeResponses[i] = dto.ToLikeResponse(&like)
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    likeResponses,
	})
}
