package services

import (
	"fmt"
	"net/http"
	"sysmgmt-next/config"
	"time"

	"github.com/conthing/utils/common"
)

// todo　删除此变量
// IsHealth 检查健康
var IsHealth bool

// HealthCheck 健康检查
// 如果有一个微服务检查失败，直接返回false
// todo 此函数改成不健康是返回 error，并且将所有不健康的服务的名字都在error信息中体现出来
// func HealthCheck() error
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

// todo 改写此函数，将健康检查和LED控制放到同一个go程里
// ScheduledHealthCheck 定时轮询任务
func ScheduledHealthCheck() {
	for {
		HealthCheck()
		time.Sleep(30 * time.Second)
	}
}
