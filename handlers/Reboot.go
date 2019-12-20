package handlers

import (
	"os/exec"

	"github.com/conthing/utils/common"
	"github.com/gin-gonic/gin"
)

// Reboot 重启
func Reboot(c *gin.Context) {

	data, err := exec.Command("bash", "-c", "reboot -f").Output() //初始化Cmd
	if err != nil {
		common.Log.Error(err)
		c.String(400, "关机失败", err)
		return
	}
	common.Log.Debug(string(data))
	c.String(200, "OK")
}
