package handlers

import (
	"net/http"
	"strings"

	"github.com/conthing/sysmgmt/db"
	"github.com/conthing/sysmgmt/dto"

	"github.com/gin-gonic/gin"
)

// GetEnviorList 获取所有环境变量
func GetEnviorList(c *gin.Context) {
	list, err := db.GetEnviorList()
	if err != nil {
		c.JSON(http.StatusOK, dto.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Resp{
		Code: http.StatusOK,
		Data: list,
	})
}

// GetEnviorList 获取所有环境变量
func GetEnvior(c *gin.Context) {
	k := c.Param("name")
	k = strings.ToUpper(strings.TrimSpace(k))
	if k == "" {
		c.JSON(http.StatusOK, dto.Resp{
			Code:    http.StatusBadRequest,
			Message: "Envior name NULL",
		})
		return
	}

	v := db.GetEnv(k)

	c.JSON(http.StatusOK, dto.Resp{
		Code: http.StatusOK,
		Data: map[string]string{k: v},
	})
}

// SetEnvior 设置环境变量
func SetEnvior(c *gin.Context) {
	envs := make(map[string]string)
	done := make(map[string]string)
	err := c.ShouldBindJSON(&envs)
	if err != nil {
		c.JSON(http.StatusOK, dto.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	for k, v := range envs {
		k = strings.ToUpper(strings.TrimSpace(k))
		v = strings.TrimSpace(v)
		if k == "" {
			c.JSON(http.StatusOK, dto.Resp{
				Code:    http.StatusBadRequest,
				Message: "Envior name NULL",
				Data:    done,
			})
			return
		}

		err = db.SetEnv(k, v)
		if err != nil {
			c.JSON(http.StatusOK, dto.Resp{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
				Data:    done,
			})
			return
		}
		done[k] = v
	}

	c.JSON(http.StatusOK, dto.Resp{
		Code: http.StatusOK,
		Data: done,
	})
}

type aliasBody struct {
	Alias string `json:"alias"`
}

// GetAlias 获取名字
func GetAlias(c *gin.Context) {
	c.JSON(http.StatusOK, dto.Resp{
		Code: http.StatusOK,
		Data: aliasBody{Alias: strings.TrimSpace(db.GetEnv("ALIAS"))},
	})
}

// SetAlias 设置名字
func SetAlias(c *gin.Context) {
	var info aliasBody
	err := c.ShouldBindJSON(&info)
	if err != nil || info.Alias == "" {
		c.JSON(http.StatusOK, dto.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	err = db.SetEnv("ALIAS", strings.TrimSpace(info.Alias))
	if err != nil {
		c.JSON(http.StatusOK, dto.Resp{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Resp{
		Code: http.StatusOK,
		Data: aliasBody{Alias: info.Alias},
	})
}

type locationBody struct {
	Location string `json:"location"`
}

// GetLocation 获取位置
func GetLocation(c *gin.Context) {
	c.JSON(http.StatusOK, dto.Resp{
		Code: http.StatusOK,
		Data: locationBody{Location: strings.TrimSpace(db.GetEnv("LOCATION"))},
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

	err = db.SetEnv("LOCATION", strings.TrimSpace(info.Location))
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
