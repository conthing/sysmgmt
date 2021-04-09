package handlers

import (
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/conthing/sysmgmt/config"
	"github.com/conthing/sysmgmt/dto"

	"github.com/conthing/utils/common"
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
	c.JSON(http.StatusOK, dto.Resp{
		Code: http.StatusOK,
		Data: timeInfo,
	})
}

// PutTime 修改时间
func PutTime(c *gin.Context) {
	var info dto.NTPInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusOK, dto.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	if info.URL == "" {
		info.URL = NTPServerURL //todo 这个服务器地址要存储，供get的时候获取
	}
	shellPath := config.Conf.ShellPath + "modifiedTime.sh"
	time := time.Unix(info.Date, 0)
	command := exec.Command(shellPath, info.Type, time.Format("2006-01-02 15:04:05"), info.URL)
	out, err := command.Output()
	common.Log.Debugf("timedatectl output:%s,err:%v", string(out), err)

	c.JSON(http.StatusOK, dto.Resp{
		Code: http.StatusOK,
		Data: info,
	})
}
