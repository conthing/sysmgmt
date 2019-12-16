package main

import (
	"sysmgmt-next/config"
	"sysmgmt-next/router"
	"sysmgmt-next/services"
)

func main() {
	config.Service()
	services.MDNS(config.Conf.MDNS)
	router.Service(config.Conf)
}
