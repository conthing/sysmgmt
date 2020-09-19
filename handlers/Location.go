package handlers

import (
	"io/ioutil"
	"net/http"
	"strings"
	"github.com/conthing/sysmgmt/dto"

	"github.com/gin-gonic/gin"
)

// GetLocation 获取位置
func GetLocation(c *gin.Context) {
	var location string

	out, err := ioutil.ReadFile("../data/.location")
	if err != nil {
		location = "Unknown"
	} else {
		location = strings.TrimSpace(string(out))
		if location == "" {
			location = "Unknown"
		}
	}

	c.JSON(http.StatusOK, dto.Resp{
		Data: locationBody{Location: location},
	})
}

type locationBody struct {
	Location string `json:"location"`
}

// SetLocation 设置位置
func SetLocation(c *gin.Context) {
	var info locationBody
	err := c.ShouldBindJSON(&info)
	if err != nil || info.Location == "" {
		c.JSON(http.StatusBadRequest, dto.Resp{
			Message: err.Error(),
		})
		return
	}

	err = ioutil.WriteFile("../data/.location", []byte(info.Location), 0666)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Resp{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Resp{
		Data: locationBody{Location: info.Location},
	})
}
