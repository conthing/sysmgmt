package handlers

import (
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"sysmgmt-next/config"
	"sysmgmt-next/dto"

	"github.com/conthing/utils/common"
	"github.com/gin-gonic/gin"
)

// GetVersion 获取版本信息
func GetVersion(c *gin.Context) {
	var globalVersion dto.VersionInfo

	command := exec.Command("cat", "../VERSION") //初始化Cmd
	out, err := command.Output()
	if err != nil {
		common.Log.Errorf("open ../VERSION failed %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	text := strings.SplitAfterN(string(out), "\n", 2) //用第一个 \n 分割字符串
	globalVersion.Version = strings.TrimSpace(text[0])
	if len(text) > 1 {
		globalVersion.Description = strings.TrimSpace(text[1])
	}

	service := config.MicroService{}
	serviceportlist := service.ServicePortlist
	var version dto.SubVersionInfo
	for i, servicename := range serviceportlist {
		resp, err := http.Get("http://localhost:" + serviceportlist[i] + "/api/v1/version")
		if err != nil || resp.StatusCode != 200 {
			continue
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		str := string(body)
		strArray := strings.Split(str, " ")
		version.Name = servicename
		version.Version = strArray[0]
		version.BuildTime = strArray[1] + " " + strArray[2] //todo 没有检查
		globalVersion.SubVersion = append(globalVersion.SubVersion, version)
	}

	c.JSON(http.StatusOK, globalVersion)
}
