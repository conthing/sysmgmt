package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/conthing/sysmgmt/dto"

	"github.com/conthing/utils/common"
)

// HealthCheckConfig 配置文件结构
type HealthCheckConfig struct {
	ControlLed       StLedControl
	MicroServiceList []MicroService
	Recovery         StRecovery
}

type StRecovery struct {
	Contains    string
	Command     string
	Parameter   []string
	Environment []string
	OutputFile  string
}

// 健康检查有失败的，status灯就闪烁，全部健康则常亮
// WWW和Link灯，每个灯对应一个URL列表。对每个URL的GET返回均正常，指示灯常亮，任何一个URL返回不正常，指示灯常灭
// URL的GET返回正常是指：HTTP返回码等于200，且body里不包含以下字符串的任意一个"err, fail, disconnect, timeout"

type StLedControl struct {
	URLForWWWLed  []string
	URLForLinkLed []string
}

// MicroService 微服务配置
type MicroService struct {
	Name         string
	Port         int
	EnableHealth bool
}

var healthCheckConfig HealthCheckConfig

func HealthCheckInit(ControlLed StLedControl, MicroServiceList []MicroService, Recovery StRecovery) {
	healthCheckConfig.ControlLed = ControlLed
	healthCheckConfig.MicroServiceList = MicroServiceList
	healthCheckConfig.Recovery = Recovery
}

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
	microservicelist := healthCheckConfig.MicroServiceList
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
	microservicelist := healthCheckConfig.MicroServiceList
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
	microservice := healthCheckConfig.ControlLed
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
	if match, err := regexp.MatchString(healthCheckConfig.Recovery.Contains, str); err == nil {
		if !match {
			common.Log.Errorf("recovery check: %q is not exist, restart...", healthCheckConfig.Recovery.Contains)
			go Restart()
		} else {
			common.Log.Debugf("recovery check: %q is exist", healthCheckConfig.Recovery.Contains)
		}
	} else {
		common.Log.Errorf("regexp MatchString %q failed: %v", healthCheckConfig.Recovery.Contains, err)
	}
}

func Restart() {
	common.Log.Debugf("exec %v para:%v env:%v > %v", healthCheckConfig.Recovery.Command, healthCheckConfig.Recovery.Parameter, healthCheckConfig.Recovery.Environment, healthCheckConfig.Recovery.OutputFile)
	command := exec.Command(healthCheckConfig.Recovery.Command, healthCheckConfig.Recovery.Parameter...)              //"/app/zap/lpr/lpr", "-d", "/app/log/lpr", "-c", "/app/zap/lpr"
	command.Env = append(os.Environ(), healthCheckConfig.Recovery.Environment...)                                     //LD_LIBRARY_PATH=/app/zap/lpr
	f, err := os.OpenFile(healthCheckConfig.Recovery.OutputFile, os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0755) // "/app/log/conthing-lpr.log"
	if err != nil {
		common.Log.Errorf("open %s failed: %v", healthCheckConfig.Recovery.OutputFile, err)
	} else {
		command.Stderr = f
		command.Stdout = f
	}

	err = command.Run()
	if err != nil {
		common.Log.Errorf("%s restart failed: %v", healthCheckConfig.Recovery.Contains, err)
		return
	}
}

// 0-其他触发led变化的事件
// 1-恢复出厂按钮按下 2-恢复出厂按钮按下后10秒 3-恢复出厂按钮松开
// 11-function按钮按下 12-function按钮按下后5秒 13-function按钮松开
// 23-identify
var buttonEventChannel = make(chan int)

func NotifyLed() {
	buttonEventChannel <- 0 // 指示灯要变化
}
func IdentifyLed() {
	buttonEventChannel <- 23 // 指示灯要变化
}

func RebootLater() {
	go func() {
		time.Sleep(time.Second * 3)
		_, err := exec.Command("reboot", "-f").Output() //初始化Cmd
		if err != nil {
			common.Log.Errorf("reboot failed: %v", err)
		}
	}()
}

// ScheduledHealthCheck 定时轮询任务
func ScheduledHealthCheck() error {
	bypass := false
	for {
		select {
		case <-time.After(time.Second * 30):
		case rst := <-buttonEventChannel:
			if rst == 0 { // 其他触发led变化的事件
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
				_, _ = exec.Command("rm", "-rf", "/app/data").Output() //复位
				RebootLater()                                          //复位
			} else if rst == 3 { // 恢复出厂按钮松开
				bypass = false
			} else if rst == 11 { // function按钮按下
				bypass = true
				_ = setLed(constLedWWW, constLedOff)
			} else if rst == 12 { // function按钮按下后5秒
				bypass = false
				_ = setLed(constLedWWW, constLedFlash)
				_ = SetNetInfo(&dto.NetInfo{DHCP: true})
				time.Sleep(1 * time.Second)
			} else if rst == 13 { // function按钮松开
				bypass = false
			} else if rst == 23 { // identify事件
				bypass = false
				_ = setLed(constLedStatus, constLedFlash)
				_ = setLed(constLedLink, constLedFlash)
				_ = setLed(constLedWWW, constLedFlash)
				continue
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
				_ = setLed(constLedStatus, constLedBlink) // ignore return error
			}
			Recovery()
		} else {
			common.Log.Debug("health check bypassed...")
		}

	}
}

// WatchDog 看门狗 todo 这里没有真正的看门狗作用，需要放到healthcheck里面
func WatchDog() error {
	wdt, err := GetWatchDog(10) //10s超时
	if err != nil {
		return fmt.Errorf("watchdog init failed: %w", err)
	}
	for {
		time.Sleep(time.Second * 4)

		err = KeepAlive(wdt) //10s超时
		if err != nil {
			return fmt.Errorf("feed dog failed: %w", err)
		} else {
			common.Log.Debug("feed dog ok")
		}
	}
}
