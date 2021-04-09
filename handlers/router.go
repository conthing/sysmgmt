package handlers

import (
	"fmt"

	"github.com/conthing/sysmgmt/auth"
	"github.com/conthing/sysmgmt/config"

	"github.com/gin-gonic/gin"
)

// Run http service
func Run(cnf *config.HTTPConfig) {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.GET("/ping", Ping)
		v1.GET("/user/ready", Ready)
		v1.POST("/user/signup", SignUp)
		v1.POST("/user/login", Login)
		v1.POST("/user/passwd", Passwd)
	}

	// jwt凭证验证接口
	authGroup := r.Group("/api/v1")
	authGroup.Use(auth.GINGuard())
	{
		authGroup.POST("/user/logout", Logout)
		authGroup.GET("/sn", GetMac)      // ✅
		authGroup.PUT("/net", PutNet)     // ✅
		authGroup.GET("/net", GetNetInfo) // ✅
		authGroup.GET("/version", GetVersion)
		authGroup.PUT("/time", PutTime)                   // ✅
		authGroup.GET("/time", GetTimeInfo)               // ✅
		authGroup.PUT("/reboot", Reboot)                  // ✅
		authGroup.POST("/update/file-upload", FileUpload) // ✅

		// 别名设置
		authGroup.GET("/alias", GetAlias)
		authGroup.POST("/alias", SetAlias)
		// 地区设置
		authGroup.GET("/location", GetLocation)
		authGroup.POST("/location", SetLocation)
		authGroup.GET("/region", GetRegion)
		authGroup.POST("/region", SetRegion)
	}
	r.Run(fmt.Sprintf(":%d", cnf.Port))
}
