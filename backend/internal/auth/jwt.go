package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("token has expired")
	ErrTokenNotFound    = errors.New("token not found")
	ErrInvalidTokenType = errors.New("invalid token type")
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// Claims represents JWT claims
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Type     string `json:"type"` // access or refresh
	jwt.RegisteredClaims
}

// TokenPair holds access and refresh tokens
type TokenPair struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	TokenType    string    `json:"token_type"`
}

// JWTManager handles JWT operations
type JWTManager struct {
	secretKey       string
	accessExpiry    time.Duration
	refreshExpiry   time.Duration
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secretKey string, accessExpiry, refreshExpiry time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// GenerateTokenPair generates both access and refresh tokens
func (m *JWTManager) GenerateTokenPair(userID uint, username, role string) (*TokenPair, error) {
	// Generate access token
	accessToken, accessExp, err := m.generateToken(userID, username, role, AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, _, err := m.generateToken(userID, username, role, RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    accessExp,
		TokenType:    "Bearer",
	}, nil
}

// generateToken creates a JWT token
func (m *JWTManager) generateToken(userID uint, username, role string, tokenType TokenType) (string, time.Time, error) {
	var expiresAt time.Time
	
	if tokenType == AccessToken {
		expiresAt = time.Now().Add(m.accessExpiry)
	} else {
		expiresAt = time.Now().Add(m.refreshExpiry)
	}

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		Type:     string(tokenType),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "dental-marketplace",
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// ValidateToken validates a JWT token and returns claims
func (m *JWTManager) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			// Verify signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(m.secretKey), nil
		},
	)

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ValidateAccessToken specifically validates an access token
func (m *JWTManager) ValidateAccessToken(tokenString string) (*Claims, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Type != string(AccessToken) {
		return nil, ErrInvalidTokenType
	}

	return claims, nil
}

// ValidateRefreshToken specifically validates a refresh token
func (m *JWTManager) ValidateRefreshToken(tokenString string) (*Claims, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims.Type != string(RefreshToken) {
		return nil, ErrInvalidTokenType
	}

	return claims, nil
}

// RefreshAccessToken generates a new access token from a refresh token
func (m *JWTManager) RefreshAccessToken(refreshToken string) (string, time.Time, error) {
	claims, err := m.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", time.Time{}, err
	}

	// Generate new access token
	accessToken, expiresAt, err := m.generateToken(claims.UserID, claims.Username, claims.Role, AccessToken)
	if err != nil {
		return "", time.Time{}, err
	}

	return accessToken, expiresAt, nil
}

// ExtractUserID extracts user ID from token
func (m *JWTManager) ExtractUserID(tokenString string) (uint, error) {
	claims, err := m.ValidateAccessToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// ExtractRole extracts user role from token
func (m *JWTManager) ExtractRole(tokenString string) (string, error) {
	claims, err := m.ValidateAccessToken(tokenString)
	if err != nil {
		return "", err
	}
	return claims.Role, nil
}
