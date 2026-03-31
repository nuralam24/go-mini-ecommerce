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

type AdminHandler struct {
	store *sqlc.Store
}

func NewAdminHandler(store *sqlc.Store) *AdminHandler {
	return &AdminHandler{store: store}
}

// RegisterAdmin godoc
// @Summary Register a new admin
// @Description Create a new admin account (only for initial setup)
// @Tags admin
// @Accept json
// @Produce json
// @Param admin body map[string]string true "Admin registration data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/v1/admin/register [post]
func (h *AdminHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
		Name     string `json:"name" validate:"required"`
	}

	if err := utils.DecodeJSON(r, &req); err != nil {
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.New(apierrors.ErrCodeInvalidRequest, "Invalid request body"))
		return
	}

	if err := validator.Validate(req); err != nil {
		validationErrors := validator.FormatValidationErrors(err)
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.NewWithDetails(apierrors.ErrCodeValidationFailed, "Validation failed", validationErrors))
		return
	}

	existingAdmin, _ := h.store.GetAdminByEmail(r.Context(), req.Email)
	if existingAdmin != nil {
		apierrors.RespondWithError(w, http.StatusConflict, apierrors.New(apierrors.ErrCodeConflict, "Admin with this email already exists"))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to hash password")
		apierrors.RespondWithError(w, http.StatusInternalServerError, apierrors.New(apierrors.ErrCodeInternalError, "Failed to hash password"))
		return
	}

	_, err = h.store.CreateAdmin(r.Context(), req.Email, hashedPassword, req.Name)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to create admin")
		apierrors.RespondWithError(w, http.StatusInternalServerError, apierrors.New(apierrors.ErrCodeDatabaseError, "Failed to create admin"))
		return
	}

	logger.Log.Info().Str("email", req.Email).Msg("Admin created successfully")
	apierrors.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Admin created successfully"})
}

// LoginAdmin godoc
// @Summary Admin login
// @Description Authenticate admin and return JWT token
// @Tags admin
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 401 {object} map[string]string
// @Router /api/v1/admin/login [post]
func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.New(apierrors.ErrCodeInvalidRequest, "Invalid request body"))
		return
	}

	if err := validator.Validate(req); err != nil {
		validationErrors := validator.FormatValidationErrors(err)
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.NewWithDetails(apierrors.ErrCodeValidationFailed, "Validation failed", validationErrors))
		return
	}

	admin, err := h.store.GetAdminByEmail(r.Context(), req.Email)
	if err != nil || admin == nil {
		apierrors.RespondWithError(w, http.StatusUnauthorized, apierrors.New(apierrors.ErrCodeUnauthorized, "Invalid email or password"))
		return
	}

	if !utils.CheckPasswordHash(req.Password, admin.Password) {
		apierrors.RespondWithError(w, http.StatusUnauthorized, apierrors.New(apierrors.ErrCodeUnauthorized, "Invalid email or password"))
		return
	}

	token, err := utils.GenerateToken(admin.ID, admin.Email, "admin")
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to generate token")
		apierrors.RespondWithError(w, http.StatusInternalServerError, apierrors.New(apierrors.ErrCodeInternalError, "Failed to generate token"))
		return
	}

	logger.Log.Info().Str("email", req.Email).Str("role", "admin").Msg("Admin logged in")
	apierrors.RespondWithJSON(w, http.StatusOK, models.LoginResponse{
		Token: token,
		User: models.UserResponse{
			ID:    admin.ID,
			Email: admin.Email,
			Name:  admin.Name,
		},
	})
}
