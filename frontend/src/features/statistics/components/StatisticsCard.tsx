/**
 * StatisticsCard - 統計サマリーカード
 * カウントアップアニメーションとグラデーションスタイルを使用
 */
import { Card, CardContent } from "@/components/ui/card";
import { useCountUp } from "@/features/common/hooks";
import type { StatisticsResponse } from "../hooks/useStatistics";

export type StatisticsCardProps = {
  data: StatisticsResponse;
};

export function StatisticsCard({ data }: StatisticsCardProps) {
  const achievementRate = data.totalDays > 0
    ? Math.round((data.achievedDays / data.totalDays) * 100)
    : 0;

  // カウントアップアニメーション用の値
  const animatedTarget = useCountUp({
    end: data.targetCalories,
    duration: 1000,
    startOnMount: true,
  });
  const animatedAverage = useCountUp({
    end: data.averageCalories,
    duration: 1000,
    startOnMount: true,
  });
  const animatedAchieved = useCountUp({
    end: data.achievedDays,
    duration: 1000,
    startOnMount: true,
  });
  const animatedTotal = useCountUp({
    end: data.totalDays,
    duration: 1000,
    startOnMount: true,
  });
  const animatedOver = useCountUp({
    end: data.overDays,
    duration: 1000,
    startOnMount: true,
  });
  const animatedRate = useCountUp({
    end: achievementRate,
    duration: 1000,
    startOnMount: true,
  });

  return (
    <div className="grid gap-4 md:grid-cols-4">
      {/* 目標カロリー */}
      <Card className="opacity-0 animate-fade-in-up transition-transform hover:scale-[1.02] hover:shadow-lg">
        <CardContent className="h-full flex flex-col items-center justify-center pt-6 pb-6">
          <p className="text-sm text-muted-foreground mb-2">目標</p>
          <p className="text-4xl font-extrabold tabular-nums tracking-tight bg-gradient-to-br from-slate-700 to-slate-500 bg-clip-text text-transparent drop-shadow-sm">
            {animatedTarget.toLocaleString()}
          </p>
          <span className="text-xs font-medium text-muted-foreground mt-1">
            kcal/日
          </span>
        </CardContent>
      </Card>

      {/* 平均摂取カロリー */}
      <Card className="opacity-0 animate-fade-in-up animation-delay-100 transition-transform hover:scale-[1.02] hover:shadow-lg">
        <CardContent className="h-full flex flex-col items-center justify-center pt-6 pb-6">
          <p className="text-sm text-muted-foreground mb-2">平均摂取</p>
          <p className="text-4xl font-extrabold tabular-nums tracking-tight bg-gradient-to-br from-blue-600 to-cyan-500 bg-clip-text text-transparent drop-shadow-sm">
            {animatedAverage.toLocaleString()}
          </p>
          <span className="text-xs font-medium text-muted-foreground mt-1">
            kcal/日
          </span>
        </CardContent>
      </Card>

      {/* 達成日数 */}
      <Card className="opacity-0 animate-fade-in-up animation-delay-200 transition-transform hover:scale-[1.02] hover:shadow-lg">
        <CardContent className="h-full flex flex-col items-center justify-center pt-6 pb-6">
          <p className="text-sm text-muted-foreground mb-2">達成日数</p>
          <p className="text-4xl font-extrabold tabular-nums tracking-tight bg-gradient-to-br from-emerald-500 to-teal-400 bg-clip-text text-transparent drop-shadow-sm">
            {animatedAchieved}
            <span className="text-lg font-normal text-muted-foreground">
              /{animatedTotal}
            </span>
          </p>
          <span className="text-xs font-medium text-muted-foreground mt-1">
            日 ({animatedRate}%)
          </span>
        </CardContent>
      </Card>

      {/* 超過日数 */}
      <Card className="opacity-0 animate-fade-in-up animation-delay-300 transition-transform hover:scale-[1.02] hover:shadow-lg">
        <CardContent className="h-full flex flex-col items-center justify-center pt-6 pb-6">
          <p className="text-sm text-muted-foreground mb-2">超過日数</p>
          <p
            className={`text-4xl font-extrabold tabular-nums tracking-tight bg-clip-text text-transparent drop-shadow-sm ${
              data.overDays > 0
                ? "bg-gradient-to-br from-red-500 to-orange-500"
                : "bg-gradient-to-br from-gray-400 to-gray-300"
            }`}
          >
            {animatedOver}
          </p>
          <span className="text-xs font-medium text-muted-foreground mt-1">
            日
          </span>
        </CardContent>
      </Card>
    </div>
  );
}
