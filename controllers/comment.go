package controllers

import (
	"net/http"
	"strconv"
	"taskive/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CommentController struct {
	commentService *services.CommentService
	validate       *validator.Validate
}

func NewCommentController(commentService *services.CommentService) *CommentController {
	return &CommentController{
		commentService: commentService,
		validate:       validator.New(),
	}
}

func (c *CommentController) Create(ctx *gin.Context) {
	taskID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	var input services.CreateCommentInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := c.validate.Struct(input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := ctx.GetUint("user_id")
	comment, err := c.commentService.Create(uint(taskID), userID, input)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, comment)
}

func (c *CommentController) GetTaskComments(ctx *gin.Context) {
	taskID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	comments, err := c.commentService.GetTaskComments(uint(taskID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, comments)
}

func (c *CommentController) Delete(ctx *gin.Context) {
	commentID, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment ID"})
		return
	}

	if err := c.commentService.Delete(uint(commentID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusNoContent)
} 