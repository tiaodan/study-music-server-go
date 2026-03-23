package controller

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"study-music-server-go/models"
	"study-music-server-go/service"

	"github.com/gin-gonic/gin"
)

type SongController struct {
	songService *service.SongService
}

func NewSongController() *SongController {
	return &SongController{
		songService: service.NewSongService(),
	}
}

func (c *SongController) AddSong(ctx *gin.Context) {
	var req models.SongRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.songService.AddSong(&req)
	ctx.JSON(http.StatusOK, resp)
}

func (c *SongController) UpdateSong(ctx *gin.Context) {
	var req models.SongRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.songService.UpdateSong(&req)
	ctx.JSON(http.StatusOK, resp)
}

func (c *SongController) DeleteSong(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	resp := c.songService.DeleteSong(uint(id))
	ctx.JSON(http.StatusOK, resp)
}

func (c *SongController) SongOfId(ctx *gin.Context) {
	// 优先从 path 参数获取，其次从 query 参数获取
	idStr := ctx.Param("id")
	if idStr == "" {
		idStr = ctx.Query("id")
	}
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	resp := c.songService.SongOfId(uint(id))
	songData, ok := resp.Data.(*models.Song)
	if !ok {
		ctx.JSON(http.StatusOK, resp)
		return
	}

	// 如果是 nas 来源，返回音频流
	if songData.UrlSource == "nas" && songData.Url != "" {
		// SMB路径直接使用
		localPath := songData.Url

		// 打开文件
		file, err := os.Open(localPath)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "文件打开失败: " + err.Error()})
			return
		}
		defer file.Close()

		// 获取文件信息
		fileInfo, err := file.Stat()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "获取文件信息失败"})
			return
		}

		// 设置响应头
		ctx.Header("Content-Type", "audio/mpeg")
		ctx.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
		ctx.Header("Accept-Ranges", "bytes")

		// 返回音频流
		ctx.Header("Content-Disposition", "inline; filename="+filepath.Base(localPath))
		_, err = io.Copy(ctx.Writer, file)
		if err != nil {
			return
		}
		return
	}

	ctx.JSON(http.StatusOK, resp)
}

func (c *SongController) SongOfSingerId(ctx *gin.Context) {
	singerIdStr := ctx.Query("singerId")
	singerId, err := strconv.ParseUint(singerIdStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid singerId"})
		return
	}
	resp := c.songService.SongOfSingerId(uint(singerId))
	ctx.JSON(http.StatusOK, resp)
}

func (c *SongController) SongOfName(ctx *gin.Context) {
	name := ctx.Query("name")
	resp := c.songService.SongOfName(name)
	ctx.JSON(http.StatusOK, resp)
}

func (c *SongController) AllSong(ctx *gin.Context) {
	resp := c.songService.AllSong()
	ctx.JSON(http.StatusOK, resp)
}
