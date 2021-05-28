package services

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

// 1-status 2-www 3-link
const (
	constLedCommand string = "/usr/test/led-pwm-test"
	constLedStatus  string = "/dev/led-pwm1"
	constLedWWW     string = "/dev/led-pwm2"
	constLedLink    string = "/dev/led-pwm3"
	constLedOff     byte   = byte(0)
	constLedOn      byte   = byte(1)
	constLedFlash   byte   = byte(2)
)

// todo 网关identify，整机组装支持
// 原型恢复后，达到这样的效果：如果调用方式是setLed(constLedStatus,constLedFlash)，函数里就会执行/usr/test/led-pwm-start /dev/led-pwm1 ...
// 所以setLed的函数体里面，是根据入参“选择”exec不通的内容，并判断返回是否正常
// setLed 设置led的开关闪状态
func setLed(led string, status byte) error {
	var cmd *exec.Cmd
	if status == constLedOff {
		cmd = exec.Command(constLedCommand, "0", led)
	} else if status == constLedOn {
		cmd = exec.Command(constLedCommand, "1", led, "200000000", "199999999")
	} else if status == constLedFlash {
		cmd = exec.Command(constLedCommand, "1", led, "200000000", "100000000")
	}

	if cmd == nil {
		return fmt.Errorf("LED operation not supported")
	}
	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}

// resp的Status是返回码，body里才是字符串，字符串判断的依据是包含，不是等于，参考原来的http.Get的地方怎么写的
// GET此url，在HTTP返回码等于200，且body里不包含以下字符串的任意一个"err, fail, disconnect, timeout"时无错误，否则返回error
func CheckURL(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("CheckURL failed: %v", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	str := string(body)
	if resp.StatusCode != 200 || strings.Contains(str, "err") || strings.Contains(str, "fail") || strings.Contains(str, "disconnect") || strings.Contains(str, "timeout") {
		return fmt.Errorf("%s response failed: code:%d, body%v", url, resp.StatusCode, str)
	}
	return nil
}
