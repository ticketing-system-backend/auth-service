package main

import (
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/ticketing-system-backend/auth-service/config"
	"github.com/ticketing-system-backend/auth-service/controller"
	docs "github.com/ticketing-system-backend/auth-service/docs"
	"github.com/ticketing-system-backend/auth-service/middleware"
	"github.com/ticketing-system-backend/auth-service/model"
	seeder "github.com/ticketing-system-backend/auth-service/seed"
)

// @title Go Gin Rest API
// @version 1.0
// @description A rest API in Go using Gin framework
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format **Bearer &lt;token&gt;**

func main() {
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})
	
	config.ConnectionDatabase()

	config.DB.AutoMigrate(&model.User{}, &model.Role{})

	seeder.SeedRoles()
	seeder.SeedSuperAdmin()

	docs.SwaggerInfo.BasePath = "/"

	auth := r.Group("/login")
	{
		auth.POST("/dashboard", controller.DashboardLogin)
		auth.POST("/mobile", controller.MobileLogin)
	}

	users := r.Group("/users", middleware.JWTAuth(), middleware.RequireRole("admin", "superadmin"))
	{
		users.GET("/", controller.GetAllUsers)
		users.GET("/:id", controller.GetUserById)
		users.POST("/", controller.CreateUser)
		users.PUT("/:id", controller.UpdateUser)
	}

	roles := r.Group("/roles", middleware.JWTAuth(), middleware.RequireRole("admin", "superadmin"))
	{
		roles.GET("/", controller.GetAllRoles)
		roles.GET("/:id", controller.GetRoleById)
		roles.POST("/", controller.CreateRole)
		roles.PUT("/:id", controller.UpdateRole)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":" + os.Getenv("PORT"))
}
