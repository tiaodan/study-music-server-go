package controller

import (
	"io"
	"net/http"
	"strconv"
	"study-music-server-go/models"
	"study-music-server-go/service"

	"github.com/gin-gonic/gin"
)

type S3Controller struct {
	s3Service *service.S3Service
}

func NewS3Controller() *S3Controller {
	return &S3Controller{
		s3Service: service.NewS3Service(),
	}
}

// ListFolders GET /s3/folders?prefix=
func (c *S3Controller) ListFolders(ctx *gin.Context) {
	prefix := ctx.Query("prefix")
	resp := c.s3Service.ListFolders(prefix)
	ctx.JSON(http.StatusOK, resp)
}

// CreateFolder POST /s3/folder
func (c *S3Controller) CreateFolder(ctx *gin.Context) {
	var req models.CreateFolderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.s3Service.CreateFolder(req.Path)
	ctx.JSON(http.StatusOK, resp)
}

// DeleteFolder DELETE /s3/folder
func (c *S3Controller) DeleteFolder(ctx *gin.Context) {
	var req models.DeleteFolderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.s3Service.DeleteFolder(req.Path, req.Recursive)
	ctx.JSON(http.StatusOK, resp)
}

// GetFolderInfo GET /s3/folder/info?path=
func (c *S3Controller) GetFolderInfo(ctx *gin.Context) {
	path := ctx.Query("path")
	resp := c.s3Service.GetFolderInfo(path)
	ctx.JSON(http.StatusOK, resp)
}

// ListObjects GET /s3/objects?path=&prefix=&limit=
func (c *S3Controller) ListObjects(ctx *gin.Context) {
	path := ctx.Query("path")
	prefix := ctx.Query("prefix")
	limit := parseIntDefault(ctx.Query("limit"), 100)
	resp := c.s3Service.ListObjects(path, prefix, limit)
	ctx.JSON(http.StatusOK, resp)
}

// UploadFile POST /s3/upload (multipart/form-data)
func (c *S3Controller) UploadFile(ctx *gin.Context) {
	folder := ctx.PostForm("folder")
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "文件上传失败"})
		return
	}
	defer file.Close()

	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "读取文件失败"})
		return
	}

	resp := c.s3Service.UploadFile(folder, header.Filename, content)
	ctx.JSON(http.StatusOK, resp)
}

// DeleteObject DELETE /s3/object
func (c *S3Controller) DeleteObject(ctx *gin.Context) {
	var req models.DeleteObjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.s3Service.DeleteObject(req.Key)
	ctx.JSON(http.StatusOK, resp)
}

// GetObjectInfo GET /s3/object/info?key=
func (c *S3Controller) GetObjectInfo(ctx *gin.Context) {
	key := ctx.Query("key")
	resp := c.s3Service.GetObjectInfo(key)
	ctx.JSON(http.StatusOK, resp)
}

// DownloadObject GET /s3/object/download?key=&expire=
func (c *S3Controller) DownloadObject(ctx *gin.Context) {
	key := ctx.Query("key")
	expire := parseIntDefault(ctx.Query("expire"), 15)
	resp := c.s3Service.DownloadObject(key, expire)
	ctx.JSON(http.StatusOK, resp)
}

// CopyObject POST /s3/object/copy
func (c *S3Controller) CopyObject(ctx *gin.Context) {
	var req models.CopyObjectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.s3Service.CopyObject(req.SourceKey, req.DestKey)
	ctx.JSON(http.StatusOK, resp)
}

// BatchDeleteObjects DELETE /s3/objects/batch
func (c *S3Controller) BatchDeleteObjects(ctx *gin.Context) {
	var req models.BatchDeleteRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.s3Service.BatchDeleteObjects(req.Keys)
	ctx.JSON(http.StatusOK, resp)
}

// BatchCopyObjects POST /s3/objects/batch/copy
func (c *S3Controller) BatchCopyObjects(ctx *gin.Context) {
	var req models.BatchCopyRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.s3Service.BatchCopyObjects(req.Items)
	ctx.JSON(http.StatusOK, resp)
}

// BatchGetObjectsInfo POST /s3/objects/batch/info
func (c *S3Controller) BatchGetObjectsInfo(ctx *gin.Context) {
	var req models.BatchGetInfoRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.s3Service.BatchGetObjectsInfo(req.Keys)
	ctx.JSON(http.StatusOK, resp)
}

// parseIntDefault 解析整数，失败返回默认值
func parseIntDefault(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return val
}