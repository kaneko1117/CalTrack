/**
 * useRequest - 汎用APIフック（SWRベース）
 * GET用とMutation用を別々のフックとして提供
 */
import useSWR from "swr";
import useSWRMutation from "swr/mutation";
import { fetcher, mutate } from "@/lib/swr";
import type { ApiErrorResponse } from "@/lib/api";

type MutationMethod = "POST" | "PUT" | "DELETE";

type UseRequestOptions<T> = {
  onSuccess?: (data: T) => void;
  onError?: (error: ApiErrorResponse) => void;
};

type UseRequestGetReturn<T> = {
  data: T | undefined;
  error: ApiErrorResponse | undefined;
  isLoading: boolean;
  isValidating: boolean;
  mutate: () => Promise<T | undefined>;
};

type UseRequestMutationReturn<T, D> = {
  trigger: (data?: D) => Promise<T>;
  data: T | undefined;
  error: ApiErrorResponse | undefined;
  isMutating: boolean;
  reset: () => void;
};

/**
 * GETリクエスト用フック
 */
export function useRequestGet<T>(url: string | null): UseRequestGetReturn<T> {
  const { data, error, isLoading, isValidating, mutate: swrMutate } = useSWR<T, ApiErrorResponse>(
    url,
    fetcher,
    { revalidateOnFocus: false }
  );
  return {
    data,
    error,
    isLoading,
    isValidating,
    mutate: async () => swrMutate(),
  };
}

/**
 * POST/PUT/DELETEリクエスト用フック
 */
export function useRequestMutation<T, D = void>(
  url: string,
  method: MutationMethod,
  options?: UseRequestOptions<T>
): UseRequestMutationReturn<T, D> {
  const { trigger, data, error, isMutating, reset } = useSWRMutation<
    T,
    ApiErrorResponse,
    string,
    { method: MutationMethod; data?: D }
  >(url, mutate, {
    onSuccess: options?.onSuccess,
    onError: options?.onError,
  });

  return {
    trigger: async (requestData?: D) => trigger({ method, data: requestData }),
    data,
    error,
    isMutating,
    reset,
  };
}

/**
 * 後方互換性のためのオーバーロード関数（非推奨）
 * 新規コードでは useRequestGet / useRequestMutation を使用すること
 */
export function useRequest<T>(url: string | null): UseRequestGetReturn<T>;
export function useRequest<T, D = void>(
  url: string,
  options: UseRequestOptions<T> & { method: MutationMethod }
): UseRequestMutationReturn<T, D>;
export function useRequest<T, D = void>(
  _url: string | null,
  _options?: UseRequestOptions<T> & { method?: MutationMethod }
): UseRequestGetReturn<T> | UseRequestMutationReturn<T, D> {
  // 注: この関数は型チェック用であり、実行時には使用しない
  // 実際のコードでは useRequestGet または useRequestMutation を直接使用すること
  throw new Error(
    "useRequest は非推奨です。useRequestGet または useRequestMutation を使用してください。"
  );
}
