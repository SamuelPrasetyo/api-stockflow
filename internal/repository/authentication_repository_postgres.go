package repository

import (
	"api-stockflow/internal/domain"
	"context"
	"database/sql"
	"time"
)

type authenticationRepository struct {
	db *sql.DB
}

// NewAuthenticationRepository creates a new authentication repository
func NewAuthenticationRepository(db *sql.DB) AuthenticationRepository {
	return &authenticationRepository{db: db}
}

func (r *authenticationRepository) AddToken(ctx context.Context, token string, userID string) error {
	query := `
		INSERT INTO authentications (token, user_id, created_at)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.ExecContext(ctx, query, token, userID, time.Now())
	return err
}

func (r *authenticationRepository) CheckTokenAvailability(ctx context.Context, token string) error {
	query := `SELECT token FROM authentications WHERE token = $1`

	var storedToken string
	err := r.db.QueryRowContext(ctx, query, token).Scan(&storedToken)

	if err == sql.ErrNoRows {
		return domain.ErrRefreshTokenNotFound
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *authenticationRepository) DeleteToken(ctx context.Context, token string) error {
	query := `DELETE FROM authentications WHERE token = $1`

	result, err := r.db.ExecContext(ctx, query, token)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrRefreshTokenNotFound
	}

	return nil
}
