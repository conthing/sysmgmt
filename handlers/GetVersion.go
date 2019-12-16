package handlers

import (
	"io/ioutil"
	"net/http"
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
		version.Version = string(body)
		versionList = append(versionList, version)
	}
	c.JSON(200, versionList)
}
