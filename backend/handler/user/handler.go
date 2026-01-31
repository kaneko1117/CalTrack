package user

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	domainErrors "caltrack/domain/errors"
	"caltrack/handler/common"
	"caltrack/handler/user/dto"
	"caltrack/usecase"
)

type UserHandler struct {
	usecase *usecase.UserUsecase
}

func NewUserHandler(uc *usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: uc}
}

func (h *UserHandler) Register(c echo.Context) error {
	var req dto.RegisterUserRequest
	if err := c.Bind(&req); err != nil {
		return common.RespondError(c, http.StatusBadRequest, common.CodeInvalidRequest, "Invalid request body")
	}

	user, parseErr, validationErrs := req.ToDomain()
	if parseErr != nil {
		return common.RespondError(c, http.StatusBadRequest, common.CodeValidationError, "Invalid birth date format. Use YYYY-MM-DD")
	}
	if validationErrs != nil {
		details := extractErrorMessages(validationErrs)
		return common.RespondValidationError(c, details)
	}

	registeredUser, err := h.usecase.Register(c.Request().Context(), user)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusCreated, dto.RegisterUserResponse{
		UserID: registeredUser.ID().String(),
	})
}

func (h *UserHandler) handleError(c echo.Context, err error) error {
	if errors.Is(err, domainErrors.ErrEmailAlreadyExists) {
		return common.RespondError(c, http.StatusConflict, common.CodeEmailAlreadyExists, err.Error())
	}

	return common.RespondError(c, http.StatusInternalServerError, common.CodeInternalError, "Internal server error")
}

func extractErrorMessages(errs []error) []string {
	details := make([]string, len(errs))
	for i, err := range errs {
		details[i] = err.Error()
	}
	return details
}
