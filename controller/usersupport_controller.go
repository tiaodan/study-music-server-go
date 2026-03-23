package controller

import (
	"net/http"
	"strconv"
	"study-music-server-go/models"
	"study-music-server-go/service"

	"github.com/gin-gonic/gin"
)

type UserSupportController struct {
	userSupportService *service.UserSupportService
}

func NewUserSupportController() *UserSupportController {
	return &UserSupportController{
		userSupportService: service.NewUserSupportService(),
	}
}

func (c *UserSupportController) AddUserSupport(ctx *gin.Context) {
	var req models.UserSupportRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.userSupportService.AddUserSupport(&req)
	ctx.JSON(http.StatusOK, resp)
}

func (c *UserSupportController) DeleteUserSupport(ctx *gin.Context) {
	userIdStr := ctx.Query("userId")
	commentIdStr := ctx.Query("commentId")
	userId, _ := strconv.ParseUint(userIdStr, 10, 32)
	commentId, _ := strconv.ParseUint(commentIdStr, 10, 32)
	resp := c.userSupportService.DeleteUserSupport(uint(userId), uint(commentId))
	ctx.JSON(http.StatusOK, resp)
}
