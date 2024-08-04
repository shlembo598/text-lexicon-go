package auth

import (
	"github.com/google/uuid"
	"golang.org/x/net/context"

	"github.com/shlembo598/text-lexicon-go/internal/models"
)

type Repository interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID uuid.UUID) error
	GetById(ctx context.Context, userID uuid.UUID) (user *models.User, err error)
	FindByEmail(ctx context.Context, user *models.User) (*models.User, error)
}
