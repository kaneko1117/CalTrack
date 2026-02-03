/**
 * useTodayRecords - 今日のカロリー記録取得フック
 */
import { useRequest } from "@/features/common/hooks";

type RecordItem = {
  itemId: string;
  name: string;
  calories: number;
};

type Record = {
  id: string;
  eatenAt: string;
  items: RecordItem[];
};

export type TodayRecordsResponse = {
  date: string;
  totalCalories: number;
  targetCalories: number;
  difference: number;
  records: Record[];
};

export function useTodayRecords() {
  const { data, error, isLoading, mutate } = useRequest<TodayRecordsResponse>(
    "/api/v1/records/today"
  );

  return {
    data,
    error,
    isLoading,
    refetch: mutate,
  };
}
