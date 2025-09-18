package router

import (
	"net/http"

	_ "github.com/Temisaputra/warOnk/docs" // wajib untuk register doc
	"github.com/Temisaputra/warOnk/internal/delivery/handler"
	"github.com/Temisaputra/warOnk/internal/delivery/middleware"
	"github.com/Temisaputra/warOnk/internal/infrastructure/config"
	"github.com/Temisaputra/warOnk/pkg/auth"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type Handlers struct {
	SupplierHandler *handler.SupplierHandler
	GroupHandler    *handler.GroupHandler
	UserHandler     *handler.UserHandler
	AuthHandler     *handler.AuthHandler
	ApprovalHandler *handler.ApprovalHandler
	Logger          *zap.Logger
	JwtService      auth.JwtService
}

// NewRouter bikin router dan register semua endpoint
func NewRouter(handlers *Handlers) http.Handler {
	router := mux.NewRouter()
	config := config.Get()
	// contoh custom metric
	opsProcessed := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "my_app_processed_ops_total",
		Help: "Total number of processed events",
	})

	prometheus.MustRegister(opsProcessed)

	authMW := middleware.NewAuthMiddleware(handlers.JwtService, opsProcessed)

	router.Handle("/metrics", promhttp.Handler())
	router.Use(middleware.LoggingMiddleware(handlers.Logger, config)) // <- inject logger

	// Swagger endpoint
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	api := router.PathPrefix("/api/war-onk").Subrouter()

	// ---------------- Public ----------------
	// Auth endpoints
	api.HandleFunc("/register", handlers.AuthHandler.Register).Methods("POST")
	api.HandleFunc("/login", handlers.AuthHandler.Login).Methods("POST")

	// Supplier endpoints
	api.HandleFunc("/suppliers", handlers.SupplierHandler.GetAllSupplier).Methods("GET")
	api.HandleFunc("/supplier/{id}", handlers.SupplierHandler.GetSupplierByID).Methods("GET")
	api.HandleFunc("/supplier-create", handlers.SupplierHandler.CreateSupplier).Methods("POST")
	api.HandleFunc("/supplier-update/{id}", handlers.SupplierHandler.UpdateSupplier).Methods("PUT")
	api.HandleFunc("/supplier-delete/{id}", handlers.SupplierHandler.DeleteSupplier).Methods("DELETE")
	api.HandleFunc("/supplier-block-unblock/{id}", handlers.SupplierHandler.BlockUnblockSupplier).Methods("PATCH")
	api.HandleFunc("/supplier-export", handlers.SupplierHandler.ExportSupplier).Methods("GET")

	// Group endpoints
	api.HandleFunc("/groups", handlers.GroupHandler.GetAllGroup).Methods("GET")
	api.HandleFunc("/group/{id}", handlers.GroupHandler.GetGroupByID).Methods("GET")
	api.HandleFunc("/group-create", handlers.GroupHandler.CreateGroup).Methods("POST")
	api.HandleFunc("/group-update/{id}", handlers.GroupHandler.UpdateGroup).Methods("PUT")
	api.HandleFunc("/group-delete/{id}", handlers.GroupHandler.DeleteGroup).Methods("DELETE")

	// ---------------- Protected ----------------
	protected := api.PathPrefix("").Subrouter()
	protected.Use(authMW.Authorization)

	// User endpoints
	protected.HandleFunc("/users", handlers.UserHandler.GetAllUser).Methods("GET")
	protected.HandleFunc("/user/{id}", handlers.UserHandler.GetUserByID).Methods("GET")
	protected.HandleFunc("/user-create", handlers.UserHandler.CreateUser).Methods("POST")
	protected.HandleFunc("/user-update/{id}", handlers.UserHandler.UpdateUser).Methods("PUT")
	protected.HandleFunc("/user-delete/{id}", handlers.UserHandler.DeleteUser).Methods("DELETE")

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"POST", "GET", "PUT", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Mode"},
		MaxAge:             60,
		AllowCredentials:   true,
		OptionsPassthrough: false,
	})

	return c.Handler(router)
}
