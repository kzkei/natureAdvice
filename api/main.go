package main

import (
	"log"

	"github.com/kzkei/natureAdvice/api/config"
	"github.com/kzkei/natureAdvice/api/server"

	"github.com/kzkei/natureAdvice/api/database"
)

func main() {
	// load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// connect to database
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// test ping
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping db: %v", err)
	}
	log.Println("db connected")

	// set up router
	router := server.SetupRouter(db) // returns *gin.Engine

	// Gin has own Run method
	router.Run(":" + cfg.APIPort)
}
