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
		common.Log.Debugf("%v port :%d", microservice.Name, microservice.Port)
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
	if len(microservice.URLForWWWLed) <= 0 { // 如果未定义WWW，则此指示灯用来指示DHCP开关
		var info dto.NetInfo
		if err := GetCurrentNetInfo(&info); err == nil {
			if !info.DHCP {
				ledStatus = constLedOff
			}
		}
	} else {
		for _, wwwurl := range microservice.URLForWWWLed {
			err := CheckURL(wwwurl)
			if err != nil {
				common.Log.Errorf("CheckURL %s failed: %v", wwwurl, err)
				ledStatus = constLedOff
			} else {
				common.Log.Debugf("CheckURL %s pass", wwwurl)
			}
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
	if !strings.Contains(str, config.Conf.Recovery.Contains) { // todo eroom同时也是目录名，这里判断不出
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

var buttonEventChannel = make(chan int) // 0-恢复出厂按钮松开 1-恢复出厂按钮按下 2-恢复出厂按钮按下后10秒 10-function按钮松开 11-function按钮按下 12-function按钮按下后5秒

// ScheduledHealthCheck 定时轮询任务
func ScheduledHealthCheck() {
	go func() {
		bypass := false
		for {

			select {
			case <-time.After(time.Second * 30):
			case rst := <-buttonEventChannel:
				if rst == 0 { // 恢复出厂按钮松开
					bypass = false
				} else if rst == 1 { // 恢复出厂按钮按下
					bypass = true
					_ = setLed(constLedStatus, constLedOff)
					_ = setLed(constLedLink, constLedOff)
					_ = setLed(constLedWWW, constLedOff)
				} else if rst == 2 { // 恢复出厂按钮按下后10秒
					bypass = true
					_ = setLed(constLedStatus, constLedFlash)
					_ = setLed(constLedLink, constLedFlash)
					_ = setLed(constLedWWW, constLedFlash)
					_ = SetNetInfo(&dto.NetInfo{DHCP: false, Address: "192.168.0.101", Netmask: "255.255.255.0", Gateway: "192.168.0.1"})
					exec.Command("rm", "-rf", "/app/data").Output() //复位
					time.Sleep(3 * time.Second)
					exec.Command("reboot").Output() //复位
				} else if rst == 10 { // function按钮松开
					bypass = false
				} else if rst == 11 { // function按钮按下
					bypass = true
					_ = setLed(constLedWWW, constLedOff)
				} else if rst == 12 { // function按钮按下后5秒
					bypass = true
					_ = setLed(constLedWWW, constLedFlash)
					_ = SetNetInfo(&dto.NetInfo{DHCP: true})
					time.Sleep(1 * time.Second)
				} else {
					bypass = false
				}
			}

			if !bypass {
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
			} else {
				common.Log.Debug("health check bypassed...")
			}

		}
	}()
}
