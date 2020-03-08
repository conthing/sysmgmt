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
	count := 0
	stringname1 := ""
	stringname2 := ""
	microservicelist := config.Conf.MicroServiceList
	for _, microservice := range microservicelist {
		if microservice.EnableHealth == true {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/v1/ping", microservice.Port))
			if err != nil || resp.StatusCode != 200 {
				count++
				if count == 1 {
					stringname1 = microservice.Name
					//common.Log.Error("微服务不健康的有: ", microservice.Name)
				} else if count == 2 {
					stringname2 = microservice.Name
					//common.Log.Error("微服务不健康的有: ", microservice.Name)
					common.Log.Error("%s %s health check failed", stringname1, stringname2)
					return err
				}
			}
			defer resp.Body.Close()
		} else {
			common.Log.Info("微服务运行健康的有: ", microservice.Name)
			return nil
		}
	}
	return nil
}

func ControlLed() error {
	microservicelist := config.Conf.MicroServiceList
	for _, microservice := range microservicelist {
		if microservice.Type == "mesh" { //todo 问题二:485总线正常不知怎么确定  当时说LINK灯一方面获取到正确的mesh或485总线正常则常亮，否则灭掉
			err := CheckURL("http://localhost:" + string(microservice.Port) + "/api/v1/mesh")
			if err != nil {
				setLed(constLedLink, constLedOff)
				common.Log.Error("http://localhost:port/api/v1/mesh: ", err)
				return err
			} else {
				setLed(constLedLink, constLedOn)
				return nil
			}
		}
		if microservice.Type == "status" {
			err := CheckURL("http://localhost:" + string(microservice.Port) + "/api/v1/status")
			if err != nil {
				setLed(constLedWWW, constLedOff)
				common.Log.Error("http://localhost:port/api/v1/status: ", err)
				return err
			} else {
				setLed(constLedWWW, constLedOn)
				return nil
			}
		}
		if microservice.Type == "ping" {
			err := CheckURL("http://localhost:" + string(microservice.Port) + "/api/v1/ping")
			if err != nil {
				setLed(constLedStatus, constLedFlash)
				common.Log.Error("http://localhost:port/api/v1/ping: ", err)
				return err
			} else {
				setLed(constLedStatus, constLedOn)
				return nil
			}
		}
	}
	return nil
}

// todo 改写此函数，将健康检查和LED控制放到同一个go程里
// ScheduledHealthCheck 定时轮询任务
func ScheduledHealthCheck() {
	go func() {
		ControlLed()
		//setLed() // todo review 因为setLed不对，所以这个写法不对；应该讲你现在实现的setLed里的一些内容搬过来
		HealthCheck()
		time.Sleep(30 * time.Second)
	}()
}
