package controller

import (
	"study-music-server-go/common"
	"study-music-server-go/models"
	"study-music-server-go/service"
	"github.com/gin-gonic/gin"
)

type ImportController struct {
	importService *service.ImportService
}

func NewImportController() *ImportController {
	return &ImportController{
		importService: service.NewImportService(),
	}
}

// FormatName 名字格式化
// POST /api/import/format-name
func (c *ImportController) FormatName(ctx *gin.Context) {
	var req models.FormatNameRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, common.Error("参数错误"))
		return
	}

	if req.Path == "" {
		ctx.JSON(400, common.Error("路径不能为空"))
		return
	}

	result := c.importService.FormatName(req.Path)
	ctx.JSON(200, result)
}

// MoveFile 移动文件到HDD
// POST /api/import/move-file
func (c *ImportController) MoveFile(ctx *gin.Context) {
	var req models.MoveFileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, common.Error("参数错误"))
		return
	}

	if req.From == "" || req.To == "" {
		ctx.JSON(400, common.Error("源路径和目标路径不能为空"))
		return
	}

	result := c.importService.MoveFile(req.From, req.To)
	ctx.JSON(200, result)
}

// ImportSongs 规整进数据库
// POST /api/import/songs
func (c *ImportController) ImportSongs(ctx *gin.Context) {
	var req models.ImportSongsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, common.Error("参数错误"))
		return
	}

	if req.Path == "" {
		ctx.JSON(400, common.Error("路径不能为空"))
		return
	}

	result := c.importService.ImportSongs(req.Path)
	ctx.JSON(200, result)
}
