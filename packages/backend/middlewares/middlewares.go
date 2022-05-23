package middlewares

import (
	"backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func parseAuthorizationHeader(a string) (string, error) {
	if a == "" {
		return "", fmt.Errorf("authorization 为空字符串")
	}

	t := strings.Split(a, " ")

	if len(t) < 2 {
		return "", fmt.Errorf("authorization 格式不正确")
	}

	return t[1], nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Authorization 提取 JWT
		token, err := parseAuthorizationHeader(c.Request.Header.Get("Authorization"))

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusOK, gin.H{
				"msg":  "authentication failed",
				"code": 4001,
				"data": nil,
			})
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(token)

		if err != nil {
			log.Println(err)
			c.JSON(http.StatusOK, gin.H{
				"msg":  "authentication failed",
				"code": 4001,
				"data": nil,
			})
			c.Abort()
			return
		}

		c.Set("uid", claims.Uid)
		c.Set("publicAddress", claims.PublicAddress)
	}
}
