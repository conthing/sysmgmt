package handlers

import (
	"net/http"
	"os/exec"
	"time"

	"github.com/conthing/sysmgmt/dto"
	"github.com/conthing/sysmgmt/services"

	"github.com/gin-gonic/gin"
)

// 原方式问题：
// 写入 NTP=FALSE 无效，至少易庐板子如此
// systemctl start systemd-timesyncd
// 会出现以下错误， journalctl -f
// May 28 19:56:37 conthing systemd[23991]: systemd-timesyncd.service: Failed to set up special execution directory in /var/lib: File exists
// May 28 19:56:37 conthing systemd[23991]: systemd-timesyncd.service: Failed at step STATE_DIRECTORY spawning /lib/systemd/systemd-timesyncd: File exists

// GetTimeInfo 获取时间信息
func GetTimeInfo(c *gin.Context) {
	var timeInfo dto.TimeInfo
	if enable, err := services.GetNtpEnable(); err == nil {
		timeInfo.NTPEnable = enable
	} else {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}
	if server, err := services.GetNtpServer(); err == nil {
		timeInfo.NTPServer = server
	} else {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	timeInfo.Time = time.Now().Unix()
	c.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Data: timeInfo,
	})
}

// SetTime 修改时间
func SetTime(c *gin.Context) {
	var info dto.TimeInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusOK, Response{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	if info.NTPEnable {
		if info.NTPServer != "" {
			if err := services.SetNtpServer(info.NTPServer); err != nil {
				c.JSON(http.StatusOK, Response{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				})
				return
			}
		}
		if err := services.SetNtpEnable(true); err != nil {
			c.JSON(http.StatusOK, Response{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}
	} else {
		if err := services.SetNtpEnable(false); err != nil {
			c.JSON(http.StatusOK, Response{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}
		if info.Time != 0 {
			time := time.Unix(info.Time, 0)
			command := exec.Command("date", "-s", time.Format("2006-01-02 15:04:05"))
			if _, err := command.Output(); err != nil {
				c.JSON(http.StatusOK, Response{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				})
				return
			}
			command = exec.Command("hwclock", "-w")
			if _, err := command.Output(); err != nil {
				c.JSON(http.StatusOK, Response{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				})
				return
			}
		}

	}

	c.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Data: info,
	})
}
