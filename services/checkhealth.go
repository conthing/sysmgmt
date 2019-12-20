package services

import (
	"net/http"
	"os/exec"
	"sysmgmt-next/config"
	"time"
)

func CheckServiceHealth(lightlist []config.MicroService) {
	time.Sleep(time.Second * 30)
	for _, s := range lightlist {
		LED(s.URL, s.LED)
	}
}

func LED(url string, name string) {
	resp, err := http.Get(url)
	if err != nil {
		command := exec.Command("/usr/test/led-hrtimer-close", name)
		command.Output()
	}
	if resp.StatusCode == 200 {
		command := exec.Command("/usr/test/led-pwm-start-percentage", name, "2", "1")
		command.Output()
	} else {
		command := exec.Command("/usr/test/led-hrtimer-close", name)
		command.Output()
	}
}
