package handlers

import (
	"net/http"

	"github.com/conthing/sysmgmt/dto"
	"github.com/conthing/sysmgmt/services"

	"github.com/gin-gonic/gin"
)

// GetNetInfo 获取系统网络
func GetNetInfo(c *gin.Context) {
	var info dto.NetInfo
	err := services.GetCurrentNetInfo(&info)
	if err != nil {
		c.JSON(http.StatusOK, dto.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Resp{
		Code: http.StatusOK,
		Data: info,
	})
}

// PutNet 修改IP
func PutNet(c *gin.Context) {
	var info dto.NetInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusOK, dto.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	if err := services.SetNetInfo(&info); err != nil {
		c.JSON(http.StatusOK, dto.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	services.NotifyLed()
	c.JSON(http.StatusOK, dto.Resp{
		Code: http.StatusOK,
		Data: info,
	})
}
