package main

import (
	"backend/apis"
	"backend/middlewares"
	"backend/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	r := gin.Default()
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowAllOrigins = true
	corsCfg.AddAllowHeaders("Authorization")
	r.Use(cors.New(corsCfg))

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.User{})

	r.GET("/v1/user/nonce", apis.GetLoginNonceAPI)
	r.POST("/v1/user/auth", apis.MetamaskLoginAPI)
	r.GET("/v1/user/profile", middlewares.AuthMiddleware(), apis.GetProfileAPI)
	r.PATCH("/v1/user/modify-name", middlewares.AuthMiddleware(), apis.ModifyUsernameAPI)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(":3040")
}
