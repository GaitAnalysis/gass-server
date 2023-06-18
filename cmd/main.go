package main

import (
	"fmt"
	"gass/api"
	"gass/middleware"
	"gass/repository"
	v1 "gass/service"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

func main() {
	db, err := repository.NewGassRepository()
	if err != nil {
		log.Fatal("db failed to initialize")
	}

	svc := v1.NewService(db, &http.Client{})
	server := api.NewServer(svc)

	router := gin.Default()
	router.Use(cors.AllowAll())
	router.MaxMultipartMemory = 2 << 30 // 1073 MiB Upload Limit

	publicRoutes := router.Group("/auth")
	publicRoutes.POST("/register", server.Register)
	publicRoutes.POST("/login", server.Login)
	publicRoutes.POST("/analysis", server.AddAnalysis)

	protectedRoutes := router.Group("/api")
	protectedRoutes.Use(middleware.JWTAuthMiddleware())

	protectedRoutes.POST("/upload", server.AddUpload)
	protectedRoutes.GET("/upload", server.GetUploads)
	protectedRoutes.GET("/analysis", server.GetAnalyses)

	router.Run(":" + os.Getenv("PORT"))
	fmt.Println("Server started on port 8080")
}
