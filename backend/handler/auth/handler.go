package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
	"caltrack/handler/auth/dto"
	"caltrack/handler/common"
	"caltrack/usecase"
)

// Cookie設定定数
const (
	sessionCookieName = "session_id"
	cookiePath        = "/"
	// 有効期限は7日間（秒単位）
	cookieMaxAge = int(vo.SessionDurationDays * 24 * time.Hour / time.Second)
)

// AuthHandler は認証関連のHTTPハンドラ
type AuthHandler struct {
	usecase *usecase.AuthUsecase
}

// NewAuthHandler は AuthHandler のインスタンスを生成する
func NewAuthHandler(uc *usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{usecase: uc}
}

// Login はログイン処理を行う
// POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespondError(c, http.StatusBadRequest, common.CodeInvalidRequest, "Invalid request body", nil)
		return
	}

	// Usecase実行
	output, err := h.usecase.Login(c.Request.Context(), req.ToInput())
	if err != nil {
		h.handleError(c, err)
		return
	}

	// セッションCookieを設定
	h.setSessionCookie(c, output.Session.ID().String())

	// レスポンス返却
	c.JSON(http.StatusOK, dto.NewLoginResponse(output))
}

// Logout はログアウト処理を行う
// POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	// CookieからセッションIDを取得
	sessionID, err := c.Cookie(sessionCookieName)
	if err != nil {
		// Cookieが無い場合は既にログアウト済みとして成功扱い
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
		return
	}

	// Usecase実行
	if err := h.usecase.Logout(c.Request.Context(), sessionID); err != nil {
		h.handleError(c, err)
		return
	}

	// セッションCookieを削除
	h.clearSessionCookie(c)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// setSessionCookie はセッションCookieを設定する
func (h *AuthHandler) setSessionCookie(c *gin.Context, sessionID string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		sessionCookieName,
		sessionID,
		cookieMaxAge,
		cookiePath,
		"",    // domain: 空文字でリクエストドメインを使用
		true,  // secure: HTTPS必須
		true,  // httpOnly: JavaScriptからアクセス不可
	)
}

// clearSessionCookie はセッションCookieを削除する
func (h *AuthHandler) clearSessionCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		sessionCookieName,
		"",
		-1, // maxAge: 負の値で削除
		cookiePath,
		"",
		true,
		true,
	)
}

// handleError はドメインエラーをHTTPレスポンスに変換する
func (h *AuthHandler) handleError(c *gin.Context, err error) {
	// 認証エラー（メールアドレスまたはパスワードが間違っている）
	if errors.Is(err, domainErrors.ErrInvalidCredentials) {
		common.RespondError(c, http.StatusUnauthorized, common.CodeInvalidCredentials, "Invalid email or password", nil)
		return
	}

	// 無効なセッションID
	if errors.Is(err, domainErrors.ErrInvalidSessionID) {
		common.RespondError(c, http.StatusBadRequest, common.CodeInvalidRequest, "Invalid session", nil)
		return
	}

	// セッション期限切れ
	if errors.Is(err, domainErrors.ErrSessionExpired) {
		common.RespondError(c, http.StatusUnauthorized, common.CodeSessionExpired, "Session has expired", nil)
		return
	}

	// セッションが見つからない
	if errors.Is(err, domainErrors.ErrSessionNotFound) {
		common.RespondError(c, http.StatusUnauthorized, common.CodeUnauthorized, "Session not found", nil)
		return
	}

	// 500エラー
	common.LogError("handleError", err, "method", c.Request.Method, "path", c.Request.URL.Path)
	common.RespondError(c, http.StatusInternalServerError, common.CodeInternalError, "Internal server error", nil)
}
