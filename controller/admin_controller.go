package controller

import (
	"net/http"
	"study-music-server-go/models"
	"study-music-server-go/service"

	"github.com/gin-gonic/gin"
)

type AdminController struct {
	adminService *service.AdminService
}

func NewAdminController() *AdminController {
	return &AdminController{
		adminService: service.NewAdminService(),
	}
}

func (c *AdminController) Login(ctx *gin.Context) {
	var req models.AdminRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp := c.adminService.Login(req.Username, req.Password)
	ctx.JSON(http.StatusOK, resp)
}
