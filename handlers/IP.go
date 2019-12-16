package handlers

import (
	"log"
	"net/http"
	"os/exec"
	"sysmgmt-next/config"
	"sysmgmt-next/dto"

	"github.com/conthing/utils/common"
	"github.com/gin-gonic/gin"
)

// PutIP 修改IP
func PutIP(c *gin.Context) {
	var netInfo dto.NetInfo
	if err := c.ShouldBindJSON(&netInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(netInfo)
	shellPath := config.Conf.ShellPath + "modifiedIP.sh"
	go runShell(shellPath, netInfo)
}

func runShell(shellPath string, s dto.NetInfo) {
	netname := common.GetMajorInterface()
	if netname == "" {
		log.Fatal("GetMajorInterface failed")
	} else {
		command := exec.Command(shellPath, common.GetMajorInterface(), s.Nettype, s.Address, s.Netmask, s.Gateway) //初始化Cmd
		out, _ := command.Output()
		log.Printf("modified ip output:%s", string(out))
	}
}
