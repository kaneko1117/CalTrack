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
  | "INTERNAL_ERROR";

/** エラーコード定数 */
export const ERROR_CODE_INTERNAL_ERROR: ErrorCode = "INTERNAL_ERROR";

/** エラーメッセージ定数 */
export const ERROR_MESSAGE_UNEXPECTED = "予期しないエラーが発生しました";

/** ユーザー登録リクエスト */
export interface RegisterUserRequest {
  email: string;
  password: string;
  nickname: string;
  weight: number;
  height: number;
  birthDate: string; // YYYY-MM-DD形式
  gender: Gender;
  activityLevel: ActivityLevel;
}

/** ユーザー登録レスポンス */
export interface RegisterUserResponse {
  userId: string;
}

/** APIエラーレスポンス */
export interface ApiErrorResponse {
  code: ErrorCode;
  message: string;
  details?: string[];
}
