package services

import (
	"fmt"

	"github.com/conthing/utils/common"
	"github.com/grandcat/zeroconf"
)

// 发现服务配置
type MDNSConfig struct {
	Enable bool
	Name   string
	Port   int
}

var server *zeroconf.Server

// StartMDNS 开启MDNS服务
func StartMDNS(cnf *MDNSConfig) error {
	if cnf.Enable {
		var err error
		server, err = zeroconf.Register(cnf.Name, "_workstation._tcp", "local.", cnf.Port, []string{"txtv=0", "lo=1", "la=2"}, nil)
		if err != nil {
			common.Log.Errorf("MDNS start failed: %v", err)
			return fmt.Errorf("MDNS start failed: %w", err)
		}
		common.Log.Infof("MDNS start on %d", cnf.Port)
	}
	return nil
}

// StopMDNS 关闭MDNS服务
func StopMDNS() {
	if server != nil {
		common.Log.Infof("MDNS shuting down")
		server.Shutdown()
	}
}
