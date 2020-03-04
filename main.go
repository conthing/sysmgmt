package main

import (
	"sysmgmt-next/config"
	"sysmgmt-next/router"
	"sysmgmt-next/services"
)

func main() {
	config.Service()
	services.MDNS(config.Conf.MDNS)

	// services.WatchDog()
	go services.ScheduledHealthCheck()
	go services.ScheduledLED(config.Conf.MicroServiceList)
	router.Service(config.Conf)

	defer services.StopMDNS()

}
