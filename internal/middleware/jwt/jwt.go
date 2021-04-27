package jwt

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"github.com/siangyeh8818/gin.project.template/pkg/env"
	"github.com/siangyeh8818/gin.project.template/pkg/util"
)

// JWT is jwt middleware
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = env.SUCCESS
		token := c.Query("token")
		if token == "" {
			code = env.INVALID_PARAMS
		} else {
			_, err := util.ParseToken(token)
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					code = env.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
				default:
					code = env.ERROR_AUTH_CHECK_TOKEN_FAIL
				}
			}
		}

		if code != env.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  env.GetMsg(code),
				"data": data,
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
