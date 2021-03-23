package main

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"ocr.service.authorization/app/middleware"
	UserDelivery "ocr.service.authorization/app/user/delivery"
	"ocr.service.authorization/config"
	"os"
)

// @title OCR AUTHORIZATION API
// @version 1.0
func main() {
	CONFIG, _ := config.NewConfig(nil)
	router := gin.New()
	router.Use(cors.Default())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	authMiddleware := middleware.NewAuth()
	userDelivery, err := UserDelivery.NewUserDelivery()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	routerApi := router.Group("/api/v1")
	routerApi.POST("/register", userDelivery.Create)
	routerApi.POST("/login", authMiddleware.LoginHandler)
	routerApi.POST("/reset_password", userDelivery.UpdatePassword)

	routerApiAuth := routerApi.Group("/auth")
	// Refresh time can be longer than token timeout
	routerApiAuth.GET("/refresh_token", authMiddleware.RefreshHandler)
	routerApiAuth.Use(authMiddleware.MiddlewareFunc())
	{
		routerApiAuth.POST("/user/:user_id", userDelivery.UpdateByID)
		routerApiAuth.GET("/user/:user_id", userDelivery.GetByID)
		//routerApiAuth.GET("/users", userDelivery.Gets)
	}
	// swagger
	router.GET("/swagger/swagger.json", func(c *gin.Context) {
		c.File("./docs/swagger.json")
	})
	url := ginSwagger.URL("/swagger/swagger.json") // The url pointing to API definition
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	fmt.Println("listening", CONFIG.GetString("NO_SSL_PORT"))
	if err := http.ListenAndServe(":"+CONFIG.GetString("NO_SSL_PORT"), router); err != nil {
		log.Fatal(err)
	}
}
