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
		v1.GET("/version", handlers.GetVersionList)
		// todo: 测试修改 IP
		v1.PUT("/update/IP", handlers.PutIP)
		// todo: 测试重启
		v1.PUT("/system/reboot", handlers.Reboot)
	}
	router.Run(fmt.Sprintf(":%d", cnf.Port))
}
