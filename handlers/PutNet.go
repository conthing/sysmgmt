package handlers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"sysmgmt-next/dto"

	"github.com/conthing/utils/common"
	"github.com/gin-gonic/gin"
)

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
