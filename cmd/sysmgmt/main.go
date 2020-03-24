package main

import (
	"sysmgmt-next/config"
	"sysmgmt-next/redis"
	"sysmgmt-next/router"
	"sysmgmt-next/services"
	"time"

	"github.com/conthing/utils/common"
)

func main() {
	err := config.Service()
	if err != nil {
		common.Log.Errorf("load config failed %v", err)
	}
	common.Log.Infof("config: %v", config.Conf)

	redis.Connect()
	services.MDNS(config.Conf.MDNS)

	// WatchDog()
	services.ScheduledHealthCheck()
	//go services.ScheduledLED(config.Conf.MicroServiceList)
	router.Service(*config.Conf)

	defer services.StopMDNS()

}

// WatchDog 看门狗
func WatchDog() {
	go func() {
		wdt, err := services.GetWatchDog(10) //10s超时
		if err == nil {
			for {
				select {
				case <-time.After(time.Second * 4):
					err = services.KeepAlive(wdt) //10s超时
					if err != nil {
						common.Log.Errorf("feed dog failed: %v", err)
					} else {
						common.Log.Debug("feed dog ok")
					}
				}

			}

		} else {
			common.Log.Errorf("watchdog init failed: ", err)
		}
	}()
}
