package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	ConnectDB()
	InitDB()

	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)

	// Auth routes (no middleware)
	r.Post("/login", handleLogin)
	r.Post("/register", handleRegister)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/profile", handleProfile)
	})

	// Admin routes
	r.Group(func(r chi.Router) {
		r.Use(adminMiddleware)
		r.Get("/logs", handleLogs)
	})

	// Internal router
	internal := chi.NewRouter()
	internal.Get("/health", healthCheck)

	go func() {
		log.Fatal(http.ListenAndServe("127.0.0.1:9091", internal))
	}()

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
