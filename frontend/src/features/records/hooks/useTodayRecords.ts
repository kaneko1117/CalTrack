/**
 * useTodayRecords - 今日のカロリー記録取得フック
 * 今日のカロリーサマリーと記録一覧を取得する
 */
import { useCallback } from "react";
import { get } from "@/lib/api";
import { useApi } from "@/features/common/hooks/useApi";

/** 記録内のアイテム */
type RecordItem = {
  itemId: string;
  name: string;
  calories: number;
};

/** 記録 */
type Record = {
  id: string;
  eatenAt: string;
  items: RecordItem[];
};

/** 今日のカロリー記録レスポンス */
export type TodayRecordsResponse = {
  date: string;
  totalCalories: number;
  targetCalories: number;
  difference: number;
  records: Record[];
};

/**
 * 今日のカロリー記録を取得するAPI関数
 */
const getTodayRecords = () =>
  get<TodayRecordsResponse>("/api/v1/records/today");

/**
 * useTodayRecords - 今日のカロリー記録取得フック
 * @returns { data, error, isPending, fetch }
 */
export function useTodayRecords() {
  const { execute, data, error, isPending } = useApi(getTodayRecords);

  const fetch = useCallback(() => {
    execute();
  }, [execute]);

  return { data, error, isPending, fetch };
}
