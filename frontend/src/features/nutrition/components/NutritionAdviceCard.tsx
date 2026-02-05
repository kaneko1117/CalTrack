/**
 * NutritionAdviceCard - PFC栄養アドバイス表示コンポーネント
 * AIが生成したPFCバランスに基づくアドバイスを表示する
 */
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import type { ApiErrorResponse } from "@/lib/api";

export type NutritionAdviceCardProps = {
  advice: string | null;
  isLoading: boolean;
  error: ApiErrorResponse | null;
};

/**
 * NutritionAdviceCard - 栄養アドバイスカードコンポーネント
 */
export function NutritionAdviceCard({
  advice,
  isLoading,
  error,
}: NutritionAdviceCardProps) {
  // ローディング状態
  if (isLoading) {
    return (
      <Card className="opacity-0 animate-fade-in-up">
        <CardHeader className="pb-3">
          <CardTitle className="text-lg flex items-center gap-2">
            <span>💡</span>
            <Skeleton className="h-5 w-32" />
          </CardTitle>
        </CardHeader>
        <CardContent>
          <Skeleton className="h-4 w-full mb-2" />
          <Skeleton className="h-4 w-5/6 mb-2" />
          <Skeleton className="h-4 w-4/6" />
        </CardContent>
      </Card>
    );
  }

  // エラー状態
  if (error) {
    return (
      <Card className="opacity-0 animate-fade-in-up">
        <CardHeader className="pb-3">
          <CardTitle className="text-lg flex items-center gap-2">
            <span>💡</span>
            AIによるアドバイス
          </CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-destructive">
            アドバイスの取得に失敗しました
          </p>
        </CardContent>
      </Card>
    );
  }

  // アドバイスがない場合
  if (!advice) {
    return (
      <Card className="opacity-0 animate-fade-in-up">
        <CardHeader className="pb-3">
          <CardTitle className="text-lg flex items-center gap-2">
            <span>💡</span>
            AIによるアドバイス
          </CardTitle>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground">
            食事を記録するとアドバイスが表示されます
          </p>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="opacity-0 animate-fade-in-up bg-gradient-to-br from-amber-50 to-orange-50 border-amber-200">
      <CardHeader className="pb-3">
        <CardTitle className="text-lg flex items-center gap-2 text-amber-800">
          <span>💡</span>
          AIによるアドバイス
        </CardTitle>
      </CardHeader>
      <CardContent>
        <p className="text-sm text-amber-900 leading-relaxed whitespace-pre-wrap">
          {advice}
        </p>
      </CardContent>
    </Card>
  );
}
