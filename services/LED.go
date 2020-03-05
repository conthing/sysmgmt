package services

import (
	"net/http"

	"github.com/conthing/utils/common"
)

// 1-status 2-www 3-link
const (
	constLedStatus string = "/dev/led-pwm1"
	constLedWWW    string = "/dev/led-pwm2"
	constLedLink   string = "/dev/led-pwm3"

	constLedOff   byte = byte(0)
	constLedOn    byte = byte(1)
	constLedFlash byte = byte(2)
)

// setLed 设置led的开关闪状态
func setLed(led string, status byte) error {
	//亮Link灯
	//这个函数写起来没什么思路
	return nil
}

// resp的Status是返回码，body里才是字符串，字符串判断的依据是包含，不是等于，参考原来的http.Get的地方怎么写的
// GET此url，在HTTP返回码等于200，且body里不包含以下字符串的任意一个"err, fail, disconnect, timeout"时无错误，否则返回error
func CheckURL(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		if resp.Status != "err" && resp.Status != "fail" && resp.Status != "disconnect" && resp.Status != "timeout" && resp.StatusCode == 200 {
			return nil
		} else {
			common.Log.Error(err)
			return err
		}
	}
	defer resp.Body.Close()
	return nil
}
