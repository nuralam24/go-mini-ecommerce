package handlers

import (
	"net/http"

	"go-ecommerce/internal/database/sqlc"
	apierrors "go-ecommerce/internal/errors"
	"go-ecommerce/internal/logger"
	"go-ecommerce/internal/middleware"
	"go-ecommerce/internal/models"
	"go-ecommerce/internal/utils"
	"go-ecommerce/internal/validator"
)

type UserHandler struct {
	store *sqlc.Store
}

func NewUserHandler(store *sqlc.Store) *UserHandler {
	return &UserHandler{store: store}
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User registration data"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /api/v1/users/register [post]
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.New(apierrors.ErrCodeInvalidRequest, "Invalid request body"))
		return
	}

	if err := validator.Validate(req); err != nil {
		validationErrors := validator.FormatValidationErrors(err)
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.NewWithDetails(apierrors.ErrCodeValidationFailed, "Validation failed", validationErrors))
		return
	}

	existingUser, _ := h.store.GetUserByEmail(r.Context(), req.Email)
	if existingUser != nil {
		apierrors.RespondWithError(w, http.StatusConflict, apierrors.New(apierrors.ErrCodeConflict, "User with this email already exists"))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to hash password")
		apierrors.RespondWithError(w, http.StatusInternalServerError, apierrors.New(apierrors.ErrCodeInternalError, "Failed to hash password"))
		return
	}

	user, err := h.store.CreateUser(r.Context(), req.Email, hashedPassword, req.Name, req.Phone, req.Address)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to create user")
		apierrors.RespondWithError(w, http.StatusInternalServerError, apierrors.New(apierrors.ErrCodeDatabaseError, "Failed to create user"))
		return
	}

	logger.Log.Info().Str("email", req.Email).Str("user_id", user.ID).Msg("User registered successfully")
	apierrors.RespondWithJSON(w, http.StatusCreated, models.ToUserResponse(user))
}

// LoginUser godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 401 {object} map[string]string
// @Router /api/v1/users/login [post]
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.store.GetUserByEmail(r.Context(), req.Email)
	if err != nil || user == nil {
		apierrors.RespondWithError(w, http.StatusUnauthorized, apierrors.New(apierrors.ErrCodeUnauthorized, "Invalid email or password"))
		return
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		apierrors.RespondWithError(w, http.StatusUnauthorized, apierrors.New(apierrors.ErrCodeUnauthorized, "Invalid email or password"))
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email, "user")
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to generate token")
		apierrors.RespondWithError(w, http.StatusInternalServerError, apierrors.New(apierrors.ErrCodeInternalError, "Failed to generate token"))
		return
	}

	logger.Log.Info().Str("email", req.Email).Str("user_id", user.ID).Msg("User logged in")
	apierrors.RespondWithJSON(w, http.StatusOK, models.LoginResponse{
		Token: token,
		User:  models.ToUserResponse(user),
	})
}

// GetProfile godoc
// @Summary Get user profile
// @Description Get current user's profile information
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.UserResponse
// @Failure 401 {object} map[string]string
// @Router /api/v1/users/profile [get]
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		apierrors.RespondWithError(w, http.StatusUnauthorized, apierrors.New(apierrors.ErrCodeUnauthorized, "Unauthorized"))
		return
	}

	user, err := h.store.GetUserByID(r.Context(), userID)
	if err != nil || user == nil {
		apierrors.RespondWithError(w, http.StatusNotFound, apierrors.New(apierrors.ErrCodeNotFound, "User not found"))
		return
	}

	apierrors.RespondWithJSON(w, http.StatusOK, models.ToUserResponse(user))
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update current user's profile information
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param user body map[string]interface{} true "User update data"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/users/profile [put]
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		apierrors.RespondWithError(w, http.StatusUnauthorized, apierrors.New(apierrors.ErrCodeUnauthorized, "Unauthorized"))
		return
	}

	var updateData map[string]interface{}
	if err := utils.DecodeJSON(r, &updateData); err != nil {
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.New(apierrors.ErrCodeInvalidRequest, "Invalid request body"))
		return
	}

	var name *string
	var phone, address *string
	if n, ok := updateData["name"].(string); ok && n != "" {
		name = &n
	}
	if p, ok := updateData["phone"].(string); ok {
		phone = &p
	}
	if a, ok := updateData["address"].(string); ok {
		address = &a
	}

	if name == nil && phone == nil && address == nil {
		apierrors.RespondWithError(w, http.StatusBadRequest, apierrors.New(apierrors.ErrCodeInvalidRequest, "No valid fields to update"))
		return
	}

	user, err := h.store.UpdateUser(r.Context(), userID, name, phone, address)
	if err != nil {
		logger.Log.Error().Err(err).Str("user_id", userID).Msg("Failed to update user profile")
		apierrors.RespondWithError(w, http.StatusInternalServerError, apierrors.New(apierrors.ErrCodeDatabaseError, "Failed to update profile"))
		return
	}
	if user == nil {
		apierrors.RespondWithError(w, http.StatusNotFound, apierrors.New(apierrors.ErrCodeNotFound, "User not found"))
		return
	}

	logger.Log.Info().Str("user_id", userID).Msg("User profile updated")
	apierrors.RespondWithJSON(w, http.StatusOK, models.ToUserResponse(user))
}
