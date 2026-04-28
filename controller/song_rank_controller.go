package controller

import (
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"study-music-server-go/service"

	"github.com/gin-gonic/gin"
)

type SongRankController struct {
	songRankService *service.SongRankService
}

func NewSongRankController() *SongRankController {
	return &SongRankController{
		songRankService: service.NewSongRankService(),
	}
}

// SongOfId 获取排行榜歌曲详情（或直接返回音频流）
// GET /song-rank/:id
func (c *SongRankController) SongOfId(ctx *gin.Context) {
	log.Printf("[IP:%s] 请求播放排行榜歌曲 song-rank id=%s", ctx.ClientIP(), ctx.Param("id"))
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "success": false, "message": "无效的ID"})
		return
	}

	// 先获取歌曲信息，判断是否需要返回音频流
	file, contentType, err := c.songRankService.StreamSong(uint(id))
	if err == nil && file != nil {
		// 能打开文件，直接返回音频流
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "success": false, "message": "获取文件信息失败"})
			return
		}

		// 支持 Range 请求
		rangeHeader := ctx.GetHeader("Range")
		if rangeHeader != "" {
			parts := strings.Split(rangeHeader, "=")
			if len(parts) == 2 && parts[0] == "bytes" {
				rangeParts := strings.Split(parts[1], "-")
				if len(rangeParts) >= 1 {
					start, _ := strconv.ParseInt(rangeParts[0], 10, 64)
					var end int64 = fileInfo.Size() - 1
					if len(rangeParts) == 2 && rangeParts[1] != "" {
						end, _ = strconv.ParseInt(rangeParts[1], 10, 64)
					}

					ctx.Header("Content-Range", "bytes "+strconv.FormatInt(start, 10)+"-"+strconv.FormatInt(end, 10)+"/"+strconv.FormatInt(fileInfo.Size(), 10))
					ctx.Header("Content-Length", strconv.FormatInt(end-start+1, 10))
					ctx.Header("Accept-Ranges", "bytes")
					ctx.Header("Content-Type", contentType)
					ctx.Status(http.StatusPartialContent)

					file.Seek(start, io.SeekStart)
					io.CopyN(ctx.Writer, file, end-start+1)
					return
				}
			}
		}

		// 无 Range 请求，返回完整文件
		ctx.Header("Content-Type", contentType)
		ctx.Header("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
		ctx.Header("Accept-Ranges", "bytes")
		ctx.Header("Content-Disposition", "inline; filename="+filepath.Base(file.Name()))
		io.Copy(ctx.Writer, file)
		return
	}

	// 文件不存在或无法打开，返回 JSON 格式的歌曲信息
	resp := c.songRankService.SongOfId(uint(id))
	ctx.JSON(http.StatusOK, resp)
}