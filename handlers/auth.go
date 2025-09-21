package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hrusfandi/sb-task-management/models"
	"github.com/hrusfandi/sb-task-management/utils"
	"gorm.io/gorm"
)

type AuthHandler struct {
	userRepo models.UserRepository
}

func NewAuthHandler(userRepo models.UserRepository) *AuthHandler {
	return &AuthHandler{
		userRepo: userRepo,
	}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)

	// Validate name
	if valid, msg := utils.ValidateName(req.Name); !valid {
		utils.RespondError(w, http.StatusBadRequest, msg)
		return
	}

	// Validate email
	if !utils.ValidateEmail(req.Email) {
		utils.RespondError(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	// Validate password
	if valid, msg := utils.ValidatePassword(req.Password); !valid {
		utils.RespondError(w, http.StatusBadRequest, msg)
		return
	}

	existingUser, _ := h.userRepo.GetUserByEmail(req.Email)
	if existingUser != nil {
		utils.RespondError(w, http.StatusConflict, "Email already registered")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to process password")
		return
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := h.userRepo.CreateUser(user); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	userResponse := models.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	utils.RespondCreated(w, "User registered successfully", userResponse)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" || req.Password == "" {
		utils.RespondError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	user, err := h.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.RespondError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		utils.RespondError(w, http.StatusInternalServerError, "Failed to authenticate")
		return
	}

	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	userResponse := models.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	response := AuthResponse{
		Token: token,
		User:  userResponse,
	}

	utils.RespondSuccess(w, "Login successful", response)
}