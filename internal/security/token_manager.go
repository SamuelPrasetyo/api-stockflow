package security

import (
	"api-stockflow/internal/domain"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenManager handles JWT token operations
type TokenManager interface {
	GenerateAccessToken(userID string, role domain.UserRole) (string, error)
	GenerateRefreshToken(userID string) (string, error)
	VerifyAccessToken(tokenString string) (*TokenPayload, error)
	VerifyRefreshToken(tokenString string) (*TokenPayload, error)
}

// TokenPayload represents the JWT token payload
type TokenPayload struct {
	UserID string           `json:"user_id"`
	Role   domain.UserRole  `json:"role"`
	jwt.RegisteredClaims
}

type jwtTokenManager struct {
	accessTokenKey  string
	refreshTokenKey string
	accessTokenAge  int
}

// NewJWTTokenManager creates a new JWT token manager
func NewJWTTokenManager() TokenManager {
	accessTokenAge, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_AGE"))
	if accessTokenAge == 0 {
		accessTokenAge = 3600 // default 1 hour
	}

	return &jwtTokenManager{
		accessTokenKey:  os.Getenv("ACCESS_TOKEN_KEY"),
		refreshTokenKey: os.Getenv("REFRESH_TOKEN_KEY"),
		accessTokenAge:  accessTokenAge,
	}
}

func (t *jwtTokenManager) GenerateAccessToken(userID string, role domain.UserRole) (string, error) {
	claims := TokenPayload{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(t.accessTokenAge) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.accessTokenKey))
}

func (t *jwtTokenManager) GenerateRefreshToken(userID string) (string, error) {
	claims := TokenPayload{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 30)), // 30 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(t.refreshTokenKey))
}

func (t *jwtTokenManager) VerifyAccessToken(tokenString string) (*TokenPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenPayload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.accessTokenKey), nil
	})

	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	if claims, ok := token.Claims.(*TokenPayload); ok && token.Valid {
		return claims, nil
	}

	return nil, domain.ErrInvalidToken
}

func (t *jwtTokenManager) VerifyRefreshToken(tokenString string) (*TokenPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenPayload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(t.refreshTokenKey), nil
	})

	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	if claims, ok := token.Claims.(*TokenPayload); ok && token.Valid {
		return claims, nil
	}

	return nil, domain.ErrInvalidToken
}
