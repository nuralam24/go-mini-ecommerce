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

type BrandHandler struct {
	store *sqlc.Store
}

func NewBrandHandler(store *sqlc.Store) *BrandHandler {
	return &BrandHandler{store: store}
}

// CreateBrand godoc
// @Summary Create a new brand
// @Description Create a new product brand (Admin only)
// @Tags brands
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param brand body models.CreateBrandRequest true "Brand data"
// @Success 201 {object} models.BrandResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/v1/brands [post]
func (h *BrandHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateBrandRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.New(apierrors.ErrCodeInvalidRequest, "Invalid request body"))
		return
	}

	if err := validator.Validate(req); err != nil {
		validationErrors := validator.FormatValidationErrors(err)
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.NewWithDetails(apierrors.ErrCodeValidationFailed, "Validation failed", validationErrors))
		return
	}

	existing, _ := h.store.GetBrandByName(r.Context(), req.Name)
	if existing != nil {
		apierrors.RespondWithError(w, http.StatusConflict, apierrors.New(apierrors.ErrCodeConflict, "Brand with this name already exists"))
		return
	}

	brand, err := h.store.CreateBrand(r.Context(), req.Name, req.Description)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to create brand")
		apierrors.RespondWithError(w, http.StatusInternalServerError, apierrors.New(apierrors.ErrCodeDatabaseError, "Failed to create brand"))
		return
	}

	logger.Log.Info().Str("brand_id", brand.ID).Str("name", brand.Name).Msg("Brand created")
	apierrors.RespondWithJSON(w, http.StatusCreated, models.ToBrandResponse(brand))
}

// GetBrands godoc
// @Summary Get all brands
// @Description Retrieve all product brands
// @Tags brands
// @Produce json
// @Success 200 {array} models.BrandResponse
// @Router /api/v1/brands [get]
func (h *BrandHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	brands, err := h.store.ListBrands(r.Context())
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to fetch brands")
		apierrors.RespondWithError(w, http.StatusInternalServerError, apierrors.New(apierrors.ErrCodeDatabaseError, "Failed to fetch brands"))
		return
	}

	responses := make([]models.BrandResponse, len(brands))
	for i := range brands {
		responses[i] = models.ToBrandResponse(&brands[i])
	}

	apierrors.RespondWithJSON(w, http.StatusOK, responses)
}

// GetBrand godoc
// @Summary Get brand by ID
// @Description Retrieve a specific brand by ID
// @Tags brands
// @Produce json
// @Param id path string true "Brand ID"
// @Success 200 {object} models.BrandResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/brands/{id} [get]
func (h *BrandHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.New(apierrors.ErrCodeInvalidRequest, "Brand ID is required"))
		return
	}

	brand, err := h.store.GetBrandByID(r.Context(), id)
	if err != nil || brand == nil {
		apierrors.RespondWithError(w, http.StatusNotFound, apierrors.New(apierrors.ErrCodeNotFound, "Brand not found"))
		return
	}

	apierrors.RespondWithJSON(w, http.StatusOK, models.ToBrandResponse(brand))
}

// UpdateBrand godoc
// @Summary Update brand
// @Description Update a brand (Admin only)
// @Tags brands
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Brand ID"
// @Param brand body models.UpdateBrandRequest true "Brand update data"
// @Success 200 {object} models.BrandResponse
// @Failure 404 {object} map[string]string
// @Router /api/v1/brands/{id} [put]
func (h *BrandHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.New(apierrors.ErrCodeInvalidRequest, "Brand ID is required"))
		return
	}

	var req models.UpdateBrandRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.New(apierrors.ErrCodeInvalidRequest, "Invalid request body"))
		return
	}

	if req.Name == nil && req.Description == nil {
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.New(apierrors.ErrCodeInvalidRequest, "No valid fields to update"))
		return
	}

	brand, err := h.store.UpdateBrand(r.Context(), id, req.Name, req.Description)
	if err != nil {
		logger.Log.Error().Err(err).Str("brand_id", id).Msg("Failed to update brand")
		apierrors.RespondWithError(w, http.StatusInternalServerError, apierrors.New(apierrors.ErrCodeDatabaseError, "Failed to update brand"))
		return
	}
	if brand == nil {
		apierrors.RespondWithError(w, http.StatusNotFound, apierrors.New(apierrors.ErrCodeNotFound, "Brand not found"))
		return
	}

	logger.Log.Info().Str("brand_id", id).Msg("Brand updated")
	apierrors.RespondWithJSON(w, http.StatusOK, models.ToBrandResponse(brand))
}

// DeleteBrand godoc
// @Summary Delete brand
// @Description Delete a brand (Admin only)
// @Tags brands
// @Security BearerAuth
// @Param id path string true "Brand ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/brands/{id} [delete]
func (h *BrandHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.New(apierrors.ErrCodeInvalidRequest, "Brand ID is required"))
		return
	}

	err := h.store.DeleteBrand(r.Context(), id)
	if err != nil {
		apierrors.RespondWithError(w, http.StatusNotFound, apierrors.New(apierrors.ErrCodeNotFound, "Brand not found"))
		return
	}

	logger.Log.Info().Str("brand_id", id).Msg("Brand deleted")
	apierrors.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Brand deleted successfully"})
}
