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

// PutNet 修改IP
func PutNet(c *gin.Context) {
	var info dto.NetInfo
	if err := c.ShouldBindJSON(&info); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Println(info)
	shellPath := config.Conf.ShellPath + "modifiedIP.sh"
	go runShell(shellPath, info)
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
