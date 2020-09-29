package handlers

import (
	"net/http"

	"github.com/conthing/sysmgmt/dto"
	"github.com/conthing/sysmgmt/services"

	"github.com/gin-gonic/gin"
)

// GetVersion 获取版本信息
func GetVersion(c *gin.Context) {
	c.JSON(http.StatusOK, dto.Resp{
		Code: http.StatusOK,
		Data: services.GetAllVersion(),
	})
}
