package services

import (
	"log"
	"sysmgmt-next/config"

	"github.com/grandcat/zeroconf"
)

// MDNS 服务
func MDNS(cnf config.MDNS) {
	server, err := zeroconf.Register(cnf.Name, "_workstation._tcp", "local.", cnf.Port, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("MDNS 启动成功")
	defer server.Shutdown()

}
