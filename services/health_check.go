package services

import (
	"fmt"
	"net/http"
	"os/exec"
	"sysmgmt-next/config"
	"time"

	"github.com/conthing/utils/common"
)

// IsHealth 检查健康
var IsHealth bool

// HealthCheck 健康检查
// 如果有一个微服务检查失败，直接返回false
func HealthCheck() {
	portList := config.Conf.ServicePortlist
	for _, port := range portList {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/v1/ping", port))

		if err != nil || resp.StatusCode != 200 {
			common.Log.Error("微服务不健康: ", fmt.Sprintf("http://localhost:%s/api/v1/ping", port))
			IsHealth = false
			return
		}
		defer resp.Body.Close()

	}
	IsHealth = true
}

// CheckLpr 检查 lpr
func CheckLpr() {
	out, err := exec.Command("/bin/sh", "-c", "pgrep -l lpr").Output()
	if err != nil {
		common.Log.Error("CheckLpr 执行 shell 失败: ", err)
	}
	str := string(out)
	common.Log.Info(str)
	if str == "" {
		_, err := exec.Command("cd /app/zap/lpr && LD_LIBRARY_PATH=. exec -a conthing-lpr ./lpr -d /app/log/lpr >/app/log/conthing-lpr.log 2>&1 &").Output()
		if err != nil {
			common.Log.Error("重启 lpr 失败", err)
		}
	}
}

// ScheduledHealthCheck 定时轮询任务
func ScheduledHealthCheck() {
	for {
		HealthCheck()
		CheckLpr()
		time.Sleep(30 * time.Second)
	}
}
