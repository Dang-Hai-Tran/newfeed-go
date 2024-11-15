package handler

import (
	"net/http"
	"strconv"

	"github.com/datran42/newfeed-go/internal/delivery/http/dto"
	"github.com/datran42/newfeed-go/internal/delivery/http/middleware"
	"github.com/datran42/newfeed-go/internal/domain"
	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentUsecase domain.CommentUsecase
}

func NewCommentHandler(router *gin.RouterGroup, commentUsecase domain.CommentUsecase, authMiddleware *middleware.AuthMiddleware) {
	handler := &CommentHandler{
		commentUsecase: commentUsecase,
	}

	// All routes require authentication
	protected := router.Group("")
	protected.Use(authMiddleware.AuthRequired())
	{
		protected.GET("/posts/:post_id/comments", handler.GetPostComments)
		protected.POST("/posts/:post_id/comments", handler.CreateComment)
		protected.PUT("/posts/:post_id/comments/:comment_id", handler.UpdateComment)
		protected.DELETE("/posts/:post_id/comments/:comment_id", handler.DeleteComment)
	}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
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

	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: err.Error()})
		return
	}

	comment := &domain.Comment{
		PostID:  postID,
		UserID:  userID,
		Content: req.Content,
	}

	if err := h.commentUsecase.CreateComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Message: "comment created successfully",
		Data:    dto.ToCommentResponse(comment),
	})
}

func (h *CommentHandler) GetPostComments(c *gin.Context) {
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

	comments, err := h.commentUsecase.GetPostComments(postID, pagination.Page, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	// Convert to response DTOs
	commentResponses := make([]*dto.CommentResponse, len(comments))
	for i, comment := range comments {
		commentResponses[i] = dto.ToCommentResponse(&comment)
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    commentResponses,
	})
}

func (h *CommentHandler) UpdateComment(c *gin.Context) {
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

	commentID, err := strconv.ParseUint(c.Param("comment_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: "invalid comment id"})
		return
	}

	var req dto.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: err.Error()})
		return
	}

	comment := &domain.Comment{
		ID:      commentID,
		PostID:  postID,
		UserID:  userID,
		Content: req.Content,
	}

	if err := h.commentUsecase.UpdateComment(comment); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "comment updated successfully",
		Data:    dto.ToCommentResponse(comment),
	})
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.Response{Success: false, Message: "unauthorized"})
		return
	}

	commentID, err := strconv.ParseUint(c.Param("comment_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: "invalid comment id"})
		return
	}

	// First get the comment to check ownership
	comment, err := h.commentUsecase.GetComment(commentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}
	if comment == nil {
		c.JSON(http.StatusNotFound, dto.Response{Success: false, Message: "comment not found"})
		return
	}
	if comment.UserID != userID {
		c.JSON(http.StatusForbidden, dto.Response{Success: false, Message: "not authorized to delete this comment"})
		return
	}

	if err := h.commentUsecase.DeleteComment(commentID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "comment deleted successfully",
	})
}
