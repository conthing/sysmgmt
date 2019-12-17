package handlers

import (
	"io/ioutil"
	"log"
	"strings"
	"sysmgmt-next/dto"

	"github.com/gin-gonic/gin"
)

// GetNetInfo 获取系统网络
func GetNetInfo(c *gin.Context) {
	var info dto.SystemNetInfo
	info.DHCPFlag = getDHCPFlag()
	c.JSON(200, info)
}

// DHCPFile DHCP状态读取文件
var DHCPFile = "/etc/init.d/ipaddr.sh"

func getDHCPFlag() bool {
	data, err := ioutil.ReadFile(DHCPFile)
	if err != nil {
		log.Fatal("IP 配置文件读取失败", err)
		return false
	}
	str := string(data)
	if strings.Contains(str, "udhcpc") {
		return true
	}
	return false
}
