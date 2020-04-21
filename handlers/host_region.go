package handlers

import (
	"net/http"
	"sysmgmt-next/dto"
	"sysmgmt-next/redis"

	"github.com/gin-gonic/gin"
)

// GetRegion 获取地区
func GetRegion(c *gin.Context) {
	region, err := redis.GetRegion()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Resp{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, dto.Resp{
		Data: regionBody{Region: region},
	})
}

type regionBody struct {
	Region string `json:"region"`
}

// SetRegion 设置地区
func SetRegion(c *gin.Context) {
	var info regionBody
	err := c.ShouldBindJSON(&info)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Resp{
			Message: err.Error(),
		})
		return
	}
	err = redis.SaveRegion(info.Region)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Resp{
			Message: err.Error(),
		})
		return
	}
}
