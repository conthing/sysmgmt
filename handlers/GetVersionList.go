package handlers

import (
	"io/ioutil"
	"net/http"
	"strings"
	"sysmgmt-next/config"
	"sysmgmt-next/dto"

	"github.com/gin-gonic/gin"
)

// GetVersionList 获取版本号
func GetVersionList(c *gin.Context) {
	servicenamelist := config.Conf.ServiceNamelist
	serviceportlist := config.Conf.ServicePortlist
	var version dto.VersionInfo
	versionList := make([]dto.VersionInfo, 0)
	for i, servicename := range servicenamelist {
		resp, err := http.Get("http://localhost:" + serviceportlist[i] + "/api/v1/version")
		if err != nil || resp.StatusCode != 200 {
			continue
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		version.Name = servicename
		str := string(body)
		strArray := strings.Split(str, " ")
		version.Version = strArray[0]
		version.CreatedTime = strArray[1] + " " + strArray[2]
		versionList = append(versionList, version)
	}
	c.JSON(200, versionList)
}
