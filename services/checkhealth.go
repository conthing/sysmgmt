package services

import (
	"github.com/conthing/utils/common"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"os/exec"
	"sysmgmt-next/config"
	"time"
)

func CheckServiceHealth(lightlist []config.MicroService) {
	for {
		for _, s := range lightlist {
			LED(s)
		}
		time.Sleep(time.Second * 30)
	}
}

func LED(s config.MicroService) {
	common.Log.Info(s.URL)
	resp, err := http.Get(s.URL)
	// http 连接异常
	if err != nil {
		str, _ := exec.Command("/usr/test/led-hrtimer-close", s.Name).Output()
		common.Log.Error("Get url 出错", err, "关灯 -> ", string(str))
	}
	// body 分类
	data, err := ioutil.ReadAll(resp.Body)
	common.Log.Info(string(data))
	// body 异常
	if err != nil {
		str, _ := exec.Command("/usr/test/led-hrtimer-close", s.Name).Output()
		common.Log.Error("解析 body 出错", err, "关灯 -> ", string(str))
	}
	// ping
	if s.Type == "ping" && string(data) == "pong" {
		str, _ := exec.Command("/usr/test/led-pwm-start-percentage", s.Name, "2", "1").Output()
		common.Log.Info("Get url pong", "开灯 -> ", string(str))
	}
	// status
	if s.Type == "status" && jsoniter.Get(data, "status").ToString() == "connected" {
		str, _ := exec.Command("/usr/test/led-pwm-start-percentage", s.Name, "2", "1").Output()
		common.Log.Info("Get url status connected", "开灯 -> ", string(str))

	} else {
		str, _ := exec.Command("/usr/test/led-hrtimer-close", s.Name).Output()
		common.Log.Error("Get url", err, "关灯 -> ", string(str))
	}
}
