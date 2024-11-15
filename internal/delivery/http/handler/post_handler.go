package handler

import (
	"net/http"
	"strconv"

	"github.com/Dang-Hai-Tran/newfeed-go/internal/delivery/http/dto"
	"github.com/Dang-Hai-Tran/newfeed-go/internal/delivery/http/middleware"
	"github.com/Dang-Hai-Tran/newfeed-go/internal/domain"
	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postUsecase domain.PostUsecase
}

func NewPostHandler(router *gin.RouterGroup, postUsecase domain.PostUsecase, authMiddleware *middleware.AuthMiddleware) {
	handler := &PostHandler{
		postUsecase: postUsecase,
	}

	// All routes require authentication
	protected := router.Group("")
	protected.Use(authMiddleware.AuthRequired())
	{
		protected.GET("/posts/:post_id", handler.GetPost)
		protected.POST("/posts", handler.CreatePost)
		protected.PUT("/posts/:post_id", handler.UpdatePost)
		protected.DELETE("/posts/:post_id", handler.DeletePost)
		protected.GET("/friends/:user_id/posts", handler.GetUserPosts)
		protected.GET("/users/:user_id/newsfeed", handler.GetNewsFeed)
	}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.Response{Success: false, Message: "unauthorized"})
		return
	}

	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: err.Error()})
		return
	}

	post := &domain.Post{
		UserID:   userID,
		Content:  req.Content,
		ImageURL: req.ImageURL,
	}

	if err := h.postUsecase.CreatePost(post); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Message: "post created successfully",
		Data:    dto.ToPostResponse(post),
	})
}

func (h *PostHandler) GetPost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: "invalid post id"})
		return
	}

	post, err := h.postUsecase.GetPost(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}
	if post == nil {
		c.JSON(http.StatusNotFound, dto.Response{Success: false, Message: "post not found"})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    dto.ToPostResponse(post),
	})
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
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

	var req dto.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: err.Error()})
		return
	}

	post := &domain.Post{
		ID:       postID,
		UserID:   userID,
		Content:  req.Content,
		ImageURL: req.ImageURL,
	}

	if err := h.postUsecase.UpdatePost(post); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "post updated successfully",
		Data:    dto.ToPostResponse(post),
	})
}

func (h *PostHandler) DeletePost(c *gin.Context) {
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

	// First get the post to check ownership
	post, err := h.postUsecase.GetPost(postID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}
	if post == nil {
		c.JSON(http.StatusNotFound, dto.Response{Success: false, Message: "post not found"})
		return
	}
	if post.UserID != userID {
		c.JSON(http.StatusForbidden, dto.Response{Success: false, Message: "not authorized to delete this post"})
		return
	}

	if err := h.postUsecase.DeletePost(postID); err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "post deleted successfully",
	})
}

func (h *PostHandler) GetUserPosts(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: "invalid user id"})
		return
	}

	var pagination dto.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: err.Error()})
		return
	}

	posts, err := h.postUsecase.GetUserPosts(userID, pagination.Page, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	// Convert to response DTOs
	postResponses := make([]*dto.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = dto.ToPostResponse(&post)
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    postResponses,
	})
}

func (h *PostHandler) GetNewsFeed(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.Response{Success: false, Message: "unauthorized"})
		return
	}

	var pagination dto.PaginationQuery
	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: err.Error()})
		return
	}

	posts, err := h.postUsecase.GetNewsFeed(userID, pagination.Page, pagination.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	// Convert to response DTOs
	postResponses := make([]*dto.PostResponse, len(posts))
	for i, post := range posts {
		postResponses[i] = dto.ToPostResponse(&post)
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Data:    postResponses,
	})
}
