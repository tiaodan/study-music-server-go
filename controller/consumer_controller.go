package controller

import (
	"net/http"
	"strconv"
	"study-music-server-go/models"
	"study-music-server-go/service"

	"github.com/gin-gonic/gin"
)

type ConsumerController struct {
	consumerService *service.ConsumerService
}

func NewConsumerController() *ConsumerController {
	return &ConsumerController{
		consumerService: service.NewConsumerService(),
	}
}

func (c *ConsumerController) AddUser(ctx *gin.Context) {
	var req models.ConsumerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.consumerService.AddUser(&req)
	ctx.JSON(http.StatusOK, resp)
}

func (c *ConsumerController) LoginStatus(ctx *gin.Context) {
	var req models.ConsumerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.consumerService.LoginStatus(&req)
	ctx.JSON(http.StatusOK, resp)
}

func (c *ConsumerController) LoginEmailStatus(ctx *gin.Context) {
	var req models.ConsumerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.consumerService.LoginEmailStatus(&req)
	ctx.JSON(http.StatusOK, resp)
}

func (c *ConsumerController) AllUser(ctx *gin.Context) {
	resp := c.consumerService.AllUser()
	ctx.JSON(http.StatusOK, resp)
}

func (c *ConsumerController) UserOfId(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	resp := c.consumerService.UserOfId(uint(id))
	ctx.JSON(http.StatusOK, resp)
}

func (c *ConsumerController) DeleteUser(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	resp := c.consumerService.DeleteUser(uint(id))
	ctx.JSON(http.StatusOK, resp)
}

func (c *ConsumerController) UpdateUserMsg(ctx *gin.Context) {
	var req models.ConsumerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.consumerService.UpdateUserMsg(&req)
	ctx.JSON(http.StatusOK, resp)
}

func (c *ConsumerController) UpdatePassword(ctx *gin.Context) {
	var req models.ConsumerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.consumerService.UpdatePassword(&req)
	ctx.JSON(http.StatusOK, resp)
}

func (c *ConsumerController) UpdateUserAvatar(ctx *gin.Context) {
	idStr := ctx.PostForm("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid file"})
		return
	}

	// Save file
	filename := "./img/avatorImages/" + file.Filename
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "file save failed"})
		return
	}

	resp := c.consumerService.UpdateUserAvator(filename, uint(id))
	ctx.JSON(http.StatusOK, resp)
}
