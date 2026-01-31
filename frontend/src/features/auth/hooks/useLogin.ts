/**
 * ログインフック
 * ユーザーログインに関するカスタムフック
 */

import { useState, useCallback } from "react";
import { login, ApiError } from "../api";
import type { LoginRequest, LoginResponse } from "../types";
import {
  ERROR_CODE_INTERNAL_ERROR,
  ERROR_MESSAGE_UNEXPECTED,
} from "../types";

/**
 * useLoginフックの戻り値の型
 */
export type UseLoginReturn = {
  login: (request: LoginRequest, onSuccess?: (response: LoginResponse) => void) => Promise<void>;
  isLoading: boolean;
  error: ApiError | null;
  isSuccess: boolean;
  reset: () => void;
};

/**
 * ユーザーログインフック
 * ローディング状態、エラー状態、成功状態を管理
 *
 * @returns UseLoginReturn
 *
 * @example
 * ```tsx
 * const { login, isLoading, error, isSuccess, reset } = useLogin();
 *
 * const handleSubmit = async (data: LoginRequest) => {
 *   await login(data, (response) => {
 *     // ログイン成功時の処理
 *     console.log('Logged in:', response.nickname);
 *   });
 * };
 * ```
 */
export function useLogin(): UseLoginReturn {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<ApiError | null>(null);
  const [isSuccess, setIsSuccess] = useState(false);

  const loginHandler = useCallback(async (
    request: LoginRequest,
    onSuccess?: (response: LoginResponse) => void
  ) => {
    setIsLoading(true);
    setError(null);
    setIsSuccess(false);

    try {
      const response = await login(request);
      setIsSuccess(true);
      // 成功時にコールバックを直接呼び出し
      onSuccess?.(response);
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
    login: loginHandler,
    isLoading,
    error,
    isSuccess,
    reset,
  };
}
