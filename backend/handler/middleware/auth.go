package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	domainErrors "caltrack/domain/errors"
	"caltrack/handler/common"
	"caltrack/usecase"
)

// AuthMiddleware は認証ミドルウェアを生成する
func AuthMiddleware(authUsecase *usecase.AuthUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		// CookieからセッションIDを取得
		sessionID, err := c.Cookie("session_id")
		if err != nil {
			common.RespondError(c, http.StatusUnauthorized, common.CodeUnauthorized, "Authentication required", nil)
			c.Abort()
			return
		}

		// セッションの有効性を検証
		session, err := authUsecase.ValidateSession(c.Request.Context(), sessionID)
		if err != nil {
			handleAuthError(c, err)
			c.Abort()
			return
		}

		// コンテキストにユーザー情報を設定
		c.Set("userID", session.UserID().String())
		c.Set("sessionID", session.ID().String())
		c.Next()
	}
}

// handleAuthError は認証エラーをHTTPレスポンスに変換する
func handleAuthError(c *gin.Context, err error) {
	// 無効なセッションID
	if errors.Is(err, domainErrors.ErrInvalidSessionID) {
		common.RespondError(c, http.StatusUnauthorized, common.CodeUnauthorized, "Invalid session", nil)
		return
	}

	// セッションが見つからない
	if errors.Is(err, domainErrors.ErrSessionNotFound) {
		common.RespondError(c, http.StatusUnauthorized, common.CodeUnauthorized, "Invalid session", nil)
		return
	}

	// セッション期限切れ
	if errors.Is(err, domainErrors.ErrSessionExpired) {
		common.RespondError(c, http.StatusUnauthorized, common.CodeSessionExpired, "Session has expired", nil)
		return
	}

	// その他のエラー
	common.RespondError(c, http.StatusInternalServerError, common.CodeInternalError, "Internal server error", err)
}
