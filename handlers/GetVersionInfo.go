package handlers

import (
	"io/ioutil"
	"net/http"
	"os/exec"
	"strconv"
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
		return
	}
	text := strings.SplitAfterN(string(out), "\n", 2) //用第一个 \n 分割字符串
	globalVersion.Version = strings.TrimSpace(text[0])
	if len(text) > 1 {
		globalVersion.Description = strings.TrimSpace(text[1])
	}
	globalVersion.SubVersion = append(globalVersion.SubVersion, dto.SubVersionInfo{Name: "sysmgmt", Version: common.Version, BuildTime: common.BuildTime})

	var version dto.SubVersionInfo
	microservicelist := config.Conf.MicroServiceList
	for _, microservice := range microservicelist {
		url := "http://localhost:" + strconv.FormatInt(int64(microservice.Port), 10) + "/api/v1/version"
		resp, err := http.Get(url)
		if err == nil {
			defer resp.Body.Close()
		} else {
			common.Log.Errorf("%s Get failed: %v", url, err)
			continue
		}
		if resp.StatusCode != 200 {
			common.Log.Errorf("%s Get failed: code:%d", url, resp.StatusCode)
			continue
		}
		body, _ := ioutil.ReadAll(resp.Body)
		str := string(body)
		common.Log.Debugf("%s Get: %s", url, str)
		strArry := strings.Split(str, " ")
		version.Name = microservice.Name
		version.Version = strArry[0]
		version.BuildTime = strArry[1] + " " + strArry[2]
		globalVersion.SubVersion = append(globalVersion.SubVersion, version)
	}
	c.JSON(http.StatusOK, globalVersion)
}
