package services

import (
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/conthing/utils/common"
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

// todo review 这里不对，把我给的函数原型怎么都改了
// 原型恢复后，达到这样的效果：如果调用方式是setLed(constLedStatus,constLedFlash)，函数里就会执行/usr/test/led-pwm-start /dev/led-pwm1 ...
// 所以setLed的函数体里面，是根据入参“选择”exec不通的内容，并判断返回是否正常
// setLed 设置led的开关闪状态
// todo again 领会了我的意思了，但两点改进，1.三个灯都要支持亮灭闪，2..改变if判断的瞬息可以将结构优化一下
func setLed(led string, status byte) error {
	if led == "constLedLink" { //todo again 这里为什么加引号？？？
		if status == constLedOff {
			exec.Command("/usr/test/led-hrtimer-close", constLedLink).Output()
		} else if status == constLedOn {
			exec.Command("/usr/test/led-pwm-start", constLedLink, "200000000", "200000000").Output()
		}
	} else if led == "constLedWWW" {
		if status == constLedOff {
			exec.Command("/usr/test/led-hrtimer-close", constLedWWW).Output()
		} else if status == constLedOn {
			exec.Command("/usr/test/led-pwm-start", constLedWWW, "200000000", "200000000").Output()
		}
	} else if led == "constLedStatus" {
		if status == constLedFlash {
			exec.Command("/usr/test/led-pwm-start", constLedStatus, "200000000", "100000000").Output()
		} else if status == constLedOn {
			exec.Command("/usr/test/led-pwm-start", constLedStatus, "200000000", "200000000").Output()
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
	// todo again 这个 if else 关系还是没有仔细考虑，还有日志不要用中文
	if resp.StatusCode == 200 && strings.Contains(str, "err") == false && strings.Contains(str, "fail") == false && strings.Contains(str, "disconnect") == false && strings.Contains(str, "timeout") == false {
		common.Log.Info("运行正常!")
	} else if strings.Contains(str, "err") == true || strings.Contains(str, "fail") == true || strings.Contains(str, "disconnect") == true || strings.Contains(str, "timeout") == true {
		common.Log.Error("运行不正常: ", str)
		return err
	}
	// todo review 这个 if和上面的if的关系是什么？这么写是错的
	//if strings.Contains(str, "err") == true || strings.Contains(str, "fail") == true || strings.Contains(str, "disconnect") == true || strings.Contains(str, "timeout") == true {
	//	common.Log.Error("运行不正常: ", str)
	//	return err
	//}
	return nil
}
