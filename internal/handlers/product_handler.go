package handlers

import (
	"net/http"

	"go-ecommerce/internal/database"
	"go-ecommerce/internal/models"
	"go-ecommerce/internal/utils"
)

type ProductHandler struct{}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product (Admin only)
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param product body models.CreateProductRequest true "Product data"
// @Success 201 {object} models.ProductResponse
// @Failure 400 {object} map[string]string
// @Router /api/products [post]
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProductRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	cat, _ := database.Queries.GetCategoryByID(r.Context(), req.CategoryID)
	if cat == nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	brand, _ := database.Queries.GetBrandByID(r.Context(), req.BrandID)
	if brand == nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid brand ID")
		return
	}

	product, err := database.Queries.CreateProduct(r.Context(), req.Name, req.Description, req.Price, int32(req.Stock), req.ImageURL, req.CategoryID, req.BrandID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	withDetails, _ := database.Queries.GetProductWithDetails(r.Context(), product.ID)
	if withDetails != nil {
		utils.RespondWithJSON(w, http.StatusCreated, models.ToProductResponseFromDetails(withDetails))
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, models.ToProductResponseFromProduct(product))
}

// GetProducts godoc
// @Summary Get all products
// @Description Retrieve all products with optional filters
// @Tags products
// @Produce json
// @Param category query string false "Filter by category ID"
// @Param brand query string false "Filter by brand ID"
// @Success 200 {array} models.ProductResponse
// @Router /api/products [get]
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categoryID := r.URL.Query().Get("category")
	brandID := r.URL.Query().Get("brand")
	var catPtr, brandPtr *string
	if categoryID != "" {
		catPtr = &categoryID
	}
	if brandID != "" {
		brandPtr = &brandID
	}

	products, err := database.Queries.ListProducts(r.Context(), catPtr, brandPtr)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	responses := make([]models.ProductResponse, len(products))
	for i := range products {
		responses[i] = models.ToProductResponseFromDetails(&products[i])
	}
	utils.RespondWithJSON(w, http.StatusOK, responses)
}

// GetProduct godoc
// @Summary Get product by ID
// @Description Retrieve a specific product by ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} models.ProductResponse
// @Failure 404 {object} map[string]string
// @Router /api/products/{id} [get]
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	product, err := database.Queries.GetProductWithDetails(r.Context(), id)
	if err != nil || product == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, models.ToProductResponseFromDetails(product))
}

// UpdateProduct godoc
// @Summary Update product
// @Description Update a product (Admin only)
// @Tags products
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body models.UpdateProductRequest true "Product update data"
// @Success 200 {object} models.ProductResponse
// @Failure 404 {object} map[string]string
// @Router /api/products/{id} [put]
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	var req models.UpdateProductRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	var name *string
	var description *string
	var price *float64
	var stock *int32
	var imageURL *string
	var categoryID, brandID *string
	if req.Name != nil && *req.Name != "" {
		name = req.Name
	}
	if req.Description != nil {
		description = req.Description
	}
	if req.Price != nil {
		price = req.Price
	}
	if req.Stock != nil {
		s := int32(*req.Stock)
		stock = &s
	}
	if req.ImageURL != nil {
		imageURL = req.ImageURL
	}
	if req.CategoryID != nil {
		categoryID = req.CategoryID
	}
	if req.BrandID != nil {
		brandID = req.BrandID
	}

	if name == nil && description == nil && price == nil && stock == nil && imageURL == nil && categoryID == nil && brandID == nil {
		utils.RespondWithError(w, http.StatusBadRequest, "No valid fields to update")
		return
	}

	product, err := database.Queries.UpdateProduct(r.Context(), id, name, description, price, stock, imageURL, categoryID, brandID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}
	if product == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	withDetails, _ := database.Queries.GetProductWithDetails(r.Context(), id)
	if withDetails != nil {
		utils.RespondWithJSON(w, http.StatusOK, models.ToProductResponseFromDetails(withDetails))
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, models.ToProductResponseFromProduct(product))
}

// DeleteProduct godoc
// @Summary Delete product
// @Description Delete a product (Admin only)
// @Tags products
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/products/{id} [delete]
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Product ID is required")
		return
	}

	err := database.Queries.DeleteProduct(r.Context(), id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}
