package services

import (
	"fmt"
	"net/http"
	"sysmgmt-next/config"
	"time"

	"github.com/conthing/utils/common"
)

// HealthCheck 健康检查
func HealthCheck() error {
	servicename := []string{} //todo 这里变量的命名也不合适
	microservicelist := config.Conf.MicroServiceList
	for _, microservice := range microservicelist {
		if microservice.EnableHealth == true {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/v1/ping", microservice.Port))
			if err != nil || resp.StatusCode != 200 {
				// todo again 这里日志有问题，返回Error在哪里？
				servicename = append(servicename, microservice.Name)
				common.Log.Error("%v health check failed", servicename)
			}
			defer resp.Body.Close()
		} else {
			common.Log.Info("%s health check success", microservice.Name)
			return nil // todo again 这里的return 也是不能要的 日志内容也是不对的
		}
	}
	return nil
}

// todo again 函数命名不合适，return 乱七八糟，日志位置也不对
func LedStatus() error {
	microservice := config.Conf.ControlLed
	for _, wwwurl := range microservice.URLForWWWLed {
		err := CheckURL(wwwurl)
		if err != nil {
			setLed(constLedWWW, constLedOff)
			common.Log.Error("WWW Led is off")
			return err
		} else {
			setLed(constLedWWW, constLedOn)
			common.Log.Info("WWW Led is on")
			return nil
		}
	}
	for _, linkurl := range microservice.URLForLinkLed {
		err := CheckURL(linkurl)
		if err != nil {
			setLed(constLedLink, constLedOff)
			common.Log.Error("Link Led is off")
			return err
		} else {
			setLed(constLedLink, constLedOn)
			common.Log.Info("Link Led is on")
			return nil
		}
	}
	return nil
}

// ScheduledHealthCheck 定时轮询任务
func ScheduledHealthCheck() {
	go func() {
		LedStatus()
		HealthCheck()
		time.Sleep(30 * time.Second)
	}()
}
