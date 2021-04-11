package auth

import (
	"net/http"

	"github.com/conthing/sysmgmt/db"
	"github.com/conthing/sysmgmt/dto"
	"github.com/conthing/sysmgmt/models"

	"github.com/conthing/utils/common"
	"github.com/gin-gonic/gin"
)

var tokenRandomMap = make(map[string]int)

func SetTokenRandom(userName string, tokenRandom int) {
	tokenRandomMap[userName] = tokenRandom
}

// GINGuard 权限验证守卫
func GINGuard() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求前
		tokenString := c.Request.Header.Get("Authorization")
		claims, err := ParseToken(tokenString)

		if err != nil {
			common.Log.Errorf("ParseToken failed: %v", err)
			c.JSON(http.StatusOK, dto.Resp{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		// 如果map里没有random，则从数据库中加载
		_, ok := tokenRandomMap[claims.Username]
		if !ok {
			user := models.User{Username: claims.Username}
			err = db.GetUser(&user)
			if err != nil {
				common.Log.Errorf("db.GetUser %s failed: %v", claims.Username, err)
				c.JSON(http.StatusOK, dto.Resp{
					Code:    http.StatusUnauthorized,
					Message: "Invalid token random",
				})
				c.Abort()
				return
			}
			tokenRandomMap[claims.Username] = user.TokenRandom
		}

		if claims.TokenRandom != tokenRandomMap[claims.Username] {
			common.Log.Errorf("TokenRandom %d should be %d", claims.TokenRandom, tokenRandomMap[claims.Username])
			c.JSON(http.StatusOK, dto.Resp{
				Code:    http.StatusUnauthorized,
				Message: "Invalid token random",
			})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}

// GINGuard 权限验证守卫，供其他进程使用 todo 不能让旧的token失效
func GINGuardExport() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求前
		tokenString := c.Request.Header.Get("Authorization")
		claims, err := ParseToken(tokenString)

		if err != nil {
			common.Log.Errorf("ParseToken failed: %v", err)
			c.JSON(http.StatusOK, dto.Resp{
				Code:    http.StatusUnauthorized,
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}
