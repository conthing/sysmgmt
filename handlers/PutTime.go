package handlers

import (
	"log"
	"net/http"
	"os/exec"
	"sysmgmt-next/config"
	"sysmgmt-next/dto"
	"time"

	"github.com/gin-gonic/gin"
)

// PutTime 修改时间
func PutTime(c *gin.Context) {
	var info dto.NTPInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if info.NtpURL == "" {
		info.NtpURL = NTPServerURL
	}
	shellPath := config.Conf.ShellPath + "modifiedTime.sh"
	time := time.Unix(info.Date, 0)
	command := exec.Command(shellPath, info.TimeType, time.Format("2006-01-02 15:04:05"), info.NtpURL)
	out, err := command.Output()
	log.Printf("timedatectl set-ntp true output:%s,err:%v", string(out), err)
}
