package handlers

import (
	"net/http"
	"github.com/conthing/sysmgmt/services"

	"github.com/gin-gonic/gin"
)

// GetVersion 获取版本信息
func GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, services.GetAllVersion())
}
