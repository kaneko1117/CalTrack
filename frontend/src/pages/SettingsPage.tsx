/**
 * SettingsPage - 設定ページ
 * ユーザープロフィールの編集が可能
 */
import { useState, useEffect } from "react";
import { Header } from "@/components/Header";
import { Footer } from "@/components/Footer";
import { ProfileEditForm } from "@/features/user";
import type { UpdateProfileResponse } from "@/features/user";

/**
 * CheckCircleアイコン - 成功表示用
 * SVGインラインアイコン
 */
function CheckCircleIcon({ className }: { className?: string }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      aria-hidden="true"
    >
      <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" />
      <polyline points="22 4 12 14.01 9 11.01" />
    </svg>
  );
}

/**
 * SettingsPage - 設定ページコンポーネント
 */
export function SettingsPage() {
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  /**
   * プロフィール更新成功時のコールバック
   * 成功メッセージを3秒間表示
   */
  const handleUpdateSuccess = (_result: UpdateProfileResponse) => {
    setSuccessMessage("プロフィールを更新しました");
    setTimeout(() => {
      setSuccessMessage(null);
    }, 3000);
  };

  // 成功メッセージが変化したら、3秒後に消す
  useEffect(() => {
    if (successMessage) {
      const timer = setTimeout(() => {
        setSuccessMessage(null);
      }, 3000);
      return () => clearTimeout(timer);
    }
  }, [successMessage]);

  return (
    <div className="min-h-screen flex flex-col bg-background">
      <Header />
      <main className="flex-1 container px-4 py-8 mx-auto max-w-3xl">
        <div className="space-y-6">
          {/* ページタイトル */}
          <h1 className="text-2xl font-bold">設定</h1>

          {/* 成功メッセージ */}
          {successMessage && (
            <div
              className="flex items-start gap-3 p-4 text-sm rounded-lg bg-green-50 border border-green-200"
              role="alert"
            >
              <CheckCircleIcon className="w-5 h-5 text-green-600 flex-shrink-0 mt-0.5" />
              <p className="font-medium text-green-800">{successMessage}</p>
            </div>
          )}

          {/* プロフィール編集フォーム */}
          <ProfileEditForm onSuccess={handleUpdateSuccess} />
        </div>
      </main>
      <Footer />
    </div>
  );
}
