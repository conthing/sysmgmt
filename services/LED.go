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
