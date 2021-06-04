package handlers

import (
	"net/http"

	"github.com/conthing/sysmgmt/services"
	"github.com/gin-gonic/gin"
)

// Ping 检测服务是否正常
func Ping(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}

// Identify LED闪烁
func Identify(c *gin.Context) {
	services.IdentifyLed()
	c.String(http.StatusOK, "All indicators flashing for 30 seconds")
}
