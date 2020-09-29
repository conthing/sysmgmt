package handlers

import (
	"net/http"

	"github.com/conthing/sysmgmt/dto"

	"github.com/conthing/utils/common"
	"github.com/gin-gonic/gin"
)

// GetMac 获取 Mac
func GetMac(c *gin.Context) {
	var resp dto.MacInfo
	resp.Mac = common.GetSerialNumber()

	c.JSON(http.StatusOK, dto.Resp{
		Code: http.StatusOK,
		Data: resp,
	})
}
