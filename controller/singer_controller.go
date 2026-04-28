package controller

import (
	"log"
	"net/http"
	"strconv"
	"study-music-server-go/models"
	"study-music-server-go/service"

	"github.com/gin-gonic/gin"
)

type SingerController struct {
	singerService *service.SingerService
}

func NewSingerController() *SingerController {
	return &SingerController{
		singerService: service.NewSingerService(),
	}
}

func (c *SingerController) AddSinger(ctx *gin.Context) {
	var req models.SingerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.singerService.AddSinger(&req)
	ctx.JSON(http.StatusOK, resp)
}

func (c *SingerController) UpdateSinger(ctx *gin.Context) {
	var req models.SingerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.singerService.UpdateSinger(&req)
	ctx.JSON(http.StatusOK, resp)
}

func (c *SingerController) DeleteSinger(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	resp := c.singerService.DeleteSinger(uint(id))
	ctx.JSON(http.StatusOK, resp)
}

func (c *SingerController) SingerOfId(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	resp := c.singerService.SingerOfId(uint(id))
	ctx.JSON(http.StatusOK, resp)
}

func (c *SingerController) SingerOfName(ctx *gin.Context) {
	name := ctx.Query("name")
	resp := c.singerService.SingerOfName(name)
	ctx.JSON(http.StatusOK, resp)
}

func (c *SingerController) AllSinger(ctx *gin.Context) {
	log.Printf("[IP:%s] 请求歌手列表", ctx.ClientIP())
	resp := c.singerService.AllSinger()
	ctx.JSON(http.StatusOK, resp)
}

// SingerJay 只返回周杰伦（临时测试用）
func (c *SingerController) SingerJay(ctx *gin.Context) {
	log.Printf("[IP:%s] 请求歌手列表(jay)", ctx.ClientIP())
	resp := c.singerService.SingerJay()
	ctx.JSON(http.StatusOK, resp)
}

// AlbumsOfSingerId 查询歌手的所有专辑
func (c *SingerController) AlbumsOfSingerId(ctx *gin.Context) {
	idStr := ctx.Query("singerId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid singerId"})
		return
	}
	resp := c.singerService.AlbumsBySingerId(uint(id))
	ctx.JSON(http.StatusOK, resp)
}
