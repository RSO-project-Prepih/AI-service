package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/RSO-project-Prepih/AI-service/docs"
	"github.com/RSO-project-Prepih/AI-service/handlers"
	"github.com/RSO-project-Prepih/AI-service/health"
	"github.com/RSO-project-Prepih/AI-service/prometheus"
	"github.com/gin-gonic/gin"
	"github.com/rs/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title AI Service API
// @version 1.0
// @description This is a sample server for AI Service.
// @BasePath /v1

func main() {
	log.Println("Starting the AI service...")
	r := gin.Default()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // your frontend's origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "DELETE"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           int(12 * time.Hour / time.Second),
	})

	r.Use(func(ctx *gin.Context) {
		c.HandlerFunc(ctx.Writer, ctx.Request)
		ctx.Next()
	})

	// Define the routes for the famous places
	r.GET("/famous-places", func(c *gin.Context) {
		log.Println("Fetching famous places...")
		famousPlaces, err := handlers.FetchFamousPlaces()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Println("Error fetching famous places:", err)
			return
		}
		log.Println("Famous places fetched successfully")
		c.JSON(http.StatusOK, famousPlaces)
	})

	// Define the routes for the color enhancement
	r.POST("/enhance-color", func(c *gin.Context) {
		log.Println("Starting color enhancement...")
		// Extract user ID and image ID from request
		userID := c.Query("userID")
		imageID := c.Query("imageID")

		if userID == "" || imageID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "userID and imageID are required"})
			log.Println("Error starting color enhancement:", userID, imageID)
			return
		}

		handlers.PostColorEnhancementPhoto(userID, imageID)
		c.JSON(http.StatusOK, gin.H{"message": "Completed color enhancement"})
	})

	// Define the fetching of images that are enhancement
	r.GET("/image-processing", func(c *gin.Context) {
		log.Println("Fetching image processing data...")
		imageProcessingData, err := handlers.GetImageProcessingPhotos()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Println("Error fetching image processing data:", err)
			return
		}
		log.Println("Image processing data fetched successfully")
		c.JSON(http.StatusOK, imageProcessingData)
	})

	// Define the routes for the health check
	liveHandler, readyHandler := health.HealthCheckHandler()
	r.GET("/live", gin.WrapH(liveHandler))
	r.GET("/ready", gin.WrapH(readyHandler))

	// Define the routes for the metrics
	r.GET("/metrics", gin.WrapH(prometheus.GetMetrics()))

	// Define the routes for the swagger
	r.GET("/openapi/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start the server
	srver := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srver.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	qoit := make(chan os.Signal, 1)
	signal.Notify(qoit, syscall.SIGINT, syscall.SIGTERM)
	<-qoit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srver.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
