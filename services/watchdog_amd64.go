package services

import (
	"os"

	"github.com/conthing/utils/common"
)

// GetWatchDog 获取看门狗
func GetWatchDog(timeout int) (*os.File, error) {
	common.Log.Error("Watchdog is not supported!")
	return nil, nil
}

// KeepAlive 保持
func KeepAlive(file *os.File) error {
	common.Log.Error("Watchdog is not supported!")
	return nil
}
