package handlers

import (
	"net/http"

	"github.com/conthing/sysmgmt/services"
	"github.com/gin-gonic/gin"
)

// Reboot 重启
func Reboot(c *gin.Context) {

	services.RebootLater()

	c.JSON(http.StatusOK, Response{
		Code: http.StatusOK,
		Data: "reboot in 3 seconds.",
	})
}
