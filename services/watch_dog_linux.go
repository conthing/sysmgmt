package services

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/conthing/utils/common"
)

const (
	WDIOC_KEEPALIVE  = 0x80045705
	WDIOC_SETTIMEOUT = 0xc0045706
)

// GetWatchDog 获取看门狗
func GetWatchDog(timeout int) (*os.File, error) {
	file, err := os.OpenFile("/dev/watchdog", syscall.O_RDWR, 0600)
	if err != nil {
		return nil, fmt.Errorf("failed to open watchdog: %v", err)
	}

	r, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(file.Fd()),
		uintptr(WDIOC_SETTIMEOUT),
		uintptr(unsafe.Pointer(&timeout)))

	if errno != 0 {
		return nil, os.NewSyscallError("SYS_IOCTL", errno)
	}

	if r != 0 {
		return nil, fmt.Errorf("unknown error from SYS_IOCTL")
	}
	common.Log.Warnf("Watchdog init %d sec timeout", timeout)
	return file, nil
}

// KeepAlive 保持
func KeepAlive(file *os.File) error {
	r, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(file.Fd()),
		uintptr(WDIOC_KEEPALIVE),
		uintptr(0))
	if errno != 0 {
		return os.NewSyscallError("SYS_IOCTL", errno)
	}

	if r != 0 {
		return fmt.Errorf("unknown error from SYS_IOCTL")
	}
	common.Log.Debug("feed dog")
	return nil
}
