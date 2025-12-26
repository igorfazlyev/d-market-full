package handlers

import (
	"dental-marketplace/backend/internal/auth"
	"dental-marketplace/backend/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	repo       *repository.Repository
	jwtManager *auth.JWTManager
}

func NewAuthHandler(repo *repository.Repository, jwtManager *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		repo:       repo,
		jwtManager: jwtManager,
	}
}

// LoginRequest represents login request body
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	User         UserInfo          `json:"user"`
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	ExpiresAt    string            `json:"expires_at"`
	TokenType    string            `json:"token_type"`
}

// UserInfo represents user information in response
type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Profile  interface{} `json:"profile,omitempty"`
}

// Login handles user authentication
// @Summary Login user
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Authenticate user
	user, err := h.repo.AuthenticateUser(req.Username, req.Password)
	if err != nil {
		if err == repository.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid username or password",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Authentication failed",
		})
		return
	}

	// Generate token pair
	tokenPair, err := h.jwtManager.GenerateTokenPair(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate tokens",
		})
		return
	}

	// Prepare user info with profile
	userInfo := UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	// Add role-specific profile
	switch user.Role {
	case "patient":
		if user.Patient != nil {
			userInfo.Profile = user.Patient
		}
	case "clinic":
		if user.Clinic != nil {
			userInfo.Profile = user.Clinic
		}
	case "regulator":
		if user.Regulator != nil {
			userInfo.Profile = user.Regulator
		}
	}

	c.JSON(http.StatusOK, LoginResponse{
		User:         userInfo,
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    tokenPair.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
		TokenType:    tokenPair.TokenType,
	})
}

// RefreshRequest represents refresh token request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshResponse represents refresh token response
type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresAt   string `json:"expires_at"`
	TokenType   string `json:"token_type"`
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Generate new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh token"
// @Success 200 {object} RefreshResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Generate new access token
	accessToken, expiresAt, err := h.jwtManager.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid or expired refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, RefreshResponse{
		AccessToken: accessToken,
		ExpiresAt:   expiresAt.Format("2006-01-02T15:04:05Z07:00"),
		TokenType:   "Bearer",
	})
}

// GetMe returns current user information
// @Summary Get current user
// @Description Get authenticated user information
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserInfo
// @Failure 401 {object} ErrorResponse
// @Router /api/auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	user, err := h.repo.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	userInfo := UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}

	// Add role-specific profile
	switch user.Role {
	case "patient":
		if user.Patient != nil {
			userInfo.Profile = user.Patient
		}
	case "clinic":
		if user.Clinic != nil {
			userInfo.Profile = user.Clinic
		}
	case "regulator":
		if user.Regulator != nil {
			userInfo.Profile = user.Regulator
		}
	}

	c.JSON(http.StatusOK, userInfo)
}
