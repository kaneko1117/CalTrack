/**
 * Common Helpers
 * 共通ヘルパー関数
 */
import type { Dispatch, SetStateAction } from "react";
import type { Result } from "@/domain/shared/result";
import type { DomainError } from "@/domain/shared/errors";
import type { ErrorCode } from "@/lib/api";

/** APIエラーメッセージのマッピング */
const API_ERROR_MESSAGES: Record<ErrorCode, string> = {
  INVALID_REQUEST: "リクエストが不正です",
  VALIDATION_ERROR: "入力内容に誤りがあります",
  UNAUTHORIZED: "認証が必要です",
  EMAIL_ALREADY_EXISTS: "このメールアドレスは既に登録されています",
  INVALID_CREDENTIALS: "メールアドレスまたはパスワードが間違っています",
  INTERNAL_ERROR: "予期しないエラーが発生しました",
};

/**
 * getApiErrorMessage - APIエラーコードからユーザー向けメッセージを取得
 * @param code - エラーコード
 * @returns ユーザー向けエラーメッセージ
 */
export function getApiErrorMessage(code: string): string {
  return (
    API_ERROR_MESSAGES[code as ErrorCode] ?? "予期しないエラーが発生しました"
  );
}

/**
 * createFieldHandler - VOファクトリを使った汎用フィールドハンドラ生成
 * @param field - フィールド名
 * @param factory - VOファクトリ関数（newEmail, newPassword等）
 * @param setFormState - フォーム状態のsetter
 * @param setErrors - エラー状態のsetter
 * @returns フィールド値変更ハンドラ
 */
export function createFieldHandler<T, E extends string, F extends string>(
  field: F,
  factory: (value: string) => Result<T, DomainError<E>>,
  setFormState: Dispatch<SetStateAction<Record<F, string>>>,
  setErrors: Dispatch<SetStateAction<Record<F, string | null>>>
): (value: string) => void {
  return (value: string) => {
    setFormState((prev) => ({ ...prev, [field]: value }));
    const result = factory(value);
    if (result.ok) {
      setErrors((prev) => ({ ...prev, [field]: null }));
    } else {
      setErrors((prev) => ({ ...prev, [field]: result.error.message }));
    }
  };
}

/**
 * createResetHandler - フォームリセットハンドラ生成
 * @param resetFormState - フォームリセット関数
 * @param resetErrors - エラーリセット関数
 * @returns リセットハンドラ
 */
export function createResetHandler(
  resetFormState: () => void,
  resetErrors: () => void
): () => void {
  return () => {
    resetFormState();
    resetErrors();
  };
}
