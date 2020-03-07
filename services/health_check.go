package services

import (
	"fmt"
	"net/http"
	"sysmgmt-next/config"
	"time"

	"github.com/conthing/utils/common"
)

// todo review 这里返回故障时要做到：如果两个微服务ping失败，error信息中String为"xxx yyy health check failed"，有点难，不要花太长时间，我会教你的
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
		setLed() // todo review 因为setLed不对，所以这个写法不对；应该讲你现在实现的setLed里的一些内容搬过来
		HealthCheck()
		time.Sleep(30 * time.Second)
	}()
}
