/**
 * CalorieChart - カロリー推移チャート
 */
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  type ChartConfig,
} from "@/components/ui/chart";
import { Line, LineChart, XAxis, YAxis, CartesianGrid, ReferenceLine } from "recharts";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import type { DailyStatistic } from "../hooks/useStatistics";

export type CalorieChartProps = {
  data: DailyStatistic[];
  targetCalories: number;
  isLoading?: boolean;
};

const chartConfig = {
  calories: { label: "摂取カロリー", color: "hsl(221.2 83.2% 53.3%)" },
  target: { label: "目標", color: "hsl(0 84.2% 60.2%)" },
} satisfies ChartConfig;

function formatDateLabel(dateString: string): string {
  const date = new Date(dateString);
  return `${date.getMonth() + 1}/${date.getDate()}`;
}

export function CalorieChart({ data, targetCalories, isLoading = false }: CalorieChartProps) {
  if (isLoading) {
    return (
      <Card>
        <CardHeader><CardTitle>カロリー推移</CardTitle></CardHeader>
        <CardContent><Skeleton className="h-[300px] w-full" /></CardContent>
      </Card>
    );
  }

  if (data.length === 0) {
    return (
      <Card>
        <CardHeader><CardTitle>カロリー推移</CardTitle></CardHeader>
        <CardContent>
          <div className="h-[300px] flex items-center justify-center">
            <p className="text-sm text-muted-foreground">表示するデータがありません</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  const chartData = data.map((item) => ({
    dateLabel: formatDateLabel(item.date),
    calories: item.totalCalories,
  }));

  const maxCalories = Math.max(...data.map((d) => d.totalCalories), targetCalories);
  const yAxisMax = Math.ceil((maxCalories * 1.2) / 100) * 100;

  return (
    <Card className="opacity-0 animate-fade-in-up">
      <CardHeader><CardTitle>カロリー推移</CardTitle></CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig} className="h-[300px] w-full">
          <LineChart data={chartData} margin={{ top: 20, right: 20, left: 20, bottom: 20 }}>
            <CartesianGrid strokeDasharray="3 3" className="stroke-muted" />
            <XAxis dataKey="dateLabel" tickLine={false} axisLine={false} tickMargin={8} />
            <YAxis domain={[0, yAxisMax]} tickLine={false} axisLine={false} tickMargin={8} />
            <ChartTooltip
              content={<ChartTooltipContent formatter={(value) => <span className="font-medium">{Number(value).toLocaleString()} kcal</span>} />}
            />
            <ReferenceLine
              y={targetCalories}
              stroke="var(--color-target)"
              strokeDasharray="5 5"
              strokeWidth={2}
              label={{ value: `目標: ${targetCalories.toLocaleString()} kcal`, position: "insideTopRight", fill: "hsl(0 84.2% 60.2%)", fontSize: 12 }}
            />
            <Line
              type="monotone"
              dataKey="calories"
              stroke="var(--color-calories)"
              strokeWidth={2}
              dot={{ fill: "var(--color-calories)", strokeWidth: 2, r: 4 }}
              activeDot={{ r: 6, strokeWidth: 2 }}
            />
          </LineChart>
        </ChartContainer>
        <div className="flex items-center justify-center gap-6 mt-4">
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 rounded-full" style={{ backgroundColor: chartConfig.calories.color }} />
            <span className="text-sm text-muted-foreground">摂取カロリー</span>
          </div>
          <div className="flex items-center gap-2">
            <div className="w-3 h-0.5" style={{ backgroundColor: chartConfig.target.color }} />
            <span className="text-sm text-muted-foreground">目標</span>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}
