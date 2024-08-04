package auth

import (
	"context"

	"github.com/google/uuid"

	"github.com/shlembo598/text-lexicon-go/internal/models"
)

type UseCase interface {
	Register(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	Login(ctx context.Context, user *models.User) (*models.UserWithToken, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
}
