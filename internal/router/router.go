package router

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"go-ecommerce/internal/database/sqlc"
	"go-ecommerce/internal/handlers"
	"go-ecommerce/internal/middleware"
)

type Router struct {
	mux   *http.ServeMux
	store *sqlc.Store
}

func NewRouter(store *sqlc.Store) *Router {
	return &Router{
		mux:   http.NewServeMux(),
		store: store,
	}
}

func (r *Router) RegisterRoutes() {
	userHandler := handlers.NewUserHandler(r.store)
	adminHandler := handlers.NewAdminHandler(r.store)
	categoryHandler := handlers.NewCategoryHandler(r.store)
	brandHandler := handlers.NewBrandHandler(r.store)
	productHandler := handlers.NewProductHandler(r.store)
	orderHandler := handlers.NewOrderHandler(r.store)
	healthHandler := handlers.NewHealthHandler()

	authMiddleware := middleware.AuthMiddleware

	r.mux.HandleFunc("GET /health", healthHandler.Health)
	r.mux.HandleFunc("GET /ready", healthHandler.Ready)

	r.mux.HandleFunc("POST /api/v1/users/register", userHandler.Register)
	r.mux.HandleFunc("POST /api/v1/users/login", userHandler.Login)

	r.mux.HandleFunc("POST /api/v1/admin/register", adminHandler.Register)
	r.mux.HandleFunc("POST /api/v1/admin/login", adminHandler.Login)

	r.mux.HandleFunc("GET /api/v1/categories", categoryHandler.GetAll)
	r.mux.HandleFunc("GET /api/v1/categories/{id}", categoryHandler.GetByID)

	r.mux.HandleFunc("GET /api/v1/brands", brandHandler.GetAll)
	r.mux.HandleFunc("GET /api/v1/brands/{id}", brandHandler.GetByID)

	r.mux.HandleFunc("GET /api/v1/products", productHandler.GetAll)
	r.mux.HandleFunc("GET /api/v1/products/{id}", productHandler.GetByID)

	r.mux.HandleFunc("GET /api/v1/users/profile", authMiddleware(http.HandlerFunc(userHandler.GetProfile)).ServeHTTP)
	r.mux.HandleFunc("PUT /api/v1/users/profile", authMiddleware(http.HandlerFunc(userHandler.UpdateProfile)).ServeHTTP)

	r.mux.HandleFunc("POST /api/v1/categories", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(categoryHandler.Create))).ServeHTTP)
	r.mux.HandleFunc("PUT /api/v1/categories/{id}", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(categoryHandler.Update))).ServeHTTP)
	r.mux.HandleFunc("DELETE /api/v1/categories/{id}", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(categoryHandler.Delete))).ServeHTTP)

	r.mux.HandleFunc("POST /api/v1/brands", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(brandHandler.Create))).ServeHTTP)
	r.mux.HandleFunc("PUT /api/v1/brands/{id}", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(brandHandler.Update))).ServeHTTP)
	r.mux.HandleFunc("DELETE /api/v1/brands/{id}", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(brandHandler.Delete))).ServeHTTP)

	r.mux.HandleFunc("POST /api/v1/products", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(productHandler.Create))).ServeHTTP)
	r.mux.HandleFunc("PUT /api/v1/products/{id}", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(productHandler.Update))).ServeHTTP)
	r.mux.HandleFunc("DELETE /api/v1/products/{id}", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(productHandler.Delete))).ServeHTTP)

	r.mux.HandleFunc("POST /api/v1/orders", authMiddleware(http.HandlerFunc(orderHandler.Create)).ServeHTTP)
	r.mux.HandleFunc("GET /api/v1/orders", authMiddleware(http.HandlerFunc(orderHandler.GetAll)).ServeHTTP)
	r.mux.HandleFunc("GET /api/v1/orders/{id}", authMiddleware(http.HandlerFunc(orderHandler.GetByID)).ServeHTTP)
	r.mux.HandleFunc("PUT /api/v1/orders/{id}/status", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(orderHandler.UpdateStatus))).ServeHTTP)

	r.mux.HandleFunc("POST /api/users/register", userHandler.Register)
	r.mux.HandleFunc("POST /api/users/login", userHandler.Login)
	r.mux.HandleFunc("POST /api/admin/register", adminHandler.Register)
	r.mux.HandleFunc("POST /api/admin/login", adminHandler.Login)
	r.mux.HandleFunc("GET /api/categories", categoryHandler.GetAll)
	r.mux.HandleFunc("GET /api/categories/{id}", categoryHandler.GetByID)
	r.mux.HandleFunc("GET /api/brands", brandHandler.GetAll)
	r.mux.HandleFunc("GET /api/brands/{id}", brandHandler.GetByID)
	r.mux.HandleFunc("GET /api/products", productHandler.GetAll)
	r.mux.HandleFunc("GET /api/products/{id}", productHandler.GetByID)
	r.mux.HandleFunc("GET /api/users/profile", authMiddleware(http.HandlerFunc(userHandler.GetProfile)).ServeHTTP)
	r.mux.HandleFunc("PUT /api/users/profile", authMiddleware(http.HandlerFunc(userHandler.UpdateProfile)).ServeHTTP)
	r.mux.HandleFunc("POST /api/categories", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(categoryHandler.Create))).ServeHTTP)
	r.mux.HandleFunc("PUT /api/categories/{id}", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(categoryHandler.Update))).ServeHTTP)
	r.mux.HandleFunc("DELETE /api/categories/{id}", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(categoryHandler.Delete))).ServeHTTP)
	r.mux.HandleFunc("POST /api/brands", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(brandHandler.Create))).ServeHTTP)
	r.mux.HandleFunc("PUT /api/brands/{id}", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(brandHandler.Update))).ServeHTTP)
	r.mux.HandleFunc("DELETE /api/brands/{id}", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(brandHandler.Delete))).ServeHTTP)
	r.mux.HandleFunc("POST /api/products", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(productHandler.Create))).ServeHTTP)
	r.mux.HandleFunc("PUT /api/products/{id}", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(productHandler.Update))).ServeHTTP)
	r.mux.HandleFunc("DELETE /api/products/{id}", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(productHandler.Delete))).ServeHTTP)
	r.mux.HandleFunc("POST /api/orders", authMiddleware(http.HandlerFunc(orderHandler.Create)).ServeHTTP)
	r.mux.HandleFunc("GET /api/orders", authMiddleware(http.HandlerFunc(orderHandler.GetAll)).ServeHTTP)
	r.mux.HandleFunc("GET /api/orders/{id}", authMiddleware(http.HandlerFunc(orderHandler.GetByID)).ServeHTTP)
	r.mux.HandleFunc("PUT /api/orders/{id}/status", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(orderHandler.UpdateStatus))).ServeHTTP)

	r.mux.Handle("/swagger/", httpSwagger.WrapHandler)
	r.mux.Handle("/swagger/index.html", httpSwagger.WrapHandler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
