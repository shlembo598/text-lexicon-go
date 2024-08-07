package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/shlembo598/text-lexicon-go/internal/auth"
	"github.com/shlembo598/text-lexicon-go/internal/config"
	"github.com/shlembo598/text-lexicon-go/internal/models"
	"github.com/shlembo598/text-lexicon-go/pkg/utils"
	"github.com/shlembo598/text-lexicon-go/pkg/utils/httpErrors"
	r "github.com/shlembo598/text-lexicon-go/pkg/utils/responses"
)

type authHandlers struct {
	cfg    *config.Config
	authUC auth.UseCase
}

func NewAuthHandlers(cfg *config.Config, authUC auth.UseCase) auth.Handlers {
	return &authHandlers{cfg: cfg, authUC: authUC}
}

// Register godoc
// @Summary Register new user
// @Description register new user, returns user and token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 201 {object} models.User
// @Router /auth/register [post]
func (h *authHandlers) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: tracing

		user := &models.User{}
		if err := utils.ReadRequest(c, user); err != nil {
			utils.LogResponseError(c, err)
			return c.JSON(r.ErrorResponse(err))
		}

		createdUser, err := h.authUC.Register(utils.GetRequestCtx(c), user)
		if err != nil {
			utils.LogResponseError(c, err)
			return c.JSON(r.ErrorResponse(err))
		}

		return c.JSON(http.StatusCreated, r.SuccessResponse(createdUser))
	}
}

// Login godoc
// @Summary Login new user
// @Description login user, returns user and set session
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Router /auth/login [post]
func (h *authHandlers) Login() echo.HandlerFunc {
	type Login struct {
		Email    string `json:"email" db:"email" validate:"omitempty,lte=60,email"`
		Password string `json:"password,omitempty" db:"password" validate:"required,gte=6"`
	}

	return func(c echo.Context) error {

		login := &Login{}
		if err := utils.ReadRequest(c, login); err != nil {
			utils.LogResponseError(c, err)
			return c.JSON(r.ErrorResponse(err))
		}

		userWithToken, err := h.authUC.Login(
			utils.GetRequestCtx(c), &models.User{
				Email:    login.Email,
				Password: login.Password,
			},
		)
		if err != nil {
			utils.LogResponseError(c, err)
			return c.JSON(r.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, r.SuccessResponse(userWithToken))
	}
}

// Update godoc
// @Summary Update user
// @Description update existing user
// @Tags Auth
// @Accept json
// @Param id path int true "user_id"
// @Produce json
// @Success 200 {object} models.User
// @Router /auth/{id} [put]
func (h *authHandlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			utils.LogResponseError(c, err)
			return c.JSON(r.ErrorResponse(err))
		}

		user := &models.User{}
		user.UserID = uID

		if err = utils.ReadRequest(c, user); err != nil {
			utils.LogResponseError(c, err)
			return c.JSON(r.ErrorResponse(err))
		}

		updatedUser, err := h.authUC.Update(utils.GetRequestCtx(c), user)
		if err != nil {
			utils.LogResponseError(c, err)
			return c.JSON(r.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, r.SuccessResponse(updatedUser))
	}
}

// Delete
// @Summary Delete user account
// @Description some description
// @Tags Auth
// @Accept json
// @Param id path int true "user_id"
// @Produce json
// @Success 200 {string} string	"ok"
// @Failure 500 {object} httpErrors.RestError
// @Router /auth/{id} [delete]
func (h *authHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		uId, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			utils.LogResponseError(c, err)
			return c.JSON(r.ErrorResponse(err))
		}

		if err = h.authUC.Delete(utils.GetRequestCtx(c), uId); err != nil {
			utils.LogResponseError(c, err)
			return c.JSON(r.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, r.SuccessResponse("Successfully deleted"))
	}
}

// GetUserByID godoc
// @Summary get user by id
// @Description get string by ID
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param id path int true "user_id"
// @Success 200 {object} models.User
// @Failure 500 {object} httpErrors.RestError
// @Router /auth/{id} [get]
func (h *authHandlers) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		uId, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			utils.LogResponseError(c, err)
			return c.JSON(r.ErrorResponse(err))
		}

		user, err := h.authUC.GetByID(utils.GetRequestCtx(c), uId)
		if err != nil {
			utils.LogResponseError(c, err)
			return c.JSON(r.ErrorResponse(err))
		}

		return c.JSON(http.StatusOK, r.SuccessResponse(user))
	}
}

// GetMe godoc
// @Summary Get user by id
// @Description Get current user by id
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Failure 500 {object} httpErrors.RestError
// @Router /auth/me [get]
func (h *authHandlers) GetMe() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("user").(*models.User)
		if !ok {
			utils.LogResponseError(c, httpErrors.Unauthorized)
			return c.JSON(r.ErrorResponse(httpErrors.Unauthorized))
		}

		return c.JSON(http.StatusOK, r.SuccessResponse(user))
	}
}
