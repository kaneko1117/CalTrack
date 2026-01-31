package usecase

import (
	"context"

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/repository"
	"caltrack/domain/vo"
)

// AuthUsecase は認証に関するユースケースを提供する
type AuthUsecase struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	txManager   repository.TransactionManager
}

// NewAuthUsecase は AuthUsecase のインスタンスを生成する
func NewAuthUsecase(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	txManager repository.TransactionManager,
) *AuthUsecase {
	return &AuthUsecase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		txManager:   txManager,
	}
}

// LoginInput はログイン処理の入力を表す
type LoginInput struct {
	Email    string
	Password string
}

// LoginOutput はログイン処理の出力を表す
type LoginOutput struct {
	Session *entity.Session
	User    *entity.User
}

// Login はメールアドレスとパスワードでユーザーを認証し、セッションを作成する
func (u *AuthUsecase) Login(ctx context.Context, input LoginInput) (*LoginOutput, error) {
	// メールアドレスのバリデーション
	email, err := vo.NewEmail(input.Email)
	if err != nil {
		logWarn("Login", "invalid email format", "email", input.Email)
		return nil, domainErrors.ErrInvalidCredentials
	}

	// パスワードのバリデーション
	password, err := vo.NewPassword(input.Password)
	if err != nil {
		logWarn("Login", "invalid password format")
		return nil, domainErrors.ErrInvalidCredentials
	}

	var output *LoginOutput

	err = u.txManager.Execute(ctx, func(txCtx context.Context) error {
		// ユーザーをメールアドレスで検索
		user, err := u.userRepo.FindByEmail(txCtx, email)
		if err != nil {
			logError("Login", err, "email", email.String())
			return err
		}
		if user == nil {
			logWarn("Login", "user not found", "email", email.String())
			return domainErrors.ErrInvalidCredentials
		}

		// パスワードの照合
		if !user.HashedPassword().Compare(password) {
			logWarn("Login", "password mismatch", "email", email.String())
			return domainErrors.ErrInvalidCredentials
		}

		// セッションの作成
		session, err := entity.NewSessionWithUserID(user.ID())
		if err != nil {
			logError("Login", err, "user_id", user.ID().String())
			return err
		}

		// セッションの保存
		if err := u.sessionRepo.Save(txCtx, session); err != nil {
			logError("Login", err, "session_id", session.ID().String())
			return err
		}

		output = &LoginOutput{
			Session: session,
			User:    user,
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return output, nil
}

// Logout はセッションを削除してログアウトする
func (u *AuthUsecase) Logout(ctx context.Context, sessionIDStr string) error {
	// セッションIDのバリデーション
	sessionID, err := vo.ParseSessionID(sessionIDStr)
	if err != nil {
		logWarn("Logout", "invalid session id")
		return domainErrors.ErrInvalidSessionID
	}

	err = u.txManager.Execute(ctx, func(txCtx context.Context) error {
		// セッションの削除
		if err := u.sessionRepo.DeleteByID(txCtx, sessionID); err != nil {
			logError("Logout", err, "session_id", sessionIDStr)
			return err
		}
		return nil
	})

	return err
}

// ValidateSession はセッションの有効性を検証する
func (u *AuthUsecase) ValidateSession(ctx context.Context, sessionIDStr string) (*entity.Session, error) {
	// セッションIDのバリデーション
	sessionID, err := vo.ParseSessionID(sessionIDStr)
	if err != nil {
		return nil, domainErrors.ErrInvalidSessionID
	}

	// セッションの取得
	session, err := u.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		logError("ValidateSession", err, "session_id", sessionIDStr)
		return nil, err
	}
	if session == nil {
		return nil, domainErrors.ErrSessionNotFound
	}

	// 有効期限の検証
	if err := session.ValidateNotExpired(); err != nil {
		logWarn("ValidateSession", "session expired", "session_id", sessionIDStr)
		return nil, err
	}

	return session, nil
}
