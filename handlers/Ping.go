package handlers

import "github.com/gin-gonic/gin"

// Ping 检测服务是否正常
func Ping(c *gin.Context) {
	c.String(200, "pong")
}
