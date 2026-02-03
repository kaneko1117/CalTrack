/**
 * useApi - 汎用APIフック
 * useTransitionを使用したシンプルなAPIフック
 */
import { useState, useCallback, useTransition } from "react";
import type { ApiErrorResponse } from "@/lib/api";

export type UseApiOptions<T> = {
  onSuccess?: (data: T) => void;
  onError?: (error: ApiErrorResponse) => void;
};

export type UseApiReturn<T> = {
  execute: () => void;
  data: T | null;
  error: ApiErrorResponse | null;
  isPending: boolean;
  isSuccess: boolean;
  reset: () => void;
};

export function useApi<T>(
  apiFunction: () => Promise<T>,
  options?: UseApiOptions<T>
): UseApiReturn<T> {
  const [data, setData] = useState<T | null>(null);
  const [error, setError] = useState<ApiErrorResponse | null>(null);
  const [isSuccess, setIsSuccess] = useState(false);
  const [isPending, startTransition] = useTransition();

  const execute = useCallback(() => {
    startTransition(async () => {
      try {
        setError(null);
        setIsSuccess(false);
        const result = await apiFunction();
        setData(result);
        setIsSuccess(true);
        options?.onSuccess?.(result);
      } catch (err) {
        const apiError = err as ApiErrorResponse;
        setError(apiError);
        setIsSuccess(false);
        options?.onError?.(apiError);
      }
    });
  }, [apiFunction, options]);

  const reset = useCallback(() => {
    setData(null);
    setError(null);
    setIsSuccess(false);
  }, []);

  return { execute, data, error, isPending, isSuccess, reset };
}
