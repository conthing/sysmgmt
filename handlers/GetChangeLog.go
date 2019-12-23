package handlers

import (
	"net/http"
	"sysmgmt-next/dto"

	"github.com/gin-gonic/gin"
)

// GetChangeLog 获取版本信息
func GetChangeLog(c *gin.Context) {
	var changelog dto.Changelog
	changelog.Version = "1.1.0"
	changelog.DescriptionList = []string{
		"优化 Web",
		"优化 Sysmgmt 结构，加入了指示灯功能",
		"增加对整体软件版本的管理功能",
	}
	c.JSON(http.StatusOK, changelog)
}
