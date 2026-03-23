package controller

import (
	"net/http"
	"strconv"
	"study-music-server-go/models"
	"study-music-server-go/service"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	commentService *service.CommentService
}

func NewCommentController() *CommentController {
	return &CommentController{
		commentService: service.NewCommentService(),
	}
}

func (c *CommentController) AddComment(ctx *gin.Context) {
	var req models.CommentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.commentService.AddComment(&req)
	ctx.JSON(http.StatusOK, resp)
}

func (c *CommentController) DeleteComment(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	resp := c.commentService.DeleteComment(uint(id))
	ctx.JSON(http.StatusOK, resp)
}

func (c *CommentController) CommentOfSongId(ctx *gin.Context) {
	songIdStr := ctx.Query("songId")
	songId, _ := strconv.ParseUint(songIdStr, 10, 32)
	resp := c.commentService.CommentOfSongId(uint(songId))
	ctx.JSON(http.StatusOK, resp)
}

func (c *CommentController) CommentOfSongListId(ctx *gin.Context) {
	songListIdStr := ctx.Query("songListId")
	songListId, _ := strconv.ParseUint(songListIdStr, 10, 32)
	resp := c.commentService.CommentOfSongListId(uint(songListId))
	ctx.JSON(http.StatusOK, resp)
}
