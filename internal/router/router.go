package router

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"go-ecommerce/internal/handlers"
	"go-ecommerce/internal/middleware"
)

type Router struct {
	mux *http.ServeMux
}

func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

func (r *Router) RegisterRoutes() {
	// Initialize handlers
	userHandler := handlers.NewUserHandler()
	adminHandler := handlers.NewAdminHandler()
	categoryHandler := handlers.NewCategoryHandler()
	brandHandler := handlers.NewBrandHandler()
	productHandler := handlers.NewProductHandler()
	orderHandler := handlers.NewOrderHandler()

	// Public routes - User
	r.mux.HandleFunc("POST /api/users/register", userHandler.Register)
	r.mux.HandleFunc("POST /api/users/login", userHandler.Login)

	// Public routes - Admin
	r.mux.HandleFunc("POST /api/admin/register", adminHandler.Register)
	r.mux.HandleFunc("POST /api/admin/login", adminHandler.Login)

	// Public routes - Categories
	r.mux.HandleFunc("GET /api/categories", categoryHandler.GetAll)
	r.mux.HandleFunc("GET /api/categories/{id}", categoryHandler.GetByID)

	// Public routes - Brands
	r.mux.HandleFunc("GET /api/brands", brandHandler.GetAll)
	r.mux.HandleFunc("GET /api/brands/{id}", brandHandler.GetByID)

	// Public routes - Products
	r.mux.HandleFunc("GET /api/products", productHandler.GetAll)
	r.mux.HandleFunc("GET /api/products/{id}", productHandler.GetByID)

	// Protected routes (require authentication)
	authMiddleware := middleware.AuthMiddleware

	// User protected routes
	r.mux.HandleFunc("GET /api/users/profile", authMiddleware(http.HandlerFunc(userHandler.GetProfile)).ServeHTTP)
	r.mux.HandleFunc("PUT /api/users/profile", authMiddleware(http.HandlerFunc(userHandler.UpdateProfile)).ServeHTTP)

	// Category protected routes (admin only)
	r.mux.HandleFunc("POST /api/categories", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(categoryHandler.Create))).ServeHTTP)
	r.mux.HandleFunc("PUT /api/categories/{id}", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(categoryHandler.Update))).ServeHTTP)
	r.mux.HandleFunc("DELETE /api/categories/{id}", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(categoryHandler.Delete))).ServeHTTP)

	// Brand protected routes (admin only)
	r.mux.HandleFunc("POST /api/brands", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(brandHandler.Create))).ServeHTTP)
	r.mux.HandleFunc("PUT /api/brands/{id}", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(brandHandler.Update))).ServeHTTP)
	r.mux.HandleFunc("DELETE /api/brands/{id}", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(brandHandler.Delete))).ServeHTTP)

	// Product protected routes (admin only)
	r.mux.HandleFunc("POST /api/products", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(productHandler.Create))).ServeHTTP)
	r.mux.HandleFunc("PUT /api/products/{id}", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(productHandler.Update))).ServeHTTP)
	r.mux.HandleFunc("DELETE /api/products/{id}", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(productHandler.Delete))).ServeHTTP)

	// Order protected routes
	r.mux.HandleFunc("POST /api/orders", authMiddleware(http.HandlerFunc(orderHandler.Create)).ServeHTTP)
	r.mux.HandleFunc("GET /api/orders", authMiddleware(http.HandlerFunc(orderHandler.GetAll)).ServeHTTP)
	r.mux.HandleFunc("GET /api/orders/{id}", authMiddleware(http.HandlerFunc(orderHandler.GetByID)).ServeHTTP)
	// Order status update requires admin
	r.mux.HandleFunc("PUT /api/orders/{id}/status", authMiddleware(middleware.AdminMiddleware(http.HandlerFunc(orderHandler.UpdateStatus))).ServeHTTP)

	// Swagger docs (served by http-swagger)
	// After running swagger generation (scripts/generate-swagger.sh) this will serve the UI.
	r.mux.Handle("/swagger/", httpSwagger.WrapHandler)
	r.mux.Handle("/swagger/index.html", httpSwagger.WrapHandler)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
