package main

import (
	"log"
	"net/http"

	"ecommerce-api/handler"
	"ecommerce-api/repository"
	"ecommerce-api/service"
)

func main() {
	// Load configuration
	cfg := LoadConfig()

	// Initialize database repository
	repoCfg := repository.Config{
		DBHost:     cfg.DBHost,
		DBUser:     cfg.DBUser,
		DBPassword: cfg.DBPassword,
		DBName:     cfg.DBName,
		DBPort:     cfg.DBPort,
		AdminUser:  cfg.AdminUser,
		AdminPass:  cfg.AdminPass,
	}

	postgresRepo, err := repository.NewPostgresRepository(repoCfg)
	if err != nil {
		log.Fatalf("Failed to initialize repository: %v", err)
	}

	// Create separate repository instances
	userRepo := &repository.UserRepo{PostgresRepository: postgresRepo}
	productRepo := &repository.ProductRepo{PostgresRepository: postgresRepo}
	cartRepo := &repository.CartRepo{PostgresRepository: postgresRepo}

	// Initialize services
	stripeSvc := service.NewStripeService(cfg.StripeKey)
	jwtSvc := service.NewJWTService(cfg.JWTSecret)
	ecommerceSvc := service.NewECommerceService(userRepo, productRepo, cartRepo, stripeSvc)

	// Initialize handlers
	apiHandler := &handler.APIHandler{
		Service:    ecommerceSvc,
		JWTService: jwtSvc,
	}

	// Setup routes
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/api/signup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		apiHandler.SignupHandler(w, r)
	})
	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		apiHandler.LoginHandler(w, r)
	})
	mux.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		apiHandler.GetProductsHandler(w, r)
	})

	// Authenticated routes (user)
	mux.HandleFunc("/api/cart/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.AuthMiddleware(jwtSvc, apiHandler.AddToCartHandler, false)(w, r)
	})
	mux.HandleFunc("/api/cart", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.AuthMiddleware(jwtSvc, apiHandler.ViewCartHandler, false)(w, r)
	})
	mux.HandleFunc("/api/checkout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.AuthMiddleware(jwtSvc, apiHandler.CheckoutHandler, false)(w, r)
	})

	// Admin routes
	mux.HandleFunc("/api/admin/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler.AuthMiddleware(jwtSvc, apiHandler.CreateProductHandler, true)(w, r)
	})

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
