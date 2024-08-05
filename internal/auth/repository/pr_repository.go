package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"

	"github.com/shlembo598/text-lexicon-go/internal/auth"
	"github.com/shlembo598/text-lexicon-go/internal/models"
)

type authRepo struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) auth.Repository {
	return &authRepo{db: db}
}

// Create new user
func (r *authRepo) Register(ctx context.Context, user *models.User) (*models.User, error) {
	const op = "auth.pg_repository.register"

	query, args, err := createUserQuery(user)
	if err != nil {
		return nil, fmt.Errorf("%s.query: %w", op, err)
	}

	u := &models.User{}
	if err = r.db.QueryRowxContext(
		ctx, query, args...,
	).StructScan(u); err != nil {
		return nil, fmt.Errorf("%s.StructScan: %w", op, err)
	}

	return u, nil
}

// Update existing user
func (r *authRepo) Update(ctx context.Context, user *models.User) (*models.User, error) {
	const op = "auth.pg_repository.update"

	query, args, buildErr := updateUserQuery(user)
	if buildErr != nil {
		return nil, fmt.Errorf("%s.query: %w", op, buildErr)
	}

	u := &models.User{}
	if err := r.db.GetContext(
		ctx, u, query, args...,
	); err != nil {
		return nil, fmt.Errorf("%s.GetContext: %w", op, err)
	}

	return u, nil
}

// Delete existing user
func (r *authRepo) Delete(ctx context.Context, userID uuid.UUID) error {
	const op = "auth.pg_repository.delete"

	query, args, buildErr := deleteUserQuery(userID)
	if buildErr != nil {
		return fmt.Errorf("%s.query: %w", op, buildErr)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s.ExecContext: %w", op, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s.RowsAffected: %w", op, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%s.rowsAffected: %w", op, sql.ErrNoRows)
	}

	return nil
}

// Get user by id
func (r *authRepo) GetById(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	const op = "auth.pg_repository.getById"

	query, args, buildErr := getUserQuery(userID)
	if buildErr != nil {
		return nil, fmt.Errorf("%s.query: %w", op, buildErr)
	}

	user := &models.User{}
	if err := r.db.QueryRowxContext(ctx, query, args...).StructScan(user); err != nil {
		return nil, fmt.Errorf("%s.QueryRowxContext: %w", op, buildErr)
	}

	return user, nil
}

// Find user by email
func (r *authRepo) FindByEmail(ctx context.Context, user *models.User) (*models.User, error) {
	const op = "auth.pg_repository.findByEmail"

	foundUser := &models.User{}

	query, args, buildErr := findUserByEmail(user.Email)
	if buildErr != nil {
		return nil, fmt.Errorf("%s.query: %w", op, buildErr)
	}

	if err := r.db.QueryRowxContext(ctx, query, args...).StructScan(foundUser); err != nil {
		return nil, fmt.Errorf("%s.QueryRowxContext: %w", op, buildErr)
	}

	return foundUser, nil
}
