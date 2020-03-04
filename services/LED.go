package services

import (
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"sysmgmt-next/config"
	"time"

	"github.com/conthing/utils/common"
	jsoniter "github.com/json-iterator/go"
)

// todo 删除掉前3个函数
// ScheduledLED 指示灯
func ScheduledLED(list []config.MicroService) {
	for {
		for _, s := range list {
			controlWWWAndLink(s)
		}
		controlHealthStatusLED()

		time.Sleep(time.Second * 30)
	}
}

func controlHealthStatusLED() {
	if IsHealth {
		exec.Command("/usr/test/led-pwm-start-percentage", "/dev/led-pwm1", "2", "1").Output()
	}
}

// controlWWWAndLink 控制指示灯
func controlWWWAndLink(s config.MicroService) {
	resp, err := http.Get(s.URL)
	// http 连接异常
	if err != nil {
		str, _ := exec.Command("/usr/test/led-hrtimer-close", s.LED).Output()
		common.Log.Error("Get url 出错", err, "关灯 -> ", string(str))
		return
	}
	// body 分类
	data, err := ioutil.ReadAll(resp.Body)

	// body 异常
	if err != nil {
		str, _ := exec.Command("/usr/test/led-hrtimer-close", s.LED).Output()
		common.Log.Error("解析 body 出错", err, "关灯 -> ", string(str))
		return
	}
	// status
	switch s.Type {
	case "status":
		if jsoniter.Get(data, "status").ToString() == "connected" {
			exec.Command("/usr/test/led-pwm-start-percentage", s.LED, "2", "1").Output()
		} else {
			exec.Command("/usr/test/led-hrtimer-close", s.LED).Output()
			common.Log.Error("Get not connected", "关灯")
		}
	default:
		log.Fatal("配置文件错误，类型仅支持 ping 或 status")
	}
}

// 1-status 2-www 3-link
const (
	constLedStatus string = "/dev/led-pwm1"
	//todo 其他led的定义

	constLedOff   byte = byte(0)
	constLedOn    byte = byte(1)
	constLedFlash byte = byte(2)
)

// setLed 设置led的开关闪状态
func setLed(led string, status byte) error {
	//todo 完成此函数
	return nil
}

// GET此url，在HTTP返回码等于200，且body里不包含以下字符串的任意一个"err, fail, disconnect, timeout"时无错误，否则返回error
func CheckURL(url string) error {
	// todo 完成此函数
	return nil
}
