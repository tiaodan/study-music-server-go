package controller

import (
	"net/http"
	"strconv"
	"study-music-server-go/models"
	"study-music-server-go/service"

	"github.com/gin-gonic/gin"
)

type CollectController struct {
	collectService *service.CollectService
}

func NewCollectController() *CollectController {
	return &CollectController{
		collectService: service.NewCollectService(),
	}
}

func (c *CollectController) AddCollect(ctx *gin.Context) {
	var req models.CollectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.collectService.AddCollect(&req)
	ctx.JSON(http.StatusOK, resp)
}

func (c *CollectController) DeleteCollect(ctx *gin.Context) {
	userIdStr := ctx.Query("userId")
	songIdStr := ctx.Query("songId")
	userId, _ := strconv.ParseUint(userIdStr, 10, 32)
	songId, _ := strconv.ParseUint(songIdStr, 10, 32)
	resp := c.collectService.DeleteCollect(uint(userId), uint(songId))
	ctx.JSON(http.StatusOK, resp)
}

func (c *CollectController) CollectOfUserId(ctx *gin.Context) {
	userIdStr := ctx.Query("userId")
	userId, _ := strconv.ParseUint(userIdStr, 10, 32)
	resp := c.collectService.CollectOfUserId(uint(userId))
	ctx.JSON(http.StatusOK, resp)
}
