package main

import (
	"asset-management-service/internal/database"
	"asset-management-service/internal/handler"
	"asset-management-service/internal/repository"
	"asset-management-service/internal/service"
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewAssetRepository(db)
	svc := service.NewAssetService(repo)
	handler := handler.NewAssetHandler(svc)
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // domain của frontend React/Vue
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	api := r.Group("/api")
	{
		api.GET("/assets", handler.GetAssets)
		api.GET("/assets/:ip", handler.GetAssetByIP)
		api.POST("/assets", handler.CreateAsset)
		api.PUT("/assets/:ip", handler.UpdateAsset)
		api.DELETE("/assets/:ip", handler.DeleteAsset)
		api.PATCH("assets/:ip/manage", handler.ManageAsset)

		api.GET("assets/:ip/services", handler.GetServices)
		api.POST("assets/:ip/services", handler.AddService)
		api.PUT("assets/:ip/services/:service_id", handler.UpdateService)
		api.DELETE("assets/:ip/services/:service_id", handler.DeleteService)
		api.PATCH("assets/:ip/services/:service_id/manage", handler.ManageService)
	}
	go func() {
		for {
			err := service.SyncFromLogService(context.Background(), repo)
			if err != nil {
				log.Fatal(err)
			} else {
				log.Println("sync complete")
			}
			time.Sleep(30 * time.Second)
		}
	}()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	r.Run(":" + port)
}
