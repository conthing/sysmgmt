package services

import (
	"log"
	"sysmgmt-next/config"

	"github.com/grandcat/zeroconf"
)

var server *zeroconf.Server

// MDNS 服务
func MDNS(cnf config.MDNS) {
	var err error
	server, err = zeroconf.Register(cnf.Name, "_workstation._tcp", "local.", cnf.Port, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("MDNS 启动成功")
}

// StopMDNS 关闭MDNS
func StopMDNS() {
	server.Shutdown()
}
