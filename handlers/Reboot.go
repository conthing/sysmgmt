package handlers

import (
	"os/exec"
	"time"

	"github.com/gin-gonic/gin"
)

// Reboot 重启
func Reboot(c *gin.Context) {
	go func() {
		time.Sleep(2 * time.Second)
		exec.Command("reboot", "-f") //初始化Cmd
	}()
	c.String(200, "OK")
}
