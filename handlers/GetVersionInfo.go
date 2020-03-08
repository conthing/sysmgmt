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

	var version dto.SubVersionInfo
	microservicelist := config.Conf.MicroServiceList
	for _, microservice := range microservicelist {
		resp, err := http.Get("http://localhost:" + string(microservice.Port) + "/api/v1/version")
		if err != nil || resp.StatusCode != 200 {
			continue
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		str := string(body)
		//解释:这里不是删除空格，而是以空格做为分隔符，将str(string)拆分成strArry([]string),  str的打印的信息是("1.1.4 2020-03-06 12:20:30"),strArry的打印信息是["1.1.4"， "2020-03-06"， "12:20:30"]
		strArry := strings.Split(str, "") //todo review 这里为什么把空格删了！！！！修改这些东西不可不仔细
		version.Name = microservice.Name
		version.Version = strArry[0]
		version.BuildTime = strArry[1] + "" + strArry[2]
		globalVersion.SubVersion = append(globalVersion.SubVersion, version)
	}
	c.JSON(http.StatusOK, globalVersion)
}
