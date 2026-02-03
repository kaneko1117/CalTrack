/**
 * useStatistics - カロリー統計データ取得フック
 */
import { useRequestGet } from "@/features/common/hooks";

/** 期間タイプ */
export type Period = "week" | "month";

/** 日別統計データ */
export type DailyStatistic = {
  date: string;
  totalCalories: number;
};

/** 統計レスポンス型 */
export type StatisticsResponse = {
  period: Period;
  targetCalories: number;
  averageCalories: number;
  totalDays: number;
  achievedDays: number;
  overDays: number;
  dailyStatistics: DailyStatistic[];
};

/** 期間ラベル */
export const PERIOD_LABELS: Record<Period, string> = {
  week: "週間",
  month: "月間",
};

/** エラーメッセージ */
export const ERROR_MESSAGE_FETCH_FAILED = "統計データの取得に失敗しました";

export function useStatistics(period: Period) {
  const { data, error, isLoading, mutate } = useRequestGet<StatisticsResponse>(
    `/api/v1/statistics?period=${period}`
  );

  return { data, error, isLoading, refetch: mutate };
}
