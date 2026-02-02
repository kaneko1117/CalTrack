/**
 * LoginPage - ログインページコンポーネント
 * ルーティング対応のためのページラッパー
 * Warm & Organicトーンのデザイン
 */
import { useNavigate } from "react-router-dom";
import { LoginForm } from "../features/auth/components/LoginForm";
import type { LoginResponse } from "../features/auth/components/LoginForm";

/** LoginPageコンポーネントのProps */
export type LoginPageProps = {
  /** ログイン成功時の遷移先URL */
  redirectTo?: string;
};

/**
 * LoginPage - ログインページ
 * 背景装飾とロゴエリアを含むフルページレイアウト
 */
export function LoginPage({ redirectTo = "/dashboard" }: LoginPageProps) {
  const navigate = useNavigate();

  /**
   * ログイン成功時のハンドラ
   * 指定されたパスへ遷移する
   */
  const handleSuccess = (_response: LoginResponse) => {
    navigate(redirectTo);
  };

  return (
    <div className="min-h-screen w-full flex flex-col items-center justify-center bg-background py-12 px-4 sm:px-6 lg:px-8 relative overflow-hidden">
      {/* 背景装飾 - 有機的な円形グラデーション */}
      <div
        className="absolute top-0 right-0 w-96 h-96 rounded-full opacity-30 blur-3xl -translate-y-1/2 translate-x-1/2"
        style={{
          background:
            "radial-gradient(circle, hsl(142 40% 45% / 0.4) 0%, transparent 70%)",
        }}
        aria-hidden="true"
      />
      <div
        className="absolute bottom-0 left-0 w-80 h-80 rounded-full opacity-20 blur-3xl translate-y-1/2 -translate-x-1/2"
        style={{
          background:
            "radial-gradient(circle, hsl(25 80% 55% / 0.3) 0%, transparent 70%)",
        }}
        aria-hidden="true"
      />

      {/* コンテンツ幅制限用コンテナ - PC表示時の幅を制限 */}
      <div className="w-full max-w-screen-sm mx-auto flex flex-col items-center relative z-10">
        {/* ロゴエリア */}
        <div className="mb-8 text-center">
          <h1 className="text-3xl font-bold text-primary tracking-tight">
            CalTrack
          </h1>
          <p className="mt-2 text-muted-foreground">
            あなたの健康的な食生活をサポート
          </p>
        </div>

        {/* ログインフォーム */}
        <div className="w-full max-w-md">
          <LoginForm onSuccess={handleSuccess} />
        </div>
      </div>
    </div>
  );
}
