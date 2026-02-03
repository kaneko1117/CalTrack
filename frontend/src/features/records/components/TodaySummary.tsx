/**
 * TodaySummary - 今日のカロリーサマリーコンポーネント
 * 目標・摂取・残りカロリーを表示する
 */
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import type { ApiErrorResponse } from "@/lib/api";
import type { TodayRecordsResponse } from "../hooks/useTodayRecords";

export type TodaySummaryProps = {
  data: TodayRecordsResponse | null;
  isPending: boolean;
  error: ApiErrorResponse | null;
};

/**
 * TodaySummary - 今日のカロリーサマリーコンポーネント
 */
export function TodaySummary({ data, isPending, error }: TodaySummaryProps) {
  // ローディング状態
  if (isPending && !data) {
    return (
      <div className="grid gap-4 md:grid-cols-3">
        {[...Array(3)].map((_, i) => (
          <Card key={i}>
            <CardContent className="pt-6">
              <Skeleton className="h-4 w-20 mb-2" />
              <Skeleton className="h-10 w-32" />
            </CardContent>
          </Card>
        ))}
      </div>
    );
  }

  // エラー状態
  if (error) {
    return (
      <Card>
        <CardContent className="pt-6">
          <p className="text-sm text-destructive">データの取得に失敗しました</p>
        </CardContent>
      </Card>
    );
  }

  // データがない場合
  if (!data) {
    return null;
  }

  const remainingCalories = data.targetCalories - data.totalCalories;
  const isOverTarget = remainingCalories < 0;

  return (
    <div className="grid gap-4 md:grid-cols-3">
      {/* 目標カロリー */}
      <Card>
        <CardContent className="pt-6">
          <p className="text-sm text-muted-foreground">目標</p>
          <p className="text-3xl font-bold">
            {data.targetCalories.toLocaleString()}
            <span className="text-base font-normal text-muted-foreground ml-1">kcal</span>
          </p>
        </CardContent>
      </Card>

      {/* 摂取カロリー */}
      <Card>
        <CardContent className="pt-6">
          <p className="text-sm text-muted-foreground">摂取</p>
          <p className="text-3xl font-bold text-primary">
            {data.totalCalories.toLocaleString()}
            <span className="text-base font-normal text-muted-foreground ml-1">kcal</span>
          </p>
        </CardContent>
      </Card>

      {/* 残り/超過カロリー */}
      <Card>
        <CardContent className="pt-6">
          <p className="text-sm text-muted-foreground">{isOverTarget ? "超過" : "残り"}</p>
          <p className={`text-3xl font-bold ${isOverTarget ? "text-destructive" : "text-green-600"}`}>
            {Math.abs(remainingCalories).toLocaleString()}
            <span className="text-base font-normal text-muted-foreground ml-1">kcal</span>
          </p>
        </CardContent>
      </Card>
    </div>
  );
}
