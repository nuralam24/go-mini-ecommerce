package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-ecommerce/internal/config"
	"go-ecommerce/internal/database"
	"go-ecommerce/internal/logger"
	"go-ecommerce/internal/middleware"
	"go-ecommerce/internal/router"
	"go-ecommerce/internal/utils"
	"go-ecommerce/internal/validator"

	_ "go-ecommerce/swagger"
)

// @title Go E-Commerce API
// @version 1.0
// @description A mini e-commerce API built with Go, PostgreSQL, and sqlc
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	cfg := config.Load()

	logger.Init(cfg.Env)
	logger.Log.Info().Str("port", cfg.Port).Msg("Starting Go E-Commerce API")

	utils.InitJWT(cfg.JWTSecret)
	validator.Init()

	store, err := database.Connect(cfg)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Disconnect()

	r := router.NewRouter(store)
	r.RegisterRoutes()

	handler := middleware.CORS(middleware.Logging(r))

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Log.Info().Str("port", cfg.Port).Msg("Server listening")
		logger.Log.Info().Str("swagger_url", "http://localhost:"+cfg.Port+"/swagger/index.html").Msg("Swagger UI available")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	logger.Log.Info().Msg("Server exited gracefully")
}
