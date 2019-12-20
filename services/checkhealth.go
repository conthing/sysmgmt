package services

import (
	"github.com/conthing/utils/common"
	"github.com/json-iterator/go"
	"io/ioutil"
	"log"
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
		str, _ := exec.Command("/usr/test/led-hrtimer-close", s.LED).Output()
		common.Log.Error("Get url 出错", err, "关灯 -> ", string(str))
		return
	}
	// body 分类
	data, err := ioutil.ReadAll(resp.Body)
	common.Log.Info(string(data))
	// body 异常
	if err != nil {
		str, _ := exec.Command("/usr/test/led-hrtimer-close", s.LED).Output()
		common.Log.Error("解析 body 出错", err, "关灯 -> ", string(str))
		return
	}
	// ping
	switch s.Type {
	case "ping":
		if string(data) == "pong" {
			exec.Command("/usr/test/led-pwm-start-percentage", s.LED, "2", "1").Output()
			common.Log.Info("Get URL pong", "开灯")
		} else {
			exec.Command("/usr/test/led-hrtimer-close", s.LED).Output()
			common.Log.Error("Get not pong", "关灯")
		}
	case "status":
		if jsoniter.Get(data, "status").ToString() == "connected" {
			exec.Command("/usr/test/led-pwm-start-percentage", s.LED, "2", "1").Output()
			common.Log.Info("Get URL connected", "开灯")
		} else {
			exec.Command("/usr/test/led-hrtimer-close", s.LED).Output()
			common.Log.Error("Get not connected", "关灯")
		}
	default:
		log.Fatal("配置文件错误，类型仅支持 ping 或 status")
	}
}
