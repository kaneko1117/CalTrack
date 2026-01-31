/**
 * ルーティング設定
 * アプリケーション全体のルート定義を管理
 *
 * 将来の拡張に備えて、機能ごとにルートを分割可能な構成
 */
import { createBrowserRouter, RouteObject } from "react-router-dom";
import { RootLayout } from "./RootLayout";
import { HomePage } from "../pages/HomePage";
import { RegisterPage } from "../features/auth/components/RegisterPage";
import { LoginPage } from "../features/auth/components/LoginPage";

/**
 * 認証関連のルート定義
 */
const authRoutes: RouteObject[] = [
  {
    path: "/register",
    element: <RegisterPage redirectTo="/" />,
  },
  // 将来追加予定:
  // { path: "/forgot-password", element: <ForgotPasswordPage /> },
];

/**
 * アプリケーションのメインルート定義
 */
const mainRoutes: RouteObject[] = [
  {
    path: "/",
    element: <LoginPage redirectTo="/" />,
  },
  {
    path: "/home",
    element: <HomePage />,
  },
  // 将来追加予定:
  // { path: "/dashboard", element: <DashboardPage /> },
  // { path: "/meals", element: <MealsPage /> },
];

/**
 * 全ルートをRootLayoutでラップ
 */
export const routes: RouteObject[] = [
  {
    element: <RootLayout />,
    children: [...mainRoutes, ...authRoutes],
  },
];

/**
 * アプリケーションのルーター
 */
export const router = createBrowserRouter(routes);
