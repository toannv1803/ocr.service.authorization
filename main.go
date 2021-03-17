package main

import (
	"fmt"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"ocr.service.authorization/app/middleware"
	UserDelivery "ocr.service.authorization/app/user/delivery"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	if port == "" {
		port = "8000"
	}

	authMiddleware := middleware.NewAuth()
	userDelivery, err := UserDelivery.NewUserDelivery()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	routerApi := router.Group("/api/v1")
	routerApi.POST("/register", userDelivery.Create)
	routerApi.POST("/login", authMiddleware.LoginHandler)

	routerApiAuth := routerApi.Group("/auth")
	// Refresh time can be longer than token timeout
	routerApiAuth.GET("/refresh_token", authMiddleware.RefreshHandler)
	routerApiAuth.Use(authMiddleware.MiddlewareFunc())
	{
		routerApiAuth.POST("/user/:user_id", userDelivery.UpdateByID)
		routerApiAuth.GET("/user/:user_id", userDelivery.GetByID)
		//routerApiAuth.GET("/users", userDelivery.Gets)
	}
	router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	fmt.Println("listening", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
