package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/conthing/sysmgmt/config"
	"github.com/conthing/sysmgmt/dto"

	"github.com/conthing/utils/common"
)

func GetAllVersion() (globalVersion dto.VersionInfo) {
	command := exec.Command("cat", "../VERSION") //初始化Cmd
	out, err := command.Output()
	if err != nil {
		common.Log.Errorf("open ../VERSION failed %v", err)
		out = []byte("")
	}
	text := strings.SplitAfterN(string(out), "\n", 2) //用第一个 \n 分割字符串
	globalVersion.Version = strings.TrimSpace(text[0])
	if len(text) > 1 {
		globalVersion.Description = strings.TrimSpace(text[1])
	}
	globalVersion.SubVersion = append(globalVersion.SubVersion, dto.SubVersionInfo{Name: "sysmgmt", Version: common.Version, BuildTime: common.BuildTime})

	var version dto.SubVersionInfo
	microservicelist := config.Conf.MicroServiceList
	for _, microservice := range microservicelist {
		url := "http://localhost:" + strconv.FormatInt(int64(microservice.Port), 10) + "/api/v1/version"
		resp, err := http.Get(url)
		if err == nil {
			defer resp.Body.Close()
		} else {
			common.Log.Errorf("%s Get failed: %v", url, err)
			continue
		}
		if resp.StatusCode != 200 {
			common.Log.Errorf("%s Get failed: code:%d", url, resp.StatusCode)
			continue
		}
		body, _ := ioutil.ReadAll(resp.Body)
		str := string(body)
		common.Log.Debugf("%s Get: %s", url, str)
		strArry := strings.SplitAfterN(str, " ", 2)
		version.Name = microservice.Name
		if len(strArry) > 1 {
			version.Version = strArry[0]
			version.BuildTime = strArry[1]
			globalVersion.SubVersion = append(globalVersion.SubVersion, version)
		} else {
			version.Version = "unknown"
			version.BuildTime = "unknown"
			globalVersion.SubVersion = append(globalVersion.SubVersion, version)
		}
	}
	return
}

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
	command := exec.Command("ps", "-ef")
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
