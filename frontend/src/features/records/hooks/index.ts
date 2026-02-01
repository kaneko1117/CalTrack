/**
 * カロリー記録フック
 * カロリー記録に関するカスタムフック
 */

import { useState, useCallback } from "react";
import { createRecord, ApiError } from "../api";
import type { CreateRecordRequest, CreateRecordResponse } from "../types";
import type { ErrorCode } from "@/features/common";

/** エラーコード定数 */
const ERROR_CODE_INTERNAL_ERROR: ErrorCode = "INTERNAL_ERROR";

/** エラーメッセージ定数 */
const ERROR_MESSAGE_UNEXPECTED = "予期しないエラーが発生しました";

/**
 * useCreateRecordフックの戻り値の型
 */
export type UseCreateRecordReturn = {
  createRecord: (request: CreateRecordRequest, onSuccess?: (response: CreateRecordResponse) => void) => Promise<void>;
  isLoading: boolean;
  error: ApiError | null;
  isSuccess: boolean;
  reset: () => void;
};

/**
 * カロリー記録作成フック
 * ローディング状態、エラー状態、成功状態を管理
 *
 * @returns UseCreateRecordReturn
 *
 * @example
 * ```tsx
 * const { createRecord, isLoading, error, isSuccess, reset } = useCreateRecord();
 *
 * const handleSubmit = async (data: CreateRecordRequest) => {
 *   await createRecord(data, (response) => {
 *     console.log("記録作成成功:", response);
 *   });
 * };
 * ```
 */
export function useCreateRecord(): UseCreateRecordReturn {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<ApiError | null>(null);
  const [isSuccess, setIsSuccess] = useState(false);

  const handleCreateRecord = useCallback(async (
    request: CreateRecordRequest,
    onSuccess?: (response: CreateRecordResponse) => void
  ) => {
    setIsLoading(true);
    setError(null);
    setIsSuccess(false);

    try {
      const response = await createRecord(request);
      setIsSuccess(true);
      // 成功時にコールバックを呼び出し
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
    createRecord: handleCreateRecord,
    isLoading,
    error,
    isSuccess,
    reset,
  };
}
