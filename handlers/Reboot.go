package handlers

import (
	"net/http"
	"os/exec"
	"time"

	"github.com/conthing/sysmgmt/dto"
	"github.com/conthing/utils/common"
	"github.com/gin-gonic/gin"
)

// Reboot 重启
func Reboot(c *gin.Context) {

	waitThenReboot()

	c.JSON(http.StatusOK, dto.Resp{
		Code: http.StatusOK,
		Data: "reboot in 3 seconds.",
	})
}

func waitThenReboot() {
	go func() {
		time.Sleep(time.Second * 3)
		_, err := exec.Command("reboot", "-f").Output() //初始化Cmd
		//_, err := exec.Command("bash", "-c", "reboot -f").Output() //初始化Cmd
		if err != nil {
			common.Log.Errorf("reboot failed: %v", err)
		}
	}()
}
