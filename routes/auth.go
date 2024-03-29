package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ichthoth/jwt-auth/controllers"
)

//Authentication routes
func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signup", controllers.SignUp())
	incomingRoutes.POST("users/login", controllers.Login())

}
