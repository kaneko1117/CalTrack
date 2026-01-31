/**
 * HomePage - ホームページコンポーネント
 *
 * 既存のApp.tsxの内容をページコンポーネントとして切り出し
 */
import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

/** ヘルスチェックのレスポンス型 */
type HealthStatus = {
  status: string;
  message: string;
};

/**
 * HomePage - ホームページ
 *
 * バックエンドのヘルスチェック機能を提供
 */
export function HomePage() {
  const [health, setHealth] = useState<HealthStatus | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  /**
   * バックエンドのヘルスチェックを実行
   */
  const checkHealth = async () => {
    setLoading(true);
    setError(null);
    try {
      const response = await fetch("http://localhost:8080/health");
      const data = await response.json();
      setHealth(data);
    } catch {
      setError("Failed to connect to backend");
      setHealth(null);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    checkHealth();
  }, []);

  return (
    <div className="min-h-screen bg-background flex items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>CalTrack</CardTitle>
          <CardDescription>Calorie tracking application</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="text-sm">
            <p className="font-medium">Backend Status:</p>
            {loading && <p className="text-muted-foreground">Checking...</p>}
            {error && <p className="text-destructive">{error}</p>}
            {health && (
              <p className="text-green-600">
                {health.status}: {health.message}
              </p>
            )}
          </div>
          <Button onClick={checkHealth} disabled={loading}>
            {loading ? "Checking..." : "Check Backend Health"}
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}
