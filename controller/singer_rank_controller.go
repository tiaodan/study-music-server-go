package controller

import (
	"net/http"
	"strconv"
	"study-music-server-go/models"
	"study-music-server-go/service"

	"github.com/gin-gonic/gin"
)

type SingerRankController struct {
	singerRankService *service.SingerRankService
}

func NewSingerRankController() *SingerRankController {
	return &SingerRankController{
		singerRankService: service.NewSingerRankService(),
	}
}

func (c *SingerRankController) AllSinger(ctx *gin.Context) {
	resp := c.singerRankService.AllSinger()
	ctx.JSON(http.StatusOK, resp)
}

func (c *SingerRankController) SingerOfId(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "success": false, "message": "无效的ID"})
		return
	}
	resp := c.singerRankService.SingerOfId(uint(id))
	ctx.JSON(http.StatusOK, resp)
}

func (c *SingerRankController) SingerOfName(ctx *gin.Context) {
	name := ctx.Query("name")
	resp := c.singerRankService.SingerOfName(name)
	ctx.JSON(http.StatusOK, resp)
}

func (c *SingerRankController) AddSinger(ctx *gin.Context) {
	var req models.SingerRankRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "success": false, "message": "参数错误"})
		return
	}
	resp := c.singerRankService.AddSinger(&req)
	ctx.JSON(http.StatusOK, resp)
}

func (c *SingerRankController) UpdateSinger(ctx *gin.Context) {
	var req models.SingerRankRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "success": false, "message": "参数错误"})
		return
	}
	resp := c.singerRankService.UpdateSinger(&req)
	ctx.JSON(http.StatusOK, resp)
}

func (c *SingerRankController) DeleteSinger(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "success": false, "message": "无效的ID"})
		return
	}
	resp := c.singerRankService.DeleteSinger(uint(id))
	ctx.JSON(http.StatusOK, resp)
}
