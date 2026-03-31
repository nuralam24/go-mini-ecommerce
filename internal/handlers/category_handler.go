package handlers

import (
	"net/http"

	"go-ecommerce/internal/database/sqlc"
	apierrors "go-ecommerce/internal/errors"
	"go-ecommerce/internal/logger"
	"go-ecommerce/internal/models"
	"go-ecommerce/internal/utils"
	"go-ecommerce/internal/validator"
)

type CategoryHandler struct {
	store *sqlc.Store
}

func NewCategoryHandler(store *sqlc.Store) *CategoryHandler {
	return &CategoryHandler{store: store}
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
// @Router /api/v1/categories [post]
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateCategoryRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.New(apierrors.ErrCodeInvalidRequest, "Invalid request body"))
		return
	}

	if err := validator.Validate(req); err != nil {
		validationErrors := validator.FormatValidationErrors(err)
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.NewWithDetails(apierrors.ErrCodeValidationFailed, "Validation failed", validationErrors))
		return
	}

	existing, _ := h.store.GetCategoryByName(r.Context(), req.Name)
	if existing != nil {
		apierrors.RespondWithError(w, http.StatusConflict, apierrors.New(apierrors.ErrCodeConflict, "Category with this name already exists"))
		return
	}

	category, err := h.store.CreateCategory(r.Context(), req.Name, req.Description)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to create category")
		apierrors.RespondWithError(w, http.StatusInternalServerError, apierrors.New(apierrors.ErrCodeDatabaseError, "Failed to create category"))
		return
	}

	logger.Log.Info().Str("category_id", category.ID).Str("name", category.Name).Msg("Category created")
	apierrors.RespondWithJSON(w, http.StatusCreated, models.ToCategoryResponse(category))
}

// GetCategories godoc
// @Summary Get all categories
// @Description Retrieve all product categories
// @Tags categories
// @Produce json
// @Success 200 {array} models.CategoryResponse
// @Router /api/v1/categories [get]
func (h *CategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	categories, err := h.store.ListCategories(r.Context())
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to fetch categories")
		apierrors.RespondWithError(w, http.StatusInternalServerError, apierrors.New(apierrors.ErrCodeDatabaseError, "Failed to fetch categories"))
		return
	}

	responses := make([]models.CategoryResponse, len(categories))
	for i := range categories {
		responses[i] = models.ToCategoryResponse(&categories[i])
	}

	apierrors.RespondWithJSON(w, http.StatusOK, responses)
}

// GetCategory godoc
// @Summary Get category by ID
// @Description Retrieve a specific category by ID
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} models.CategoryResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/categories/{id} [get]
func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.New(apierrors.ErrCodeInvalidRequest, "Category ID is required"))
		return
	}

	category, err := h.store.GetCategoryByID(r.Context(), id)
	if err != nil || category == nil {
		apierrors.RespondWithError(w, http.StatusNotFound, apierrors.New(apierrors.ErrCodeNotFound, "Category not found"))
		return
	}

	apierrors.RespondWithJSON(w, http.StatusOK, models.ToCategoryResponse(category))
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
// @Router /api/v1/categories/{id} [put]
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.New(apierrors.ErrCodeInvalidRequest, "Category ID is required"))
		return
	}

	var req models.UpdateCategoryRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.New(apierrors.ErrCodeInvalidRequest, "Invalid request body"))
		return
	}

	if req.Name == nil && req.Description == nil {
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.New(apierrors.ErrCodeInvalidRequest, "No valid fields to update"))
		return
	}

	category, err := h.store.UpdateCategory(r.Context(), id, req.Name, req.Description)
	if err != nil {
		logger.Log.Error().Err(err).Str("category_id", id).Msg("Failed to update category")
		apierrors.RespondWithError(w, http.StatusInternalServerError, apierrors.New(apierrors.ErrCodeDatabaseError, "Failed to update category"))
		return
	}
	if category == nil {
		apierrors.RespondWithError(w, http.StatusNotFound, apierrors.New(apierrors.ErrCodeNotFound, "Category not found"))
		return
	}

	logger.Log.Info().Str("category_id", id).Msg("Category updated")
	apierrors.RespondWithJSON(w, http.StatusOK, models.ToCategoryResponse(category))
}

// DeleteCategory godoc
// @Summary Delete category
// @Description Delete a category (Admin only)
// @Tags categories
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/categories/{id} [delete]
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.New(apierrors.ErrCodeInvalidRequest, "Category ID is required"))
		return
	}

	err := h.store.DeleteCategory(r.Context(), id)
	if err != nil {
		apierrors.RespondWithError(w, http.StatusNotFound, apierrors.New(apierrors.ErrCodeNotFound, "Category not found"))
		return
	}

	logger.Log.Info().Str("category_id", id).Msg("Category deleted")
	apierrors.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Category deleted successfully"})
}
