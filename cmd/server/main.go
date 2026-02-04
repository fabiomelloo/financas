package main

import (
	"financas/database"
	"financas/internal/controllers"
	"financas/internal/repositories"
	"financas/internal/routes"
	"financas/internal/services"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Initialize database connection
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize layers (Dependency Injection)
	repository := repositories.NewExpenseRepository(db)
	service := services.NewExpenseService(repository)
	controller := controllers.NewExpenseController(service)

	// Register routes
	routes.RegisterRoutes(controller)

	// Serve static files (CSS, JS, images)
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start HTTP server
	fmt.Println("Server running at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
