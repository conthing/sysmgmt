package router

import (
	"fmt"
	"sysmgmt-next/config"
	"sysmgmt-next/handlers"

	"github.com/gin-gonic/gin"
)

// Service 启动路由
func Service(cnf config.Config) {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/ping", handlers.Ping)
	}

	router.Run(fmt.Sprintf(":%d", cnf.Port))
}
