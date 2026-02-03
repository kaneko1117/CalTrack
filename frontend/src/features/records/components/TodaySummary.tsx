/**
 * TodaySummary - 今日のカロリーサマリーコンポーネント
 * 目標・摂取・残りカロリーと記録一覧を表示する
 */
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import type { ApiErrorResponse } from "@/lib/api";
import type { TodayRecordsResponse } from "../hooks/useTodayRecords";
import { newEatenAt } from "@/domain/valueObjects/eatenAt";

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
      <div className="space-y-6">
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
        <Card>
          <CardHeader>
            <Skeleton className="h-6 w-24" />
          </CardHeader>
          <CardContent>
            <Skeleton className="h-4 w-full mb-2" />
            <Skeleton className="h-4 w-3/4" />
          </CardContent>
        </Card>
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
    <div className="space-y-6">
      {/* サマリーカード */}
      <div className="grid gap-4 md:grid-cols-3">
        {/* 目標カロリー */}
        <Card>
          <CardContent className="pt-6">
            <p className="text-sm text-muted-foreground">目標</p>
            <p className="text-3xl font-bold">
              {data.targetCalories.toLocaleString()}
              <span className="text-base font-normal text-muted-foreground ml-1">
                kcal
              </span>
            </p>
          </CardContent>
        </Card>

        {/* 摂取カロリー */}
        <Card>
          <CardContent className="pt-6">
            <p className="text-sm text-muted-foreground">摂取</p>
            <p className="text-3xl font-bold text-primary">
              {data.totalCalories.toLocaleString()}
              <span className="text-base font-normal text-muted-foreground ml-1">
                kcal
              </span>
            </p>
          </CardContent>
        </Card>

        {/* 残り/超過カロリー */}
        <Card>
          <CardContent className="pt-6">
            <p className="text-sm text-muted-foreground">
              {isOverTarget ? "超過" : "残り"}
            </p>
            <p
              className={`text-3xl font-bold ${isOverTarget ? "text-destructive" : "text-green-600"}`}
            >
              {Math.abs(remainingCalories).toLocaleString()}
              <span className="text-base font-normal text-muted-foreground ml-1">
                kcal
              </span>
            </p>
          </CardContent>
        </Card>
      </div>

      {/* 記録一覧 */}
      <Card>
        <CardHeader></CardHeader>
        <CardContent>
          {data.records.length === 0 ? (
            <p className="text-sm text-muted-foreground">
              まだ記録がありません
            </p>
          ) : (
            <div className="space-y-4">
              {data.records.map((record) => {
                const eatenAtResult = newEatenAt(record.eatenAt);
                if (!eatenAtResult.ok) return null;
                const eatenAt = eatenAtResult.value;

                return (
                  <div
                    key={record.id}
                    className="border-b last:border-b-0 pb-4 last:pb-0"
                  >
                    <p className="text-sm font-medium text-muted-foreground mb-2">
                      {eatenAt.formattedTime()}
                      <span className="ml-2 text-xs bg-muted px-2 py-0.5 rounded">
                        {eatenAt.mealTypeLabel()}
                      </span>
                    </p>
                    <ul className="space-y-1">
                      {record.items.map((item, index) => (
                        <li
                          key={item.itemId}
                          className="flex justify-between items-center text-sm"
                        >
                          <span className="flex items-center">
                            <span className="text-muted-foreground mr-2">
                              {index === record.items.length - 1 ? "└" : "├"}
                            </span>
                            {item.name}
                          </span>
                          <span className="text-muted-foreground">
                            {item.calories.toLocaleString()} kcal
                          </span>
                        </li>
                      ))}
                    </ul>
                  </div>
                );
              })}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
