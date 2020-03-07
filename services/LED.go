package services

import (
	"github.com/conthing/utils/common"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"sysmgmt-next/config"
)

// 1-status 2-www 3-link
const (
	constLedStatus string = "/dev/led-pwm1"
	constLedWWW    string = "/dev/led-pwm2"
	constLedLink   string = "/dev/led-pwm3"

	//todo 问题一:这个怎么通过逻辑转化使led灯亮灭闪，目前没想通
	constLedOff   byte = byte(0)
	constLedOn    byte = byte(1)
	constLedFlash byte = byte(2)
)

// setLed 设置led的开关闪状态
func setLed() error {
	microservicelist := config.Conf.MicroServiceList
	for _, microservice := range microservicelist {
		if microservice.Type == "mesh" { //todo 问题二:485总线正常不知怎么确定  当时说LINK灯一方面获取到正确的mesh或485总线正常则常亮，否则灭掉
			err := CheckURL("http://localhost:" + string(microservice.Port) + "/api/v1/mesh")
			if err != nil {
				exec.Command("/usr/test/led-hrtimer-close", constLedLink, string(constLedOff)).Output()
				//exec.Command("/usr/test/led-hrtimer-close", constLedLink).Output()
				return err
			} else {
				exec.Command("/usr/test/led-pwm-start", constLedLink, string(constLedOn)).Output()
				//exec.Command("/usr/test/led-pwm-start", constLedLink, 200000000, 200000000).Output()
				//或exec.Command("/usr/test/led-pwm-start-percent", constLedLink, 200000000, 100).Output()
				return nil
			}
		}
		if microservice.Type == "status" {
			err := CheckURL("http://localhost:" + string(microservice.Port) + "/api/v1/status")
			if err != nil {
				exec.Command("/usr/test/led-hrtimer-close", constLedWWW, string(constLedOff)).Output()
				//exec.Command("/usr/test/led-hrtimer-close", constLedWWW).Output()
				return err
			} else {
				exec.Command("/usr/test/led-pwm-start", constLedWWW, string(constLedOn)).Output()
				//exec.Command("/usr/test/led-pwm-start", constLedWWW, 200000000, 200000000).Output()
				//或exec.Command("/usr/test/led-pwm-start-percent", constLedWWW, 200000000, 100).Output()
				return nil
			}
		}
		if microservice.Type == "ping" {
			err := CheckURL("http://localhost:" + string(microservice.Port) + "/api/v1/ping")
			if err != nil {
				exec.Command("/usr/test/led-pwm-start", constLedStatus, string(constLedFlash)).Output()
				//exec.Command("/usr/test/led-pwm-start", constLedStatus, 200000000, 100000000).Output()
				//或exec.Command("/usr/test/led-pwm-start-percent", constLedStatus, 200000000, 50).Output()
				return err
			} else {
				exec.Command("/usr/test/led-pwm-start", constLedStatus, string(constLedOn)).Output()
				//exec.Command("/usr/test/led-pwm-start", constLedStatus, 200000000, 200000000).Output()
				//或exec.Command("/usr/test/led-pwm-start-percent", constLedStatus, 200000000, 100).Output()
				return nil
			}
		}
	}
	return nil
}

// resp的Status是返回码，body里才是字符串，字符串判断的依据是包含，不是等于，参考原来的http.Get的地方怎么写的
// GET此url，在HTTP返回码等于200，且body里不包含以下字符串的任意一个"err, fail, disconnect, timeout"时无错误，否则返回error
func CheckURL(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		common.Log.Error("解析url错误: ", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	str := string(body)
	if resp.StatusCode == 200 && strings.Contains("err", str) == false && strings.Contains("fail", str) == false && strings.Contains("disconnect", str) == false && strings.Contains("timeout", str) == false {
		common.Log.Info("运行正常!")
	}
	if strings.Contains("err", str) == true || strings.Contains("fail", str) == true || strings.Contains("disconnect", str) == true || strings.Contains("timeout", str) == true {
		common.Log.Error("运行不正常: ", str)
		return err
	}
	return nil
}
