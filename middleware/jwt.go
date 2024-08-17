package middleware

import (
	util "ginchat/utils"
	"ginchat/utils/e"
	"github.com/gin-gonic/gin"
	"time"
)

// JWT token验证中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		code = 200
		token := c.GetHeader("Authorization")
		if token == "" {
			code = 404
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ErrorAuthToken
			} else if time.Now().Unix() > claims.ExpiresAt { //验证token是否过期
				code = e.ErrorAuthCheckTokenTimeOut
			}
		}
		if code != e.Success {
			c.JSON(200, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"token":  token,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
