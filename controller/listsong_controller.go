package controller

import (
	"net/http"
	"strconv"
	"study-music-server-go/models"
	"study-music-server-go/service"

	"github.com/gin-gonic/gin"
)

type ListSongController struct {
	listSongService *service.ListSongService
}

func NewListSongController() *ListSongController {
	return &ListSongController{
		listSongService: service.NewListSongService(),
	}
}

func (c *ListSongController) AddListSong(ctx *gin.Context) {
	var req models.ListSongRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.listSongService.AddListSong(&req)
	ctx.JSON(http.StatusOK, resp)
}

func (c *ListSongController) DeleteListSong(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)
	resp := c.listSongService.DeleteListSong(uint(id))
	ctx.JSON(http.StatusOK, resp)
}

func (c *ListSongController) ListSongOfSongListId(ctx *gin.Context) {
	songListIdStr := ctx.Query("songListId")
	songListId, _ := strconv.ParseUint(songListIdStr, 10, 32)
	resp := c.listSongService.ListSongOfSongListId(uint(songListId))
	ctx.JSON(http.StatusOK, resp)
}
