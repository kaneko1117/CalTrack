/**
 * DashboardPage - ダッシュボードページ
 * ユーザーの今日のカロリー記録を表示し、新規記録の追加が可能
 */
import { Header } from "@/components/Header";
import { Footer } from "@/components/Footer";
import { RecordDialog } from "@/features/records";
import { Card, CardContent } from "@/components/ui/card";

/**
 * DashboardPage - ダッシュボードページコンポーネント
 */
export function DashboardPage() {
  const handleRecordSuccess = () => {
    // 記録成功時の処理（今後、記録一覧の再取得などを実装予定）
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
          <Card>
            <CardContent className="pt-6">
              <p className="text-sm text-muted-foreground">今日の合計</p>
              <p className="text-4xl font-bold text-primary">0 kcal</p>
            </CardContent>
          </Card>
        </div>
      </main>
      <Footer />
    </div>
  );
}
