package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"github.com/conthing/sysmgmt/dto"

	"github.com/conthing/utils/common"
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

const networkSettingFile = "/etc/network/interfaces"

// 从 /etc/network/interfaces 文件中读取网络设置，此文件必须
//    从"auto eth0开始"到"\n\n"结束是这样的设置：
//      DHCP时
//auto eth0
//iface eth0 inet dhcp
//
//      静态IP时
//auto eth0
//iface eth0 inet static
//	address 192.168.0.101
//	netmask	255.255.255.0
//	gateway 192.168.0.1
//
func getCurrentNetInfo(info *dto.NetInfo) error {
	data, err := ioutil.ReadFile(networkSettingFile)
	if err != nil {
		return fmt.Errorf("network interfaces file read failed: %v", err)
	}
	ethx := common.GetMajorInterface() // 确定接口名称

	contents := strings.ReplaceAll(string(data), "\r", "") // 去除\r
	contents = strings.ReplaceAll(contents, "\t", " ")     // 去除Tab
	begin := strings.Index(contents, "auto "+ethx)
	if begin == -1 {
		return fmt.Errorf("network interfaces file don't contain \"auto %s\"", ethx)
	}
	end := begin + strings.Index(contents[begin:], "\n\n")
	if end == -1 {
		return fmt.Errorf("network interfaces file don't contain \"\\n\\n\"")
	}
	//common.Log.Debugf("begin:%v,end:%v", begin, end)

	setting := contents[begin:end]
	//common.Log.Debugf("setting:%v", setting)
	setting = strings.ReplaceAll(setting, "auto "+ethx, "")           // 去除"auto eth0"
	setting = strings.ReplaceAll(setting, "iface "+ethx+" inet ", "") // 去除"iface eth0 inet"
	//common.Log.Debugf("setting:%v", setting)

	strs := strings.Fields(setting)
	//common.Log.Debugf("strs:%v", strs)
	length := len(strs)
	for i, s := range strs {
		if s == "address" && i < length-1 {
			info.Address = strs[i+1]
			info.DHCP = false
		}
		if s == "netmask" && i < length-1 {
			info.Netmask = strs[i+1]
			info.DHCP = false
		}
		if s == "gateway" && i < length-1 {
			info.Gateway = strs[i+1]
			info.DHCP = false
		}
		if s == "dhcp" {
			info.DHCP = true
		}
	}

	return nil
}

// PutNet 修改IP
func PutNet(c *gin.Context) {
	var info dto.NetInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := setNetInfo(&info); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

func setNetInfo(info *dto.NetInfo) error {
	data, err := ioutil.ReadFile(networkSettingFile)
	if err != nil {
		return fmt.Errorf("network interfaces file read failed: %v", err)
	}
	ethx := common.GetMajorInterface() // 确定接口名称

	contents := strings.ReplaceAll(string(data), "\r", "") // 去除\r
	begin := strings.Index(contents, "auto "+ethx)
	end := begin + strings.Index(contents[begin:], "\n\n")
	prefix := contents[:begin]
	postfix := contents[end:]

	setting := "auto " + ethx + "\niface " + ethx + " inet "
	if info.DHCP {
		setting = setting + "dhcp"
	} else {
		setting = setting + "static\n\taddress " + info.Address + "\n\tnetmask " + info.Netmask + "\n\tgateway " + info.Gateway
	}

	//common.Log.Debugf("setting:%v", setting)
	err = ioutil.WriteFile(networkSettingFile, []byte(prefix+setting+postfix), 0777)
	if err != nil {
		return fmt.Errorf("network interfaces file write failed: %v", err)
	}

	// 使网络设置立即生效
	go flushNetSetting(ethx)

	return nil
}

func flushNetSetting(ethx string) {
	command := exec.Command("ip", "addr", "flush", "dev", ethx) //初始化Cmd
	_, _ = command.Output()
	//common.Log.Debugf("ip addr flush output:%s", string(out))

	command = exec.Command("/etc/init.d/networking", "restart") //初始化Cmd
	_, _ = command.Output()
	//common.Log.Debugf("modified ip output:%s", string(out))
}
