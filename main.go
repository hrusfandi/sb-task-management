package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hrusfandi/sb-task-management/config"
	"github.com/hrusfandi/sb-task-management/database"
	"github.com/hrusfandi/sb-task-management/routes"
)

func main() {
	config.LoadConfig()

	database.InitDB()

	r := routes.SetupRoutes(database.GetDB())

	log.Println("Task Management API is starting...")
	log.Printf("Server running on http://localhost:%s", config.AppConfig.Port)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.AppConfig.Port), r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}