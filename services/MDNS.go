package services

import (
	"sysmgmt-next/config"

	"github.com/conthing/utils/common"
	"github.com/grandcat/zeroconf"
)

var server *zeroconf.Server

// MDNS 服务
func MDNS(cnf config.MDNS) {
	var err error
	server, err = zeroconf.Register(cnf.Name, "_workstation._tcp", "local.", cnf.Port, []string{"txtv=0", "lo=1", "la=2"}, nil)
	if err != nil {
		common.Log.Errorf("MDNS start failed: %v", err)
	}
	common.Log.Infof("MDNS start on %d", cnf.Port)
}

// StopMDNS 关闭MDNS
func StopMDNS() {
	server.Shutdown()
}
