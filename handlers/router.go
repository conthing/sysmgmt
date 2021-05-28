package handlers

import (
	"fmt"

	"github.com/conthing/sysmgmt/auth"

	"github.com/gin-gonic/gin"
)

// Response 通用HTTP回复body格式，除了ping/version/status等请求可以纯文本回复以外，都必须用此格式，且正常回复的Code必须200
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Run http service
func Run(port int) error {
	gin.SetMode(gin.ReleaseMode)
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
	authGroup.Use(auth.GINGuardExport())
	{
		authGroup.POST("/user/logout", Logout)
		authGroup.GET("/sn", GetMac)
		authGroup.GET("/version", GetVersion)
		authGroup.GET("/net", GetNetInfo)
		authGroup.POST("/net", SetNet)
		authGroup.GET("/time", GetTimeInfo)
		authGroup.POST("/time", SetTime)
		authGroup.GET("/envior", GetEnviorList)
		authGroup.GET("/envior/:name", GetEnvior)
		authGroup.POST("/envior", SetEnvior)

		authGroup.POST("/reboot", Reboot)
		authGroup.POST("/upgrade", Upgrade)
		authGroup.POST("/export", Export)
		authGroup.POST("/import", Import)

		// 别名设置
		//authGroup.GET("/alias", GetAlias)
		//authGroup.POST("/alias", SetAlias)
		// 地区设置
		//authGroup.GET("/location", GetLocation)
		//authGroup.POST("/location", SetLocation)
	}
	return r.Run(fmt.Sprintf(":%d", port))
}
