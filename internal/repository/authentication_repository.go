package repository

import (
	"context"
)

// AuthenticationRepository defines the interface for authentication data operations
type AuthenticationRepository interface {
	AddToken(ctx context.Context, token string, userID string) error
	CheckTokenAvailability(ctx context.Context, token string) error
	DeleteToken(ctx context.Context, token string) error
}
