package services

import (
	"fmt"
	"github.com/conthing/utils/common"
	"net/http"
	"sysmgmt-next/config"
	"time"
)

// HealthCheck 健康检查
func HealthCheck() error {
	microservicelist := config.Conf.MicroServiceList
	for _, microservice := range microservicelist {
		if microservice.EnableHealth == true {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/v1/ping", microservice.Port))
			if err != nil || resp.StatusCode != 200 {
				common.Log.Error("微服务不健康的有: ", microservice.Name)
				return err
			}
			defer resp.Body.Close()
		} else {
			return nil
		}
	}
	return nil
}

// todo 改写此函数，将健康检查和LED控制放到同一个go程里
// ScheduledHealthCheck 定时轮询任务
func ScheduledHealthCheck() {
	go func() {
		setLed()
		HealthCheck()
		time.Sleep(30 * time.Second)
	}()
}
