package handlers

import (
	"net/http"
	"os/exec"

	"github.com/conthing/sysmgmt/dto"
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

	c.JSON(http.StatusOK, dto.Resp{
		Code: http.StatusOK,
	})
}
