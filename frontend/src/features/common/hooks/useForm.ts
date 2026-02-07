/**
 * useForm - 汎用フォームフック
 * SWRベースのuseRequestMutationを使用
 */
import { useState, useCallback } from "react";
import type { Result } from "@/domain/shared/result";
import type { DomainError } from "@/domain/shared/errors";
import type { ApiErrorResponse } from "@/lib/api";
import { createFieldHandler, createResetHandler } from "../helpers";
import { useRequestMutation } from "./useRequest";

type MutationMethod = "POST" | "PUT" | "PATCH" | "DELETE";

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
  setFormState: React.Dispatch<React.SetStateAction<Record<F, string>>>;
};

/**
 * useFormのオプション
 */
export type UseFormOptions<F extends string, T, D> = {
  /** VOファクトリ設定 */
  config: Record<F, VOFactory<string>>;
  /** フォームの初期状態 */
  initialFormState: Record<F, string>;
  /** エラーの初期状態 */
  initialErrors: Record<F, string | null>;
  /** APIエンドポイント */
  url: string;
  /** フォームデータからAPIリクエストデータへの変換 */
  transformData: (formState: Record<F, string>) => D;
  /** 成功時コールバック */
  onSuccess?: (result: T) => void;
  /** HTTPメソッド（デフォルト: "POST"） */
  method?: MutationMethod;
};

/**
 * useForm - 汎用フォームフック（SWRベース）
 * @param options - フォーム設定オプション
 */
export function useForm<F extends string, T, D = unknown>(
  options: UseFormOptions<F, T, D>,
): UseFormReturn<F> {
  const {
    config,
    initialFormState,
    initialErrors,
    url,
    transformData,
    onSuccess,
    method = "POST",
  } = options;

  const [formState, setFormState] = useState(initialFormState);
  const [errors, setErrors] = useState(initialErrors);

  // SWRベースのmutationフック
  const { trigger, isMutating, error: apiError, reset: resetMutation } = useRequestMutation<T, D>(
    url,
    method,
  );

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
        // APIエラーがあればリセット
        if (apiError) resetMutation();
      };
    },
    [config, apiError, resetMutation],
  );

  const reset = useCallback(
    () =>
      createResetHandler(
        () => setFormState(initialFormState),
        () => {
          setErrors(initialErrors);
          resetMutation();
        },
      )(),
    [initialFormState, initialErrors, resetMutation],
  );

  const isValid =
    Object.values(errors).every((e) => e === null) &&
    Object.values(formState).every((v) => v !== "");

  const handleSubmit = useCallback(
    async (e: React.FormEvent) => {
      e.preventDefault();
      if (!isValid) return;

      try {
        const requestData = transformData(formState);
        const result = await trigger(requestData);
        // trigger成功後にonSuccessを呼ぶ
        onSuccess?.(result);
      } catch {
        // エラーはapiErrorで管理される
      }
    },
    [isValid, formState, transformData, trigger, onSuccess],
  );

  return {
    formState,
    errors,
    apiError: apiError ?? null,
    handleChange,
    handleSubmit,
    reset,
    isValid,
    isPending: isMutating,
    setFormState,
  };
}
