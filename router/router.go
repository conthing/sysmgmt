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
		v1.GET("/ping", handlers.Ping)      // ✅
		v1.GET("/sn", handlers.GetMac)      // ✅
		v1.PUT("/net", handlers.PutNet)     // ✅
		v1.GET("/net", handlers.GetNetInfo) // ✅
		v1.GET("/version", handlers.GetVersion)
		v1.PUT("/time", handlers.PutTime)                   // ✅
		v1.GET("/time", handlers.GetTimeInfo)               // ✅
		v1.PUT("/reboot", handlers.Reboot)                  // ✅
		v1.POST("/update/file-upload", handlers.FileUpload) // ✅

		// 别名设置
		v1.GET("/alias", handlers.GetAlias)
		v1.POST("/alias", handlers.SetAlias)
		// 地区设置
		v1.GET("/region", handlers.GetRegion)
		v1.POST("/region", handlers.SetRegion)
	}
	router.Run(fmt.Sprintf(":%d", cnf.Port))
}
