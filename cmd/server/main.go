package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ichthoth/jwt-auth/routes"
	"github.com/joho/godotenv"
)

func Exec() error {
	err := godotenv.Load("env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.GET("api/v1", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"success": "Access granted for api v1 "})
	})
	router.GET("api/v2", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"success": "Access granted for api v2 "})
	})

	appMount := router.Run(":" + port)

	return appMount
}

func main() {

}
