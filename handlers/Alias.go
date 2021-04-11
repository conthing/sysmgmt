package handlers

import (
	"net/http"
	"strings"

	"github.com/conthing/sysmgmt/db"
	"github.com/conthing/sysmgmt/dto"
	"github.com/conthing/sysmgmt/models"

	"github.com/gin-gonic/gin"
)

type aliasBody struct {
	Alias string `json:"alias"`
}

// GetAlias 获取名字
func GetAlias(c *gin.Context) {
	var alias string
	envior := models.Envior{Name: "alias"}
	err := db.GetEnvior(&envior)
	if err != nil {
		alias = "Unknown"
	} else {
		alias = strings.TrimSpace(envior.Value)
		if alias == "" {
			alias = "Unknown"
		}
	}

	c.JSON(http.StatusOK, dto.Resp{
		Code: http.StatusOK,
		Data: aliasBody{Alias: alias},
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

	envior := models.Envior{Name: "alias", Value: strings.TrimSpace(info.Alias)}
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
		Data: aliasBody{Alias: info.Alias},
	})
}
