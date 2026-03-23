package controller

import (
	"net/http"
	"strconv"
	"study-music-server-go/models"
	"study-music-server-go/service"

	"github.com/gin-gonic/gin"
)

type SongListController struct {
	songListService *service.SongListService
}

func NewSongListController() *SongListController {
	return &SongListController{
		songListService: service.NewSongListService(),
	}
}

func (c *SongListController) AddSongList(ctx *gin.Context) {
	var req models.SongListRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.songListService.AddSongList(&req)
	ctx.JSON(http.StatusOK, resp)
}

func (c *SongListController) UpdateSongList(ctx *gin.Context) {
	var req models.SongListRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.songListService.UpdateSongList(&req)
	ctx.JSON(http.StatusOK, resp)
}

func (c *SongListController) DeleteSongList(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	resp := c.songListService.DeleteSongList(uint(id))
	ctx.JSON(http.StatusOK, resp)
}

func (c *SongListController) SongListOfId(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	resp := c.songListService.SongListOfId(uint(id))
	ctx.JSON(http.StatusOK, resp)
}

func (c *SongListController) SongListOfTitle(ctx *gin.Context) {
	title := ctx.Query("title")
	resp := c.songListService.SongListOfTitle(title)
	ctx.JSON(http.StatusOK, resp)
}

func (c *SongListController) AllSongList(ctx *gin.Context) {
	resp := c.songListService.AllSongList()
	ctx.JSON(http.StatusOK, resp)
}
