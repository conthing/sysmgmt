package handlers

import (
	"net/http"
	"sysmgmt-next/dto"
	"sysmgmt-next/redis"

	"github.com/gin-gonic/gin"
)

// GetAlias 获取名字
func GetAlias(c *gin.Context) {
	alias, err := redis.GetAlias()
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Resp{
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, dto.Resp{
		Data: aliasBody{Alias: alias},
	})
}

type aliasBody struct {
	Alias string `json:"alias"`
}

// SetAlias 设置名字
func SetAlias(c *gin.Context) {
	var info aliasBody
	err := c.ShouldBindJSON(&info)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Resp{
			Message: err.Error(),
		})
		return
	}
	err = redis.SaveAlias(info.Alias)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Resp{
			Message: err.Error(),
		})
		return
	}
}
