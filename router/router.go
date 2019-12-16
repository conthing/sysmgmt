package router

import (
	"fmt"
	"sysmgmt-next/config"
	"sysmgmt-next/handlers"

	"github.com/gin-gonic/gin"
)

// Service 启动路由
func Service(cnf config.Config) {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())
	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", handlers.Ping)
		// todo:
		v1.GET("/getMacAddr", handlers.GetMacAddr)
		v1.GET("/system/net", handlers.GetSystemNet)

		v1.GET("/version", handlers.GetVersionList)
		v1.GET("/get/time/info", handlers.GetTimeInfo)

		v1.PUT("/update/IP", handlers.PutIP)
		v1.PUT("/system/reboot", handlers.Reboot)
		v1.PUT("/modified/time", handlers.PutTime)

		v1.POST("/uploadFile", handlers.FileUpload)

	}
	router.Run(fmt.Sprintf(":%d", cnf.Port))
}
