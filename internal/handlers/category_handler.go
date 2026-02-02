package handlers

import (
	"net/http"

	"go-ecommerce/internal/database"
	"go-ecommerce/internal/models"
	"go-ecommerce/internal/utils"
)

type CategoryHandler struct{}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{}
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new product category (Admin only)
// @Tags categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param category body models.CreateCategoryRequest true "Category data"
// @Success 201 {object} models.CategoryResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/categories [post]
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCategoryRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	existing, _ := database.Queries.GetCategoryByName(r.Context(), req.Name)
	if existing != nil {
		utils.RespondWithError(w, http.StatusConflict, "Category with this name already exists")
		return
	}

	category, err := database.Queries.CreateCategory(r.Context(), req.Name, req.Description)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create category")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, models.ToCategoryResponse(category))
}

// GetCategories godoc
// @Summary Get all categories
// @Description Retrieve all product categories
// @Tags categories
// @Produce json
// @Success 200 {array} models.CategoryResponse
// @Router /api/categories [get]
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := database.Queries.ListCategories(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch categories")
		return
	}

	responses := make([]models.CategoryResponse, len(categories))
	for i := range categories {
		responses[i] = models.ToCategoryResponse(&categories[i])
	}

	utils.RespondWithJSON(w, http.StatusOK, responses)
}

// GetCategory godoc
// @Summary Get category by ID
// @Description Retrieve a specific category by ID
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} models.CategoryResponse
// @Failure 404 {object} map[string]string
// @Router /api/categories/{id} [get]
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Category ID is required")
		return
	}

	category, err := database.Queries.GetCategoryByID(r.Context(), id)
	if err != nil || category == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Category not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.ToCategoryResponse(category))
}

// UpdateCategory godoc
// @Summary Update category
// @Description Update a category (Admin only)
// @Tags categories
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body models.UpdateCategoryRequest true "Category update data"
// @Success 200 {object} models.CategoryResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/categories/{id} [put]
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Category ID is required")
		return
	}

	var req models.UpdateCategoryRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == nil && req.Description == nil {
		utils.RespondWithError(w, http.StatusBadRequest, "No valid fields to update")
		return
	}

	category, err := database.Queries.UpdateCategory(r.Context(), id, req.Name, req.Description)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update category")
		return
	}
	if category == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Category not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.ToCategoryResponse(category))
}

// DeleteCategory godoc
// @Summary Delete category
// @Description Delete a category (Admin only)
// @Tags categories
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/categories/{id} [delete]
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Category ID is required")
		return
	}

	err := database.Queries.DeleteCategory(r.Context(), id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Category not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Category deleted successfully"})
}
