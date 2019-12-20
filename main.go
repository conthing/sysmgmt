package main

import (
	"github.com/conthing/utils/common"
	"sysmgmt-next/config"
	"sysmgmt-next/router"
	"sysmgmt-next/services"
)

func main() {
	config.Service()
	services.MDNS(config.Conf.MDNS)

	// todo: 测试 WatchDog
	// services.WatchDog()
	common.Log.Info(config.Conf.MicroServiceList)
	go services.CheckServiceHealth(config.Conf.MicroServiceList)
	router.Service(config.Conf)

	defer services.StopMDNS()

}
