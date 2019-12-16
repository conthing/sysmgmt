package handlers

import (
	"sysmgmt-next/dto"

	"github.com/conthing/utils/common"
	"github.com/gin-gonic/gin"
)

// GetMacAddr 获取 Mac
func GetMacAddr(c *gin.Context) {
	var resp dto.MacInfo
	resp.Mac = common.GetSerialNumber()
	c.JSON(200, resp)
}
