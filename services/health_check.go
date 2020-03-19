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
	failservicename := ""
	successservicename := ""
	microservicelist := config.Conf.MicroServiceList
	for _, microservice := range microservicelist {
		common.Log.Debugf("ms: %v", microservice)
		if microservice.EnableHealth {
			url := fmt.Sprintf("http://localhost:%d/api/v1/ping", microservice.Port)
			resp, err := http.Get(url)
			if err != nil || resp.StatusCode != 200 {
				common.Log.Errorf("Get %s failed", url)
				failservicename += microservice.Name + ","
			} else {
				common.Log.Debugf("Get %s success", url)
				successservicename += microservice.Name + ","
			}
			if err == nil {
				defer resp.Body.Close()
			}
		}
	}
	if successservicename != "" {
		common.Log.Debugf("%v health check success", successservicename)
	}
	if failservicename != "" {
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
			common.Log.Errorf("CheckURL %s failed: %v", wwwurl, err)
			ledStatus = constLedOff
		} else {
			common.Log.Debugf("CheckURL %s pass", wwwurl)
		}
	}
	_ = setLed(constLedWWW, ledStatus) // ignore return error

	ledStatus = constLedOn
	for _, linkurl := range microservice.URLForLinkLed {
		err := CheckURL(linkurl)
		if err != nil {
			common.Log.Errorf("CheckURL %s failed: %v", linkurl, err)
			ledStatus = constLedOff
		} else {
			common.Log.Debugf("CheckURL %s pass", linkurl)
		}
	}
	_ = setLed(constLedLink, ledStatus) // ignore return error
}

// ScheduledHealthCheck 定时轮询任务
func ScheduledHealthCheck() {
	go func() {
		for {
			CtrlLED()
			err := HealthCheck()
			if err != nil {
				_ = setLed(constLedStatus, constLedFlash) // ignore return error
				//common.Log.Error(err)
			} else {
				_ = setLed(constLedStatus, constLedOn) // ignore return error
			}
			time.Sleep(30 * time.Second)
		}
	}()
}
