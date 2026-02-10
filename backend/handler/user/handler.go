package user

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
	"caltrack/handler/common"
	"caltrack/handler/user/dto"
)

// UserUsecaseInterface はUserUsecaseのインターフェース
type UserUsecaseInterface interface {
	Register(ctx context.Context, user *entity.User) (*entity.User, error)
	GetProfile(ctx context.Context, userID vo.UserID) (*entity.User, error)
	UpdateProfile(ctx context.Context, userID vo.UserID, nickname vo.Nickname, height vo.Height, weight vo.Weight, activityLevel vo.ActivityLevel) (*entity.User, error)
}

type UserHandler struct {
	usecase UserUsecaseInterface
}

func NewUserHandler(uc UserUsecaseInterface) *UserHandler {
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
		details := common.ExtractErrorMessages(validationErrs)
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

// GetProfile は認証ユーザーのプロフィールを取得する
func (h *UserHandler) GetProfile(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		common.RespondError(c, http.StatusUnauthorized, common.CodeUnauthorized, "User not authenticated", nil)
		return
	}

	userID := vo.ReconstructUserID(userIDStr.(string))

	user, err := h.usecase.GetProfile(c.Request.Context(), userID)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.NewGetProfileResponse(user))
}

// UpdateProfile は認証ユーザーのプロフィールを更新する
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		common.RespondError(c, http.StatusUnauthorized, common.CodeUnauthorized, "User not authenticated", nil)
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespondError(c, http.StatusBadRequest, common.CodeInvalidRequest, "Invalid request body", nil)
		return
	}

	userID := vo.ReconstructUserID(userIDStr.(string))

	// DTOからVOに変換
	nickname, height, weight, activityLevel, errs := req.ToDomain()
	if errs != nil {
		details := common.ExtractErrorMessages(errs)
		common.RespondValidationError(c, details)
		return
	}

	updatedUser, err := h.usecase.UpdateProfile(c.Request.Context(), userID, nickname, height, weight, activityLevel)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.NewUpdateProfileResponse(updatedUser))
}

func (h *UserHandler) handleError(c *gin.Context, err error) {
	if errors.Is(err, domainErrors.ErrEmailAlreadyExists) {
		common.RespondError(c, http.StatusConflict, common.CodeEmailAlreadyExists, err.Error(), nil)
		return
	}

	if errors.Is(err, domainErrors.ErrUserNotFound) {
		common.RespondError(c, http.StatusNotFound, common.CodeNotFound, "User not found", nil)
		return
	}

	if isValidationError(err) {
		common.RespondValidationError(c, []string{err.Error()})
		return
	}

	// 500エラーの場合はHandler層のログヘルパーでログ出力
	common.LogError("handleError", err, "method", c.Request.Method, "path", c.Request.URL.Path)
	common.RespondError(c, http.StatusInternalServerError, common.CodeInternalError, "Internal server error", nil)
}

func isValidationError(err error) bool {
	validationErrors := []error{
		domainErrors.ErrNicknameRequired,
		domainErrors.ErrNicknameTooLong,
		domainErrors.ErrWeightMustBePositive,
		domainErrors.ErrWeightTooHeavy,
		domainErrors.ErrHeightMustBePositive,
		domainErrors.ErrHeightTooTall,
		domainErrors.ErrInvalidActivityLevel,
	}
	for _, ve := range validationErrors {
		if errors.Is(err, ve) {
			return true
		}
	}
	return false
}
