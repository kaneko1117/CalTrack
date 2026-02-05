/**
 * DashboardPage - ダッシュボードページ
 * ユーザーの今日のカロリー記録を表示し、新規記録の追加が可能
 * 統計データ（期間選択・統計カード・カロリーチャート）も表示
 */
import { useState } from "react";
import { Header } from "@/components/Header";
import { Footer } from "@/components/Footer";
import { RecordDialog, TodaySummary, useTodayRecords } from "@/features/records";
import {
  useStatistics,
  ERROR_MESSAGE_FETCH_FAILED,
  type Period,
} from "@/features/statistics/hooks/useStatistics";
import { PeriodSelector } from "@/features/statistics/components/PeriodSelector";
import { StatisticsCard } from "@/features/statistics/components/StatisticsCard";
import { CalorieChart } from "@/features/statistics/components/CalorieChart";
import { useNutritionAdvice, NutritionAdviceCard } from "@/features/nutrition";

/**
 * DashboardPage - ダッシュボードページコンポーネント
 */
export function DashboardPage() {
  const { data, error, isLoading, refetch } = useTodayRecords();
  const [period, setPeriod] = useState<Period>("week");
  const {
    data: statisticsData,
    error: statisticsError,
    isLoading: statisticsLoading,
    refetch: statisticsRefetch,
  } = useStatistics(period);
  const {
    data: adviceData,
    error: adviceError,
    isLoading: adviceLoading,
    refetch: adviceRefetch,
  } = useNutritionAdvice();

  /**
   * 記録成功時のコールバック
   * データを再取得する
   */
  const handleRecordSuccess = () => {
    refetch();
    statisticsRefetch();
    adviceRefetch();
  };

  return (
    <div className="min-h-screen flex flex-col bg-background">
      <Header />
      <main className="flex-1 container px-4 py-8 mx-auto max-w-3xl">
        <div className="space-y-8">
          {/* 今日のカロリー記録セクション */}
          <section className="space-y-6">
            <div className="flex items-center justify-between">
              <h1 className="text-2xl font-bold">今日のカロリー記録</h1>
              <RecordDialog onSuccess={handleRecordSuccess} />
            </div>
            <TodaySummary data={data ?? null} isLoading={isLoading} error={error ?? null} />
          </section>

          {/* PFCアドバイスセクション */}
          <section>
            <NutritionAdviceCard
              advice={adviceData?.advice ?? null}
              isLoading={adviceLoading}
              error={adviceError ?? null}
            />
          </section>

          {/* 統計データセクション */}
          <section className="space-y-6">
            <div className="flex items-center justify-between">
              <h2 className="text-2xl font-bold">カロリー統計</h2>
              <PeriodSelector value={period} onChange={setPeriod} />
            </div>
            {statisticsError ? (
              <p className="text-sm text-destructive">{ERROR_MESSAGE_FETCH_FAILED}</p>
            ) : statisticsData ? (
              <div className="space-y-6">
                <StatisticsCard data={statisticsData} />
                <CalorieChart
                  data={statisticsData.dailyStatistics}
                  targetCalories={statisticsData.targetCalories}
                  isLoading={statisticsLoading}
                />
              </div>
            ) : null}
          </section>
        </div>
      </main>
      <Footer />
    </div>
  );
}
