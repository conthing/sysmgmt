package handlers

import (
	"net/http"
	"strings"

	"github.com/conthing/sysmgmt/db"
	"github.com/conthing/sysmgmt/dto"
	"github.com/conthing/sysmgmt/models"

	"github.com/gin-gonic/gin"
)

type locationBody struct {
	Location string `json:"location"`
}

// GetLocation 获取位置
func GetLocation(c *gin.Context) {
	var location string
	envior := models.Envior{Name: "location"}
	err := db.GetEnvior(&envior)
	if err != nil {
		location = "Unknown"
	} else {
		location = strings.TrimSpace(envior.Value)
		if location == "" {
			location = "Unknown"
		}
	}

	c.JSON(http.StatusOK, dto.Resp{
		Code: http.StatusOK,
		Data: locationBody{Location: location},
	})
}

// SetLocation 设置位置
func SetLocation(c *gin.Context) {
	var info locationBody
	err := c.ShouldBindJSON(&info)
	if err != nil || info.Location == "" {
		c.JSON(http.StatusOK, dto.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	envior := models.Envior{Name: "location", Value: strings.TrimSpace(info.Location)}
	err = db.SetEnvior(&envior)
	if err != nil {
		c.JSON(http.StatusOK, dto.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Resp{
		Code: http.StatusOK,
		Data: locationBody{Location: info.Location},
	})
}
