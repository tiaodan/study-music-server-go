package controller

import (
	"net/http"
	"strconv"
	"study-music-server-go/models"
	"study-music-server-go/service"

	"github.com/gin-gonic/gin"
)

type RankListController struct {
	rankListService *service.RankListService
}

func NewRankListController() *RankListController {
	return &RankListController{
		rankListService: service.NewRankListService(),
	}
}

func (c *RankListController) AddRankList(ctx *gin.Context) {
	var req models.RankListRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.rankListService.AddRankList(&req)
	ctx.JSON(http.StatusOK, resp)
}

func (c *RankListController) RankListOfSongListId(ctx *gin.Context) {
	songListIdStr := ctx.Query("songListId")
	songListId, _ := strconv.ParseUint(songListIdStr, 10, 32)
	resp := c.rankListService.RankListOfSongListId(uint(songListId))
	ctx.JSON(http.StatusOK, resp)
}
