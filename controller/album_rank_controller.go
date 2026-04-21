package controller

import (
	"net/http"
	"strconv"
	"study-music-server-go/models"
	"study-music-server-go/service"

	"github.com/gin-gonic/gin"
)

type AlbumRankController struct {
	albumRankService *service.AlbumRankService
}

func NewAlbumRankController() *AlbumRankController {
	return &AlbumRankController{
		albumRankService: service.NewAlbumRankService(),
	}
}

func (c *AlbumRankController) AllAlbum(ctx *gin.Context) {
	resp := c.albumRankService.AllAlbum()
	ctx.JSON(http.StatusOK, resp)
}

func (c *AlbumRankController) AlbumOfId(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "success": false, "message": "无效的ID"})
		return
	}
	resp := c.albumRankService.AlbumOfId(uint(id))
	ctx.JSON(http.StatusOK, resp)
}

func (c *AlbumRankController) AlbumsOfSingerId(ctx *gin.Context) {
	idStr := ctx.Query("singerId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "success": false, "message": "无效的ID"})
		return
	}
	resp := c.albumRankService.AlbumsOfSingerId(uint(id))
	ctx.JSON(http.StatusOK, resp)
}

func (c *AlbumRankController) AddAlbum(ctx *gin.Context) {
	var req models.AlbumRankRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "success": false, "message": "参数错误"})
		return
	}
	resp := c.albumRankService.AddAlbum(&req)
	ctx.JSON(http.StatusOK, resp)
}

func (c *AlbumRankController) UpdateAlbum(ctx *gin.Context) {
	var req models.AlbumRankRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "success": false, "message": "参数错误"})
		return
	}
	resp := c.albumRankService.UpdateAlbum(&req)
	ctx.JSON(http.StatusOK, resp)
}

func (c *AlbumRankController) DeleteAlbum(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "success": false, "message": "无效的ID"})
		return
	}
	resp := c.albumRankService.DeleteAlbum(uint(id))
	ctx.JSON(http.StatusOK, resp)
}
