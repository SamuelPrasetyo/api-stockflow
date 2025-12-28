package usecase

import (
	"api-stockflow/internal/domain"
	"api-stockflow/internal/repository"
	"api-stockflow/internal/security"
	"context"

	"github.com/google/uuid"
)

// AuthUseCase defines authentication use case interface
type AuthUseCase interface {
	Login(ctx context.Context, username, password string) (*LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	Register(ctx context.Context, req *RegisterRequest) error
}

// LoginResponse represents login response
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RefreshTokenResponse represents refresh token response
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Username string           `json:"username"`
	Password string           `json:"password"`
	Fullname string           `json:"fullname"`
	Role     domain.UserRole  `json:"role"`
}

type authUseCase struct {
	userRepo           repository.UserRepository
	authRepo           repository.AuthenticationRepository
	tokenManager       security.TokenManager
	passwordHash       security.PasswordHash
}

// NewAuthUseCase creates a new authentication use case
func NewAuthUseCase(
	userRepo repository.UserRepository,
	authRepo repository.AuthenticationRepository,
	tokenManager security.TokenManager,
	passwordHash security.PasswordHash,
) AuthUseCase {
	return &authUseCase{
		userRepo:     userRepo,
		authRepo:     authRepo,
		tokenManager: tokenManager,
		passwordHash: passwordHash,
	}
}

func (u *authUseCase) Login(ctx context.Context, username, password string) (*LoginResponse, error) {
	// Get user by username
	user, err := u.userRepo.GetByUsername(ctx, username)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	// Verify password
	err = u.passwordHash.Compare(user.Password, password)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	// Generate access token
	accessToken, err := u.tokenManager.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := u.tokenManager.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Store refresh token
	err = u.authRepo.AddToken(ctx, refreshToken, user.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *authUseCase) RefreshToken(ctx context.Context, refreshToken string) (*RefreshTokenResponse, error) {
	// Verify refresh token
	payload, err := u.tokenManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// Check if refresh token exists in database
	err = u.authRepo.CheckTokenAvailability(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	// Get user to get the role
	user, err := u.userRepo.GetByID(ctx, payload.UserID)
	if err != nil {
		return nil, err
	}

	// Generate new access token
	accessToken, err := u.tokenManager.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenResponse{
		AccessToken: accessToken,
	}, nil
}

func (u *authUseCase) Logout(ctx context.Context, refreshToken string) error {
	// Verify refresh token
	_, err := u.tokenManager.VerifyRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	// Delete refresh token from database
	err = u.authRepo.DeleteToken(ctx, refreshToken)
	if err != nil {
		return err
	}

	return nil
}

func (u *authUseCase) Register(ctx context.Context, req *RegisterRequest) error {
	// Generate user ID
	userID := "user-" + uuid.New().String()

	// Hash password
	hashedPassword, err := u.passwordHash.Hash(req.Password)
	if err != nil {
		return err
	}

	// Create user
	user := &domain.User{
		ID:       userID,
		Username: req.Username,
		Password: hashedPassword,
		Fullname: req.Fullname,
		Role:     req.Role,
	}

	err = u.userRepo.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
