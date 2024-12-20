package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Doreen-Onyango/zingiratech/backend/internal/handlers"
	"github.com/Doreen-Onyango/zingiratech/backend/internal/middlewares"
	"github.com/Doreen-Onyango/zingiratech/backend/internal/utils"
)

// InitRoutes initializes all application routes and serves static files.
func InitRoutes(mux *http.ServeMux) error {
	// Resolve the static files directory
	dir, err := utils.GetProjectRootPath("frontend", "static")
	if err != nil {
		return fmt.Errorf("failed to resolve static directory: %w", err)
	}

	// Log the static directory for debugging
	log.Printf("Static files will be served from: %s\n", dir)

	// Serve static files under /static/
	fs := http.FileServer(http.Dir(dir))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Register other application routes
	registerRoutes(mux)

	log.Println("Routes initialized successfully")
	return nil
}

// registerRoutes sets up route handlers for the application.
func registerRoutes(mux *http.ServeMux) {
	// Public routes
	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/about", handlers.AboutHandler)
	mux.HandleFunc("/login", handlers.LoginHandler)
	mux.HandleFunc("/signup", handlers.SignupHandler)

	// Protected routes with middleware
	protectedRoutes := http.NewServeMux()
	protectedRoutes.HandleFunc("/dashboard", handlers.DashboardHandler)
	protectedRoutes.HandleFunc("/schedule-pickup", handlers.SchedulePickupHandler)

	// Apply auth middleware to protected routes
	mux.Handle("/dashboard/", middlewares.AuthMiddleware(protectedRoutes))
	mux.Handle("/api/", middlewares.AuthMiddleware(protectedRoutes))

	log.Println("All application routes registered successfully")
}
