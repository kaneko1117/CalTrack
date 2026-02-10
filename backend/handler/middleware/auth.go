package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
	"caltrack/handler/common"
)

// AuthSessionValidator はセッション検証用のインターフェース
type AuthSessionValidator interface {
	ValidateSession(ctx context.Context, sessionID vo.SessionID) (*entity.Session, error)
}

// AuthMiddleware は認証ミドルウェアを生成する
func AuthMiddleware(authUsecase AuthSessionValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		// CookieからセッションIDを取得
		sessionIDStr, err := c.Cookie("session_id")
		if err != nil {
			common.RespondError(c, http.StatusUnauthorized, common.CodeUnauthorized, "Authentication required", nil)
			c.Abort()
			return
		}

		// セッションIDをVOに変換
		sessionID, err := vo.ParseSessionID(sessionIDStr)
		if err != nil {
			common.RespondError(c, http.StatusUnauthorized, common.CodeUnauthorized, "Invalid session", nil)
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
