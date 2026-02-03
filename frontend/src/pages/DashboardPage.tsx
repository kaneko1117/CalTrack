/**
 * DashboardPage - ダッシュボードページ
 * ユーザーの今日のカロリー記録を表示し、新規記録の追加が可能
 */
import { Header } from "@/components/Header";
import { Footer } from "@/components/Footer";
import { RecordDialog, TodaySummary, useTodayRecords } from "@/features/records";

/**
 * DashboardPage - ダッシュボードページコンポーネント
 */
export function DashboardPage() {
  const { data, error, isLoading, refetch } = useTodayRecords();

  /**
   * 記録成功時のコールバック
   * データを再取得する
   */
  const handleRecordSuccess = () => {
    refetch();
  };

  return (
    <div className="min-h-screen flex flex-col bg-background">
      <Header />
      <main className="flex-1 container px-4 py-8 mx-auto max-w-2xl">
        <div className="space-y-6">
          <div className="flex items-center justify-between">
            <h1 className="text-2xl font-bold">今日のカロリー記録</h1>
            <RecordDialog onSuccess={handleRecordSuccess} />
          </div>
          <TodaySummary data={data ?? null} isLoading={isLoading} error={error ?? null} />
        </div>
      </main>
      <Footer />
    </div>
  );
}
