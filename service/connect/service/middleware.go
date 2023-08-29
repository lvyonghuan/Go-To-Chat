package service

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func checkToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取token
		token := c.GetHeader("Authorization")
		if token == "" {
			respError(c, 401, errors.New("token为空"))
			return
		}

		//解析token
		claims, err := parseToken(token)
		if err != nil {
			respError(c, 401, err)
			return
		}

		//验证token
		if !claims.Valid {
			respError(c, 401, errors.New("token无效"))
			return
		}

		c.Set("user_id", claims.Claims.(jwt.MapClaims)["userID"])

		//验证通过，继续处理请求
		c.Next()
	}
}

// 解析token
func parseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
}
