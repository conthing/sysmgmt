package handlers

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/conthing/sysmgmt/dto"

	"github.com/gin-gonic/gin"
)

// GetAlias 获取名字
func GetAlias(c *gin.Context) {
	var alias string

	out, err := ioutil.ReadFile("../data/.alias")
	if err != nil {
		alias = "Unknown"
	} else {
		alias = strings.TrimSpace(string(out))
		if alias == "" {
			alias = "Unknown"
		}
	}

	c.JSON(http.StatusOK, dto.Resp{
		Code: http.StatusOK,
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
	if err != nil || info.Alias == "" {
		c.JSON(http.StatusOK, dto.Resp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	err = ioutil.WriteFile("../data/.alias", []byte(info.Alias), 0666)
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
