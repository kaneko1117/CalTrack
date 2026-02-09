/**
 * useTodayPfc - 今日のPFC摂取量・目標値取得フック
 */
import { useRequestGet } from "@/features/common/hooks";

/** PFC栄養素の値 */
export type PfcValues = {
  protein: number;
  fat: number;
  carbs: number;
};

/** 今日のPFCレスポンス型 */
export type TodayPfcResponse = {
  date: string;
  current: PfcValues;
  target: PfcValues;
};

/** エラーメッセージ */
export const ERROR_MESSAGE_PFC_FETCH_FAILED = "PFCデータの取得に失敗しました";

/**
 * 今日のPFC摂取量・目標値を取得するフック
 * @returns { data, error, isLoading, refetch }
 */
export function useTodayPfc() {
  const { data, error, isLoading, isValidating, mutate } =
    useRequestGet<TodayPfcResponse>("/api/v1/nutrition/today-pfc");

  return { data, error, isLoading: isLoading || isValidating, refetch: mutate };
}
