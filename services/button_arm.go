package services

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"syscall"
	"time"
	"unsafe"

	"github.com/conthing/utils/common"
)

// 1-status 2-www 3-link
const (
	constButtonEvent        string = "/dev/input/event2"
	constResetButtonCode    uint16 = uint16(0x100)
	constFunctionButtonCode uint16 = uint16(0x99)
	constKeyEventType       uint16 = uint16(1)
)

var eventsize = int(unsafe.Sizeof(InputEvent{}))

// 恢复出厂按钮按下时的定时器
var holdTimer *time.Timer

//input文件中存储的数据结构
type InputEvent struct {
	Time  syscall.Timeval // time in seconds since epoch at which event occurred
	Type  uint16          // event type - one of ecodes.EV_*
	Code  uint16          // event code related to the event type
	Value int32           // event value related to the event type
}

// resp的Status是返回码，body里才是字符串，字符串判断的依据是包含，不是等于，参考原来的http.Get的地方怎么写的
// GET此url，在HTTP返回码等于200，且body里不包含以下字符串的任意一个"err, fail, disconnect, timeout"时无错误，否则返回error
func ButtonSevcie() error {
	f, err := os.Open(constButtonEvent)
	if err != nil {
		return fmt.Errorf("open button event file failed: %w", err)
	}

	//把file读取到缓冲区中
	defer f.Close()

	for {
		event := InputEvent{}
		buffer := make([]byte, eventsize)

		_, err := f.Read(buffer)
		if err != nil {
			return fmt.Errorf("read button event file failed: %w", err)
		}

		b := bytes.NewBuffer(buffer)
		err = binary.Read(b, binary.LittleEndian, &event)
		if err != nil {
			return fmt.Errorf("convert to event struct failed: %w", err)
		}

		ButtonProcess(&event)
	}

}

// 开始长按，状态灯熄灭，长按达到10秒时状态灯闪烁并执行恢复出厂设置，此时可以松开按键，执行完成后重启设备

func ButtonProcess(event *InputEvent) {
	if event.Type == constKeyEventType {
		common.Log.Debugf("key event: %d code:0x%04x, value:%d", event.Type, event.Code, event.Value)
		if event.Code == constResetButtonCode {
			if event.Value == 0 { // 松开
				if holdTimer != nil {
					holdTimer.Stop()
					holdTimer = nil
					common.Log.Debugf("holdtime canceled")
				}

				resetToDefaultEventChannel <- 0
			} else { // 按下
				resetToDefaultEventChannel <- 1
				holdTimer = time.AfterFunc(10*time.Second, func() {
					resetToDefaultEventChannel <- 2
				})
			}
		}

	}
}
