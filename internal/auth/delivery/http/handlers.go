package http

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/shlembo598/text-lexicon-go/internal/auth"
	"github.com/shlembo598/text-lexicon-go/internal/config"
	"github.com/shlembo598/text-lexicon-go/internal/models"
	"github.com/shlembo598/text-lexicon-go/pkg/utils"
	"github.com/shlembo598/text-lexicon-go/pkg/utils/responses"
)

type authHandlers struct {
	cfg    *config.Config
	authUC auth.UseCase
}

func NewAuthHandlers(cfg *config.Config, authUC auth.UseCase) auth.Handlers {
	return &authHandlers{cfg: cfg, authUC: authUC}
}

// TODO: swagger specification
func (h *authHandlers) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: tracing

		user := &models.User{}
		if err := utils.ReadRequest(c, user); err != nil {
			utils.LogResponseError(c, err)
			return responses.ErrorResponse(c, http.StatusBadRequest, err)
		}

		createdUser, err := h.authUC.Register(utils.GetRequestCtx(c), user)
		if err != nil {
			utils.LogResponseError(c, err)
			return responses.ErrorResponse(c, http.StatusInternalServerError, err)
		}

		return responses.SuccessResponse(c, http.StatusCreated, createdUser)
	}
}

// TODO: swagger specification
func (h *authHandlers) Login() echo.HandlerFunc {
	type Login struct {
		Email    string `json:"email" db:"email" validate:"omitempty,lte=60,email"`
		Password string `json:"password,omitempty" db:"password" validate:"required,gte=6"`
	}

	return func(c echo.Context) error {

		login := &Login{}
		if err := utils.ReadRequest(c, login); err != nil {
			utils.LogResponseError(c, err)
			return responses.ErrorResponse(c, http.StatusBadRequest, err)
		}

		userWithToken, err := h.authUC.Login(
			utils.GetRequestCtx(c), &models.User{
				Email:    login.Email,
				Password: login.Password,
			},
		)
		if err != nil {
			utils.LogResponseError(c, err)
			return responses.ErrorResponse(c, http.StatusUnauthorized, err)
		}

		return responses.SuccessResponse(c, http.StatusOK, userWithToken)
	}
}

// TODO: swagger specification
func (h *authHandlers) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		uID, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			utils.LogResponseError(c, err)
			return responses.ErrorResponse(c, http.StatusBadRequest, err)
		}

		user := &models.User{}
		user.UserID = uID

		if err = utils.ReadRequest(c, user); err != nil {
			utils.LogResponseError(c, err)
			return responses.ErrorResponse(c, http.StatusBadRequest, err)
		}

		updatedUser, err := h.authUC.Update(utils.GetRequestCtx(c), user)
		if err != nil {
			utils.LogResponseError(c, err)
			return responses.ErrorResponse(c, http.StatusInternalServerError, err)
		}

		return responses.SuccessResponse(c, http.StatusOK, updatedUser)
	}
}

// TODO: swagger specification
func (h *authHandlers) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		uId, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			utils.LogResponseError(c, err)
			return responses.ErrorResponse(c, http.StatusBadRequest, err)
		}

		if err = h.authUC.Delete(utils.GetRequestCtx(c), uId); err != nil {
			utils.LogResponseError(c, err)
			return responses.ErrorResponse(c, http.StatusInternalServerError, err)
		}

		return c.NoContent(http.StatusOK)
	}
}

// TODO: swagger specification
func (h *authHandlers) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		uId, err := uuid.Parse(c.Param("user_id"))
		if err != nil {
			utils.LogResponseError(c, err)
			return responses.ErrorResponse(c, http.StatusBadRequest, err)
		}

		user, err := h.authUC.GetByID(utils.GetRequestCtx(c), uId)
		if err != nil {
			utils.LogResponseError(c, err)
			return responses.ErrorResponse(c, http.StatusNotFound, err)
		}

		return responses.SuccessResponse(c, http.StatusOK, user)
	}
}

// TODO: swagger specification
func (h *authHandlers) GetMe() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, ok := c.Get("User").(*models.User)
		if !ok {
			err := errors.New("Unauthorized")
			utils.LogResponseError(c, err)
			return responses.ErrorResponse(c, http.StatusUnauthorized, err)
		}

		return responses.SuccessResponse(c, http.StatusOK, user)
	}
}
