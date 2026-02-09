/**
 * PfcProgressCard - 今日のPFC摂取量プログレスカード
 * P(タンパク質)、F(脂質)、C(炭水化物)それぞれの進捗バーを表示する
 */
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { Skeleton } from "@/components/ui/skeleton";
import { PFC_NUTRIENT_OPTIONS } from "@/domain/valueObjects";
import { useCountUp } from "@/features/common/hooks";
import type { ApiErrorResponse } from "@/lib/api";
import type { TodayPfcResponse } from "../hooks/useTodayPfc";

export type PfcProgressCardProps = {
  data: TodayPfcResponse | null;
  isLoading: boolean;
  error: ApiErrorResponse | null;
};

/** PFC各栄養素の表示設定 */
type NutrientConfig = {
  key: "protein" | "fat" | "carbs";
  label: string;
  shortLabel: string;
  colorClass: string;
  bgColorClass: string;
};

const NUTRIENT_CONFIGS: NutrientConfig[] = PFC_NUTRIENT_OPTIONS.map((opt) => {
  const colors: Record<string, { colorClass: string; bgColorClass: string }> = {
    protein: { colorClass: "bg-blue-500", bgColorClass: "bg-blue-100" },
    fat: { colorClass: "bg-amber-500", bgColorClass: "bg-amber-100" },
    carbs: { colorClass: "bg-emerald-500", bgColorClass: "bg-emerald-100" },
  };
  return {
    key: opt.value,
    label: opt.label,
    shortLabel: opt.shortLabel,
    ...colors[opt.value],
  };
});

/** 進捗ステータス */
type ProgressStatus = "normal" | "optimal" | "over";

/** 進捗率からステータスを判定する */
function getProgressStatus(rawPercentage: number): ProgressStatus {
  if (rawPercentage > 100) return "over";
  if (rawPercentage >= 80) return "optimal";
  return "normal";
}

/**
 * 単一栄養素の進捗バー表示コンポーネント
 */
function NutrientProgressBar({
  config,
  current,
  target,
}: {
  config: NutrientConfig;
  current: number;
  target: number;
}) {
  const rawPercentage = target > 0 ? Math.round((current / target) * 100) : 0;
  const clampedPercentage = Math.min(rawPercentage, 100);
  const status = getProgressStatus(rawPercentage);

  const animatedCurrent = useCountUp({ end: Math.round(current), duration: 1000, startOnMount: true });
  const animatedPercentage = useCountUp({ end: rawPercentage, duration: 1000, startOnMount: true });

  // ステータスに応じたスタイル
  const badgeColorClass = status === "over" ? "bg-red-500" : status === "optimal" ? "bg-green-500" : config.colorClass;
  const barTrackClass = status === "over" ? "bg-red-100" : status === "optimal" ? "bg-green-100" : config.bgColorClass;
  const barIndicatorClass = status === "over" ? "bg-red-500" : status === "optimal" ? "bg-green-500" : config.colorClass;
  const textColorClass = status === "over" ? "text-red-600" : status === "optimal" ? "text-green-600" : "text-muted-foreground";

  return (
    <div className="space-y-2">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <span
            className={`inline-flex items-center justify-center w-6 h-6 rounded-full text-xs font-bold text-white ${badgeColorClass}`}
          >
            {config.shortLabel}
          </span>
          <span className="text-sm font-medium">{config.label}</span>
        </div>
        <span className={`text-sm tabular-nums ${textColorClass}`}>
          {animatedCurrent}g / {Math.round(target)}g ({animatedPercentage}%)
        </span>
      </div>
      <Progress
        value={clampedPercentage}
        className={`h-3 ${barTrackClass}`}
        indicatorClassName={barIndicatorClass}
        aria-label={`${config.label}の進捗`}
      />
    </div>
  );
}

/**
 * PfcProgressCard - PFC摂取量プログレスカード
 */
export function PfcProgressCard({ data, isLoading, error }: PfcProgressCardProps) {
  if (isLoading) {
    return (
      <Card className="opacity-0 animate-fade-in-up">
        <CardHeader className="pb-3">
          <CardTitle className="text-lg">今日のPFCバランス</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          {[0, 1, 2].map((i) => (
            <div key={i} className="space-y-2">
              <div className="flex items-center justify-between">
                <Skeleton className="h-4 w-24" />
                <Skeleton className="h-4 w-32" />
              </div>
              <Skeleton className="h-3 w-full rounded-full" />
            </div>
          ))}
        </CardContent>
      </Card>
    );
  }

  if (error) {
    return (
      <Card className="opacity-0 animate-fade-in-up">
        <CardHeader className="pb-3">
          <CardTitle className="text-lg">今日のPFCバランス</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-destructive">PFCデータの取得に失敗しました</p>
        </CardContent>
      </Card>
    );
  }

  if (!data) {
    return (
      <Card className="opacity-0 animate-fade-in-up">
        <CardHeader className="pb-3">
          <CardTitle className="text-lg">今日のPFCバランス</CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">
            食事を記録するとPFCバランスが表示されます
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="opacity-0 animate-fade-in-up">
      <CardHeader className="pb-3">
        <CardTitle className="text-lg">今日のPFCバランス</CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {NUTRIENT_CONFIGS.map((config) => (
          <NutrientProgressBar
            key={config.key}
            config={config}
            current={data.current[config.key]}
            target={data.target[config.key]}
          />
        ))}
      </CardContent>
    </Card>
  );
}
