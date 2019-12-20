package services

import (
	"net/http"
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

	for port := range portList {
		resp, err := http.Get("http://localhost:" + string(port) + "/api/v1/ping")

		if err != nil || resp.StatusCode != 200 {
			IsHealth = false
			return
		}
		defer resp.Body.Close()

	}
	IsHealth = true
}

// ScheduledHealthCheck 定时轮询任务
func ScheduledHealthCheck() {
	for {
		HealthCheck()
		common.Log.Info("健康检查结果: ", IsHealth)
		time.Sleep(30 * time.Second)
	}
}
