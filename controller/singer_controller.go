package controller

import (
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
	resp := c.singerService.AllSinger()
	ctx.JSON(http.StatusOK, resp)
}
