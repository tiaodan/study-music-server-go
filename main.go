package main

import (
	"fmt"
	"log"
	"study-music-server-go/config"
	"study-music-server-go/mapper"
	"study-music-server-go/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	if err := mapper.InitDB(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	fmt.Println("Database connected successfully")

	// Create Gin router
	r := gin.Default()

	// Apply middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS middleware
	r.Use(func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
		}
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
		c.Header("Access-Control-Max-Age", "3600")
		c.Header("Access-Control-Allow-Headers", "x_requested_with, Authorization, Content-Type, token")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Setup routes first
	routes.SetupRoutes(r)

	// Static file serving - use NoRoute to avoid conflict with API routes
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		// Handle /song/*filepath
		if len(path) > 6 && path[:6] == "/song/" {
			c.File("./resource/song" + path[5:])
			return
		}
		// Handle /img/*filepath
		if len(path) > 5 && path[:5] == "/img/" {
			c.File("./resource/img" + path[4:])
			return
		}
	})

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	fmt.Printf("Server starting on %s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
