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
		// v1.GET("/ping", handlers.Ping)                      // ✅
		// v1.GET("/mac", handlers.GetMac)                     // ✅
		// v1.PUT("/net", handlers.PutNet)                     // ✅
		// v1.GET("/net/info", handlers.GetNetInfo)            // ✅
		// v1.GET("/version-list", handlers.GetVersionList)    // ✅
		// v1.PUT("/time", handlers.PutTime)                   // ✅
		// v1.GET("/time/info", handlers.GetTimeInfo)          // ✅
		// v1.PUT("/reboot", handlers.Reboot)                  // ❌
		v1.POST("/update/file-upload", handlers.FileUpload) // ❌

	}
	router.Run(fmt.Sprintf(":%d", cnf.Port))
}
