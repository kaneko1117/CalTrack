/**
 * 共通型定義
 */

/** バックエンドからのエラーコード */
export type ErrorCode =
  | "INVALID_REQUEST"
  | "VALIDATION_ERROR"
  | "UNAUTHORIZED"
  | "INTERNAL_ERROR";

/** APIエラーレスポンス */
export type ApiErrorResponse = {
  code: ErrorCode;
  message: string;
  details?: string[];
};
