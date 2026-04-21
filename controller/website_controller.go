package controller

import (
	"net/http"
	"study-music-server-go/service"

	"github.com/gin-gonic/gin"
)

type WebsiteController struct {
	websiteService *service.WebsiteService
}

func NewWebsiteController() *WebsiteController {
	return &WebsiteController{
		websiteService: service.NewWebsiteService(),
	}
}

func (c *WebsiteController) AllWebsite(ctx *gin.Context) {
	resp := c.websiteService.AllWebsite()
	ctx.JSON(http.StatusOK, resp)
}