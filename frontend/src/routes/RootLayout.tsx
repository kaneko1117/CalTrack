/**
 * RootLayout - アプリケーション全体のレイアウトコンポーネント
 *
 * 将来的には共通ヘッダー、フッター、ナビゲーションなどを追加可能
 */
import { Outlet } from "react-router-dom";

/**
 * RootLayout - ルートレイアウト
 *
 * 全ページに共通のレイアウトを提供
 * Outletで子ルートのコンテンツをレンダリング
 */
export function RootLayout() {
  return (
    <>
      {/* 将来追加予定: <Header /> */}
      <main>
        <Outlet />
      </main>
      {/* 将来追加予定: <Footer /> */}
    </>
  );
}
