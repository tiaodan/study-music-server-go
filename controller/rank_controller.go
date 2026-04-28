package controller

import (
	"log"
	"net/http"
	"strconv"
	"study-music-server-go/models"
	"study-music-server-go/service"

	"github.com/gin-gonic/gin"
)

type RankController struct {
	rankService *service.RankService
}

func NewRankController() *RankController {
	return &RankController{
		rankService: service.NewRankService(),
	}
}

// ImportRank 导入排行榜数据
// POST /rank/import
func (c *RankController) ImportRank(ctx *gin.Context) {
	var req models.RankImportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "success": false, "message": "参数错误: " + err.Error()})
		return
	}
	resp := c.rankService.ImportRank(&req)
	if !resp.Success {
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

// GetRankList 获取榜单列表
// GET /rank/list?websiteId=2
func (c *RankController) GetRankList(ctx *gin.Context) {
	log.Printf("[IP:%s] 请求榜单列表 websiteId=%s", ctx.ClientIP(), ctx.Query("websiteId"))
	idStr := ctx.Query("websiteId")
	websiteId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "success": false, "message": "无效的websiteId"})
		return
	}
	resp := c.rankService.GetRankList(uint(websiteId))
	ctx.JSON(http.StatusOK, resp)
}

// GetRankDetail 获取榜单详情
// GET /rank/detail?websiteId=2&rankName=top500
func (c *RankController) GetRankDetail(ctx *gin.Context) {
	log.Printf("[IP:%s] 请求榜单详情 websiteId=%s rankName=%s", ctx.ClientIP(), ctx.Query("websiteId"), ctx.Query("rankName"))
	websiteIdStr := ctx.Query("websiteId")
	rankName := ctx.Query("rankName")
	websiteId, err := strconv.ParseUint(websiteIdStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "success": false, "message": "无效的websiteId"})
		return
	}
	if rankName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "success": false, "message": "rankName不能为空"})
		return
	}
	resp := c.rankService.GetRankDetail(uint(websiteId), rankName)
	ctx.JSON(http.StatusOK, resp)
}

// CheckRank 校验排行榜数据（不入库，只对比文件）
// POST /rank/check
func (c *RankController) CheckRank(ctx *gin.Context) {
	var req models.RankCheckRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "success": false, "message": "参数错误: " + err.Error()})
		return
	}
	resp := c.rankService.CheckRank(&req)
	if !resp.Success {
		ctx.JSON(http.StatusInternalServerError, resp)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}