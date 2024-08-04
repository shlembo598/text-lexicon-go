package usecase

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/net/context"

	"github.com/shlembo598/text-lexicon-go/internal/auth"
	"github.com/shlembo598/text-lexicon-go/internal/config"
	"github.com/shlembo598/text-lexicon-go/internal/models"
	"github.com/shlembo598/text-lexicon-go/pkg/utils"
	"github.com/shlembo598/text-lexicon-go/pkg/utils/httpErrors"
)

type authUC struct {
	cfg      *config.Config
	authRepo auth.Repository
}

func NewAuthUserCase(cfg *config.Config, authRepo auth.Repository) auth.UseCase {
	return &authUC{cfg: cfg, authRepo: authRepo}
}

func (u *authUC) Register(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	const op = "auth.userCase.register"

	existsUser, err := u.authRepo.FindByEmail(ctx, user)
	if existsUser != nil || err == nil {
		return nil, errors.New(httpErrors.ErrEmailAlreadyExists)
	}

	if err = user.PrepareCreate(); err != nil {
		return nil, fmt.Errorf("%s.PrepareCreate: %w", op, err)
	}

	createdUser, err := u.authRepo.Register(ctx, user)
	if err != nil {
		return nil, err
	}

	createdUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(createdUser, u.cfg)
	if err != nil {
		return nil, fmt.Errorf("%s.GenerateJWTToken: %w", op, err)
	}

	return &models.UserWithToken{
		User:  createdUser,
		Token: token,
	}, nil
}

// Login user, returns user model with jwt token
func (u *authUC) Login(ctx context.Context, user *models.User) (*models.UserWithToken, error) {
	const op = "auth.userCase.register"

	foundUser, err := u.authRepo.FindByEmail(ctx, user)
	if err != nil {
		return nil, err
	}

	if err = foundUser.ComparePasswords(user.Password); err != nil {
		return nil, fmt.Errorf("%s.ComparePasswords: %w", op, errors.New(httpErrors.ErrUnauthorized))
	}

	foundUser.SanitizePassword()

	token, err := utils.GenerateJWTToken(foundUser, u.cfg)
	if err != nil {
		return nil, fmt.Errorf("%s.GenerateJWTToken: %w", op, err)
	}

	return &models.UserWithToken{
		User:  foundUser,
		Token: token,
	}, nil
}

func (u *authUC) Update(ctx context.Context, user *models.User) (*models.User, error) {
	const op = "auth.userCase.update"

	if err := user.PrepareUpdate(); err != nil {
		return nil, fmt.Errorf("%s.PrepareUpdate: %w", op, err)
	}

	updatedUser, err := u.authRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	updatedUser.SanitizePassword()

	return updatedUser, nil
}

func (u *authUC) Delete(ctx context.Context, userID uuid.UUID) error {
	if err := u.authRepo.Delete(ctx, userID); err != nil {
		return err
	}

	return nil
}

func (u *authUC) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := u.authRepo.GetById(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.SanitizePassword()

	return user, nil
}
