package main

import (
	"log"
	"sysmgmt-next/config"
	"sysmgmt-next/router"
	"sysmgmt-next/services"
	"time"
)

func main() {
	config.Service()
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
						log.Fatalf("feed dog failed: %v", err)
					} else {
						log.Println("feed dog ok")
					}
				}

			}

		} else {
			log.Fatal("watchdog init failed: ", err)
		}
	}()
}
