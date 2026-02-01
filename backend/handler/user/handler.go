package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

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

// Register はユーザー登録を行う
// @Summary ユーザー登録
// @Description 新規ユーザーを登録する
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.RegisterUserRequest true "ユーザー登録リクエスト"
// @Success 201 {object} dto.RegisterUserResponse "登録成功"
// @Failure 400 {object} common.ErrorResponse "バリデーションエラー"
// @Failure 409 {object} common.ErrorResponse "メールアドレス重複"
// @Failure 500 {object} common.ErrorResponse "サーバーエラー"
// @Router /users [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespondError(c, http.StatusBadRequest, common.CodeInvalidRequest, "Invalid request body", nil)
		return
	}

	user, parseErr, validationErrs := req.ToDomain()
	if parseErr != nil {
		common.RespondError(c, http.StatusBadRequest, common.CodeValidationError, "Invalid birth date format. Use YYYY-MM-DD", nil)
		return
	}
	if validationErrs != nil {
		details := extractErrorMessages(validationErrs)
		common.RespondValidationError(c, details)
		return
	}

	registeredUser, err := h.usecase.Register(c.Request.Context(), user)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.RegisterUserResponse{
		UserID: registeredUser.ID().String(),
	})
}

func (h *UserHandler) handleError(c *gin.Context, err error) {
	if errors.Is(err, domainErrors.ErrEmailAlreadyExists) {
		common.RespondError(c, http.StatusConflict, common.CodeEmailAlreadyExists, err.Error(), nil)
		return
	}

	// 500エラーの場合はHandler層のログヘルパーでログ出力
	common.LogError("handleError", err, "method", c.Request.Method, "path", c.Request.URL.Path)
	common.RespondError(c, http.StatusInternalServerError, common.CodeInternalError, "Internal server error", nil)
}

func extractErrorMessages(errs []error) []string {
	details := make([]string, len(errs))
	for i, err := range errs {
		details[i] = err.Error()
	}
	return details
}
