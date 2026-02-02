package handlers

import (
	"net/http"
	"strings"

	"go-ecommerce/internal/database"
	"go-ecommerce/internal/models"
	"go-ecommerce/internal/utils"
)

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
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
// @Router /api/admin/register [post]
func (h *AdminHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}

	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if !strings.Contains(req.Email, "@") {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	existingAdmin, _ := database.Queries.GetAdminByEmail(r.Context(), req.Email)
	if existingAdmin != nil {
		utils.RespondWithError(w, http.StatusConflict, "Admin with this email already exists")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	_, err = database.Queries.CreateAdmin(r.Context(), req.Email, hashedPassword, req.Name)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create admin")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"message": "Admin created successfully"})
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
// @Router /api/admin/login [post]
func (h *AdminHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := utils.DecodeJSON(r, &req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	admin, err := database.Queries.GetAdminByEmail(r.Context(), req.Email)
	if err != nil || admin == nil {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if !utils.CheckPasswordHash(req.Password, admin.Password) {
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := utils.GenerateToken(admin.ID, admin.Email, "admin")
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.LoginResponse{
		Token: token,
		User: models.UserResponse{
			ID:    admin.ID,
			Email: admin.Email,
			Name:  admin.Name,
		},
	})
}
