package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// todo 和 handlers中的定义重复
// Response 通用HTTP回复body格式，除了ping/version/status等请求可以纯文本回复以外，都必须用此格式，且正常回复的Code必须200
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// GINGuardExport 权限验证守卫，供其他进程使用 todo 不能让旧的token失效
func GINGuardExport() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求前
		tokenString := c.Request.Header.Get("Authorization")
		claims, err := ParseToken(tokenString)

		if err != nil {
			c.JSON(http.StatusOK, Response{
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
