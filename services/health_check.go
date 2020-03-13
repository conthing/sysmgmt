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
	failservicename := []string{} //todo 这里变量的命名也不合适
	successservicename := []string{}
	microservicelist := config.Conf.MicroServiceList
	for _, microservice := range microservicelist {
		if microservice.EnableHealth == true {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/v1/ping", microservice.Port))
			if err != nil || resp.StatusCode != 200 {
				// todo again 这里日志有问题，返回Error在哪里？
				failservicename = append(failservicename, microservice.Name)
			} else {
				successservicename = append(successservicename, microservice.Name)
			}
			defer resp.Body.Close()
		}
	}
	common.Log.Debugf("%v health check success", successservicename)
	if failservicename != nil {
		common.Log.Error("%v health check failed", failservicename)
		return fmt.Errorf("%v HealthCheck failed", failservicename)
	}
	return nil
}

// todo again 函数命名不合适，return 乱七八糟，日志位置也不对
func CtrlLED() error {
	microservice := config.Conf.ControlLed
	ledStatus := constLedOn
	for _, wwwurl := range microservice.URLForWWWLed {
		err := CheckURL(wwwurl)
		if err != nil {
			ledStatus = constLedOff
		}
	}
	setLed(constLedWWW, ledStatus)

	ledStatus = constLedOn
	for _, linkurl := range microservice.URLForLinkLed {
		err := CheckURL(linkurl)
		if err != nil {
			ledStatus = constLedOff
		}
	}
	setLed(constLedLink, ledStatus)
	return nil
}

// ScheduledHealthCheck 定时轮询任务
func ScheduledHealthCheck() {
	go func() {
		CtrlLED()
		err := HealthCheck()
		if err != nil {
			setLed(constLedStatus, constLedFlash)
			//common.Log.Error(err)
		} else {
			setLed(constLedStatus, constLedOn)
		}
		time.Sleep(30 * time.Second)
	}()
}
