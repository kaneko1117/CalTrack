/**
 * useForm - 汎用フォームフック
 */
import { useState, useCallback, useTransition } from "react";
import type { Result } from "@/domain/shared/result";
import type { DomainError } from "@/domain/shared/errors";
import type { ApiErrorResponse } from "@/lib/api";
import { createFieldHandler, createResetHandler } from "../helpers";

/**
 * VOファクトリの型
 */
type VOFactory<E extends string> = (
  value: string,
) => Result<unknown, DomainError<E>>;

/**
 * useFormの戻り値
 */
export type UseFormReturn<F extends string> = {
  formState: Record<F, string>;
  errors: Record<F, string | null>;
  apiError: ApiErrorResponse | null;
  handleChange: (field: F) => (value: string) => void;
  handleSubmit: (e: React.FormEvent) => void;
  reset: () => void;
  isValid: boolean;
  isPending: boolean;
};

/**
 * useForm - 汎用フォームフック
 * @param config - フィールドごとのVOファクトリ設定
 * @param initialFormState - フォームの初期状態
 * @param initialErrors - エラーの初期状態
 * @param onSubmit - 送信処理（フォームデータを受け取り、結果を返す）
 * @param onSuccess - 成功時コールバック
 */
export function useForm<F extends string, T>(
  config: Record<F, VOFactory<string>>,
  initialFormState: Record<F, string>,
  initialErrors: Record<F, string | null>,
  onSubmit: (formState: Record<F, string>) => Promise<T>,
  onSuccess?: (result: T) => void,
): UseFormReturn<F> {
  const [formState, setFormState] = useState(initialFormState);
  const [errors, setErrors] = useState(initialErrors);
  const [apiError, setApiError] = useState<ApiErrorResponse | null>(null);
  const [isPending, startTransition] = useTransition();

  const handleChange = useCallback(
    (field: F) => {
      const handler = createFieldHandler(
        field,
        config[field],
        setFormState,
        setErrors,
      );
      return (value: string) => {
        handler(value);
        if (apiError) setApiError(null);
      };
    },
    [config, apiError],
  );

  const reset = useCallback(
    () =>
      createResetHandler(
        () => setFormState(initialFormState),
        () => {
          setErrors(initialErrors);
          setApiError(null);
        },
      )(),
    [initialFormState, initialErrors],
  );

  const isValid =
    Object.values(errors).every((e) => e === null) &&
    Object.values(formState).every((v) => v !== "");

  const handleSubmit = useCallback(
    (e: React.FormEvent) => {
      e.preventDefault();
      if (!isValid) return;

      startTransition(async () => {
        try {
          const result = await onSubmit(formState);
          onSuccess?.(result);
        } catch (err) {
          setApiError(err as ApiErrorResponse);
        }
      });
    },
    [isValid, formState, onSubmit, onSuccess],
  );

  return {
    formState,
    errors,
    apiError,
    handleChange,
    handleSubmit,
    reset,
    isValid,
    isPending,
  };
}
