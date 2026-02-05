/**
 * useNutritionAdvice - PFC栄養アドバイス取得フック
 */
import { useRequestGet } from "@/features/common/hooks";

/** 栄養アドバイスレスポンス型 */
export type NutritionAdviceResponse = {
  advice: string;
};

/** エラーメッセージ */
export const ERROR_MESSAGE_FETCH_FAILED = "アドバイスの取得に失敗しました";

/**
 * 栄養アドバイスを取得するフック
 * @returns { data, error, isLoading, refetch }
 */
export function useNutritionAdvice() {
  const { data, error, isLoading, mutate } =
    useRequestGet<NutritionAdviceResponse>("/api/v1/nutrition/advice");

  return { data, error, isLoading, refetch: mutate };
}
