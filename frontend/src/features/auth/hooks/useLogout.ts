/**
 * ログアウトフック
 * ユーザーログアウトに関するカスタムフック
 */

import { useState, useCallback } from "react";
import { logout, ApiError } from "../api";
import {
  ERROR_CODE_INTERNAL_ERROR,
  ERROR_MESSAGE_UNEXPECTED,
} from "../types";

/**
 * useLogoutフックの戻り値の型
 */
export type UseLogoutReturn = {
  logout: (onSuccess?: () => void) => Promise<void>;
  isLoading: boolean;
  error: ApiError | null;
  isSuccess: boolean;
  reset: () => void;
};

/**
 * ユーザーログアウトフック
 * ローディング状態、エラー状態、成功状態を管理
 *
 * @returns UseLogoutReturn
 *
 * @example
 * ```tsx
 * const { logout, isLoading, error, isSuccess, reset } = useLogout();
 *
 * const handleLogout = async () => {
 *   await logout(() => {
 *     // ログアウト成功時の処理
 *     navigate('/login');
 *   });
 * };
 * ```
 */
export function useLogout(): UseLogoutReturn {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<ApiError | null>(null);
  const [isSuccess, setIsSuccess] = useState(false);

  const logoutHandler = useCallback(async (
    onSuccess?: () => void
  ) => {
    setIsLoading(true);
    setError(null);
    setIsSuccess(false);

    try {
      await logout();
      setIsSuccess(true);
      // 成功時にコールバックを直接呼び出し
      onSuccess?.();
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err);
      } else {
        setError(
          new ApiError(ERROR_CODE_INTERNAL_ERROR, ERROR_MESSAGE_UNEXPECTED, 500)
        );
      }
    } finally {
      setIsLoading(false);
    }
  }, []);

  const reset = useCallback(() => {
    setIsLoading(false);
    setError(null);
    setIsSuccess(false);
  }, []);

  return {
    logout: logoutHandler,
    isLoading,
    error,
    isSuccess,
    reset,
  };
}
