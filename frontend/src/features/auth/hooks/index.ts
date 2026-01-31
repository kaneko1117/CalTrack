/**
 * 認証フック
 * ユーザー認証に関するカスタムフック
 */

import { useState, useCallback } from "react";
import { registerUser, ApiError } from "../api";
import type { RegisterUserRequest } from "../types";
import {
  ERROR_CODE_INTERNAL_ERROR,
  ERROR_MESSAGE_UNEXPECTED,
} from "../types";

/**
 * useRegisterUserフックの戻り値の型
 */
export interface UseRegisterUserReturn {
  register: (request: RegisterUserRequest) => Promise<void>;
  isLoading: boolean;
  error: ApiError | null;
  isSuccess: boolean;
  reset: () => void;
}

/**
 * ユーザー登録フック
 * ローディング状態、エラー状態、成功状態を管理
 *
 * @returns UseRegisterUserReturn
 *
 * @example
 * ```tsx
 * const { register, isLoading, error, isSuccess, reset } = useRegisterUser();
 *
 * const handleSubmit = async (data: RegisterUserRequest) => {
 *   await register(data);
 * };
 * ```
 */
export function useRegisterUser(): UseRegisterUserReturn {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<ApiError | null>(null);
  const [isSuccess, setIsSuccess] = useState(false);

  const register = useCallback(async (request: RegisterUserRequest) => {
    setIsLoading(true);
    setError(null);
    setIsSuccess(false);

    try {
      await registerUser(request);
      setIsSuccess(true);
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
    register,
    isLoading,
    error,
    isSuccess,
    reset,
  };
}
