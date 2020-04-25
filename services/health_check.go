package services

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
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

func Recovery() {
	command := exec.Command("ps", "-a")
	out, err := command.Output()
	if err != nil {
		common.Log.Errorf("exec ps failed: %v", err)
		return
	}
	str := string(out)
	if !strings.Contains(str, config.Conf.Recovery.Contains) { //"lpr -d"
		common.Log.Errorf("recovery check: %s is not exist, restart...", config.Conf.Recovery.Contains)
		go Restart()
	} else {
		common.Log.Debugf("recovery check: %s is exist", config.Conf.Recovery.Contains)
	}
}

func Restart() {
	common.Log.Debugf("exec %v para:%v env:%v > %v", config.Conf.Recovery.Command, config.Conf.Recovery.Parameter, config.Conf.Recovery.Environment, config.Conf.Recovery.OutputFile)
	command := exec.Command(config.Conf.Recovery.Command, config.Conf.Recovery.Parameter...)                    //"/app/zap/lpr/lpr", "-d", "/app/log/lpr", "-c", "/app/zap/lpr"
	command.Env = append(os.Environ(), config.Conf.Recovery.Environment...)                                     //LD_LIBRARY_PATH=/app/zap/lpr
	f, err := os.OpenFile(config.Conf.Recovery.OutputFile, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0755) // "/app/log/conthing-lpr.log"
	if err != nil {
		common.Log.Errorf("open %s failed: %v", config.Conf.Recovery.OutputFile, err)
	} else {
		command.Stderr = f
		command.Stdout = f
	}

	err = command.Run()
	if err != nil {
		common.Log.Errorf("%s restart failed: %v", config.Conf.Recovery.Contains, err)
		return
	}
}

// ScheduledHealthCheck 定时轮询任务
func ScheduledHealthCheck() {
	go func() {
		for {
			common.Log.Debug("health check...")
			CtrlLED()
			err := HealthCheck()
			if err != nil {
				_ = setLed(constLedStatus, constLedFlash) // ignore return error
				//common.Log.Error(err)
			} else {
				_ = setLed(constLedStatus, constLedOn) // ignore return error
			}
			Recovery()

			time.Sleep(30 * time.Second)
		}
	}()
}
