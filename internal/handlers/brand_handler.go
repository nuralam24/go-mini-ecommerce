package handlers

import (
	"net/http"

	"go-ecommerce/internal/database"
	"go-ecommerce/internal/models"
	"go-ecommerce/internal/utils"
)

type BrandHandler struct{}

func NewBrandHandler() *BrandHandler {
	return &BrandHandler{}
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
// @Router /api/brands [post]
func (h *BrandHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateBrandRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	existing, _ := database.Queries.GetBrandByName(r.Context(), req.Name)
	if existing != nil {
		utils.RespondWithError(w, http.StatusConflict, "Brand with this name already exists")
		return
	}

	brand, err := database.Queries.CreateBrand(r.Context(), req.Name, req.Description)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create brand")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, models.ToBrandResponse(brand))
}

// GetBrands godoc
// @Summary Get all brands
// @Description Retrieve all product brands
// @Tags brands
// @Produce json
// @Success 200 {array} models.BrandResponse
// @Router /api/brands [get]
func (h *BrandHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	brands, err := database.Queries.ListBrands(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to fetch brands")
		return
	}

	responses := make([]models.BrandResponse, len(brands))
	for i := range brands {
		responses[i] = models.ToBrandResponse(&brands[i])
	}

	utils.RespondWithJSON(w, http.StatusOK, responses)
}

// GetBrand godoc
// @Summary Get brand by ID
// @Description Retrieve a specific brand by ID
// @Tags brands
// @Produce json
// @Param id path string true "Brand ID"
// @Success 200 {object} models.BrandResponse
// @Failure 404 {object} map[string]string
// @Router /api/brands/{id} [get]
func (h *BrandHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Brand ID is required")
		return
	}

	brand, err := database.Queries.GetBrandByID(r.Context(), id)
	if err != nil || brand == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Brand not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.ToBrandResponse(brand))
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
// @Router /api/brands/{id} [put]
func (h *BrandHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Brand ID is required")
		return
	}

	var req models.UpdateBrandRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == nil && req.Description == nil {
		utils.RespondWithError(w, http.StatusBadRequest, "No valid fields to update")
		return
	}

	brand, err := database.Queries.UpdateBrand(r.Context(), id, req.Name, req.Description)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update brand")
		return
	}
	if brand == nil {
		utils.RespondWithError(w, http.StatusNotFound, "Brand not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.ToBrandResponse(brand))
}

// DeleteBrand godoc
// @Summary Delete brand
// @Description Delete a brand (Admin only)
// @Tags brands
// @Security BearerAuth
// @Param id path string true "Brand ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/brands/{id} [delete]
func (h *BrandHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Brand ID is required")
		return
	}

	err := database.Queries.DeleteBrand(r.Context(), id)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Brand not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Brand deleted successfully"})
}
