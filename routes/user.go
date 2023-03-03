package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ichthoth/jwt-auth/controllers"
	"github.com/ichthoth/jwt-auth/middleware"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users", controllers.GetUser())
	incomingRoutes.GET("/users/:user_id", controllers.GetUserID())
}
