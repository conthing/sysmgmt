package handlers

import (
	"os/exec"
	"strings"
	"sysmgmt-next/dto"
	"time"

	"github.com/gin-gonic/gin"
)

// NTPServerURL 地址
const NTPServerURL = "cn.pool.ntp.org"

// GetTimeInfo 获取时间信息
func GetTimeInfo(c *gin.Context) {
	var timeInfo dto.TimeInfo
	timeInfo.NtpURL = NTPServerURL
	timeInfo.Time = time.Now().Unix()
	command := exec.Command("/bin/sh", "-c", `timedatectl | grep "NTP enabled"`)
	out, _ := command.Output()
	if len(out) == 0 {
		command = exec.Command("/bin/sh", "-c", `timedatectl | grep "Network time on"`)
		out, _ = command.Output()
	}
	if len(out) == 0 {
		command = exec.Command("/bin/sh", "-c", `timedatectl | grep "System clock synchronized"`)
		out, _ = command.Output()
	}
	status := (strings.Split(string(out), ": "))
	if status[1] == "yes\n" {
		timeInfo.Ntpstatus = true
	} else {
		timeInfo.Ntpstatus = false
	}
	c.JSON(200, timeInfo)
}
