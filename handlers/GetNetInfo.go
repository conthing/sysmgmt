package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sysmgmt-next/dto"

	"github.com/gin-gonic/gin"
)

// GetNetInfo 获取系统网络
func GetNetInfo(c *gin.Context) {
	var info dto.NetInfo
	err := getCurrentNetInfo(&info)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, info)
}

// DHCPFile DHCP状态读取文件
var DHCPFile = "/etc/init.d/ipaddr.sh"

func getCurrentNetInfo(info *dto.NetInfo) error {
	data, err := ioutil.ReadFile(DHCPFile)
	if err != nil {
		return fmt.Errorf("IP 配置文件读取失败 %v", err)
	}
	strs := strings.Fields(string(data))
	//common.Log.Debugf("strs:%v", strs)
	length := len(strs)
	for i, s := range strs {
		if s == "netmask" && i >= 1 && i < length-1 {
			info.Address = strs[i-1]
			info.Netmask = strs[i+1]
			info.DHCP = false
		}
		if s == "gw" && i < length-1 {
			info.Gateway = strs[i+1]
			info.DHCP = false
		}
		if s == "udhcpc" {
			info.DHCP = true
		}
	}

	return nil
}
