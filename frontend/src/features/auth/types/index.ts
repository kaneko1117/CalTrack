/**
 * 認証関連の型定義
 * ユーザー登録に関する型を定義
 */

/** 性別 */
export type Gender = "male" | "female" | "other";

/** 活動レベル */
export type ActivityLevel =
  | "sedentary"
  | "light"
  | "moderate"
  | "active"
  | "veryActive";

/** バックエンドからのエラーコード */
export type ErrorCode =
  | "INVALID_REQUEST"
  | "VALIDATION_ERROR"
  | "EMAIL_ALREADY_EXISTS"
  | "INVALID_CREDENTIALS"
  | "INTERNAL_ERROR";

/** エラーコード定数 */
export const ERROR_CODE_INTERNAL_ERROR: ErrorCode = "INTERNAL_ERROR";
export const ERROR_CODE_INVALID_CREDENTIALS: ErrorCode = "INVALID_CREDENTIALS";

/** エラーメッセージ定数 */
export const ERROR_MESSAGE_UNEXPECTED = "予期しないエラーが発生しました";
export const ERROR_MESSAGE_INVALID_CREDENTIALS = "メールアドレスまたはパスワードが間違っています";

/** ユーザー登録リクエスト */
export type RegisterUserRequest = {
  email: string;
  password: string;
  nickname: string;
  weight: number;
  height: number;
  birthDate: string; // YYYY-MM-DD形式
  gender: Gender;
  activityLevel: ActivityLevel;
};

/** ユーザー登録レスポンス */
export type RegisterUserResponse = {
  userId: string;
};

/** APIエラーレスポンス */
export type ApiErrorResponse = {
  code: ErrorCode;
  message: string;
  details?: string[];
};

/** ログインリクエスト */
export type LoginRequest = {
  email: string;
  password: string;
};

/** ログインレスポンス */
export type LoginResponse = {
  userId: string;
  email: string;
  nickname: string;
};

/** ログアウトレスポンス */
export type LogoutResponse = {
  message: string;
};
