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
	failservicename := []string{}
	successservicename := []string{}
	microservicelist := config.Conf.MicroServiceList
	for _, microservice := range microservicelist {
		if microservice.EnableHealth {
			resp, err := http.Get(fmt.Sprintf("http://localhost:%d/api/v1/ping", microservice.Port))
			if err != nil || resp.StatusCode != 200 {
				failservicename = append(failservicename, microservice.Name)
			} else {
				successservicename = append(successservicename, microservice.Name)
			}
			defer resp.Body.Close()
		}
	}
	common.Log.Debugf("%v health check success", successservicename)
	if failservicename != nil {
		common.Log.Errorf("%v health check failed", failservicename)
		return fmt.Errorf("%v HealthCheck failed", failservicename)
	}
	return nil
}

func CtrlLED() {
	microservice := config.Conf.ControlLed
	ledStatus := constLedOn
	for _, wwwurl := range microservice.URLForWWWLed {
		err := CheckURL(wwwurl)
		if err != nil {
			ledStatus = constLedOff
		}
	}
	_ = setLed(constLedWWW, ledStatus) // ignore return error

	ledStatus = constLedOn
	for _, linkurl := range microservice.URLForLinkLed {
		err := CheckURL(linkurl)
		if err != nil {
			ledStatus = constLedOff
		}
	}
	_ = setLed(constLedLink, ledStatus) // ignore return error
}

// ScheduledHealthCheck 定时轮询任务
func ScheduledHealthCheck() {
	go func() {
		CtrlLED()
		err := HealthCheck()
		if err != nil {
			_ = setLed(constLedStatus, constLedFlash) // ignore return error
			//common.Log.Error(err)
		} else {
			_ = setLed(constLedStatus, constLedOn) // ignore return error
		}
		time.Sleep(30 * time.Second)
	}()
}
