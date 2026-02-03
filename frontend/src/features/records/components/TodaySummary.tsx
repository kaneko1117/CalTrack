/**
 * TodaySummary - 今日のカロリーサマリーコンポーネント
 * 目標・摂取・残りカロリーと記録一覧を表示する
 */
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import type { ApiErrorResponse } from "@/lib/api";
import type { TodayRecordsResponse } from "../hooks/useTodayRecords";
import { newEatenAt, type MealType } from "@/domain/valueObjects/eatenAt";
import { ProgressRing } from "./ProgressRing";
import { useCountUp } from "../hooks/useCountUp";

/**
 * 食事タイプごとの絵文字マップ
 */
const MEAL_TYPE_EMOJI: Record<MealType, string> = {
  breakfast: "🌅",
  lunch: "☀️",
  snack: "🍪",
  dinner: "🌙",
  lateNight: "🌃",
};

/**
 * mealTypeに応じたスタイルクラスを返す
 */
function getMealTypeStyle(mealType: MealType): string {
  const styles: Record<MealType, string> = {
    breakfast: "bg-orange-100 text-orange-700",
    lunch: "bg-green-100 text-green-700",
    snack: "bg-purple-100 text-purple-700",
    dinner: "bg-blue-100 text-blue-700",
    lateNight: "bg-gray-100 text-gray-700",
  };
  return styles[mealType];
}

export type TodaySummaryProps = {
  data: TodayRecordsResponse | null;
  isLoading: boolean;
  error: ApiErrorResponse | null;
};

/**
 * TodaySummary - 今日のカロリーサマリーコンポーネント
 */
export function TodaySummary({ data, isLoading, error }: TodaySummaryProps) {
  // カウントアップアニメーション用の値
  const animatedTarget = useCountUp({
    end: data?.targetCalories ?? 0,
    duration: 1000,
    startOnMount: !!data,
  });
  const animatedTotal = useCountUp({
    end: data?.totalCalories ?? 0,
    duration: 1000,
    startOnMount: !!data,
  });
  const animatedDifference = useCountUp({
    end: Math.abs(data?.difference ?? 0),
    duration: 1000,
    startOnMount: !!data,
  });

  // ローディング状態
  if (isLoading && !data) {
    return (
      <div className="space-y-6">
        <div className="grid gap-4 md:grid-cols-4">
          {[...Array(4)].map((_, i) => (
            <Card key={i}>
              <CardContent className="pt-6 pb-6">
                <Skeleton className="h-4 w-20 mb-2" />
                <Skeleton className="h-10 w-32" />
              </CardContent>
            </Card>
          ))}
        </div>
        <Card>
          <CardContent className="pt-6">
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

  const isOverTarget = data.difference < 0;
  const progressPercent = Math.min(
    (data.totalCalories / data.targetCalories) * 100,
    100
  );

  return (
    <div className="space-y-6">
      {/* サマリーカード */}
      <div className="grid gap-4 md:grid-cols-4">
        {/* プログレスリング */}
        <Card className="opacity-0 animate-fade-in-up transition-transform hover:scale-[1.02] hover:shadow-lg">
          <CardContent className="h-full flex flex-col items-center justify-center pt-6 pb-6">
            <p className="text-sm text-muted-foreground mb-3">達成率</p>
            <ProgressRing progress={progressPercent} />
          </CardContent>
        </Card>

        {/* 目標カロリー */}
        <Card className="opacity-0 animate-fade-in-up animation-delay-100 transition-transform hover:scale-[1.02] hover:shadow-lg">
          <CardContent className="h-full flex flex-col items-center justify-center pt-6 pb-6">
            <p className="text-sm text-muted-foreground mb-2">目標</p>
            <p className="text-4xl font-extrabold tabular-nums tracking-tight bg-gradient-to-br from-slate-700 to-slate-500 bg-clip-text text-transparent drop-shadow-sm">
              {animatedTarget.toLocaleString()}
            </p>
            <span className="text-xs font-medium text-muted-foreground mt-1">
              kcal
            </span>
          </CardContent>
        </Card>

        {/* 摂取カロリー */}
        <Card className="opacity-0 animate-fade-in-up animation-delay-200 transition-transform hover:scale-[1.02] hover:shadow-lg">
          <CardContent className="h-full flex flex-col items-center justify-center pt-6 pb-6">
            <p className="text-sm text-muted-foreground mb-2">摂取</p>
            <p className="text-4xl font-extrabold tabular-nums tracking-tight bg-gradient-to-br from-blue-600 to-cyan-500 bg-clip-text text-transparent drop-shadow-sm">
              {animatedTotal.toLocaleString()}
            </p>
            <span className="text-xs font-medium text-muted-foreground mt-1">
              kcal
            </span>
          </CardContent>
        </Card>

        {/* 残り/超過カロリー */}
        <Card className="opacity-0 animate-fade-in-up animation-delay-300 transition-transform hover:scale-[1.02] hover:shadow-lg">
          <CardContent className="h-full flex flex-col items-center justify-center pt-6 pb-6">
            <p className="text-sm text-muted-foreground mb-2">
              {isOverTarget ? "超過" : "残り"}
            </p>
            <p
              className={`text-4xl font-extrabold tabular-nums tracking-tight bg-clip-text text-transparent drop-shadow-sm ${
                isOverTarget
                  ? "bg-gradient-to-br from-red-500 to-orange-500"
                  : "bg-gradient-to-br from-emerald-500 to-teal-400"
              }`}
            >
              {animatedDifference.toLocaleString()}
            </p>
            <span className="text-xs font-medium text-muted-foreground mt-1">
              kcal
            </span>
          </CardContent>
        </Card>
      </div>

      {/* 記録一覧 */}
      <Card className="opacity-0 animate-fade-in-up animation-delay-400">
        <CardContent className="pt-6">
          {data.records.length === 0 ? (
            <p className="text-sm text-muted-foreground">
              まだ記録がありません
            </p>
          ) : (
            <div className="space-y-4">
              {data.records.map((record, recordIndex) => {
                const eatenAtResult = newEatenAt(record.eatenAt);
                if (!eatenAtResult.ok) return null;
                const eatenAt = eatenAtResult.value;
                const mealType = eatenAt.mealType();

                return (
                  <div
                    key={record.id}
                    className="border-b last:border-b-0 pb-4 last:pb-0 opacity-0 animate-fade-in"
                    style={{ animationDelay: `${500 + recordIndex * 100}ms` }}
                  >
                    <p className="text-sm font-medium text-muted-foreground mb-2">
                      <span className="mr-1">{MEAL_TYPE_EMOJI[mealType]}</span>
                      {eatenAt.formattedTime()}
                      <span
                        className={`ml-2 text-xs px-2 py-0.5 rounded ${getMealTypeStyle(mealType)}`}
                      >
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
