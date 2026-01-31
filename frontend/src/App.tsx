/**
 * App - アプリケーションのルートコンポーネント
 *
 * React Routerを使用してルーティングを管理
 */
import { RouterProvider } from "react-router-dom";
import { router } from "./routes";

/**
 * App - メインアプリケーションコンポーネント
 *
 * RouterProviderでルーティング機能を提供
 */
function App() {
  return <RouterProvider router={router} />;
}

export default App;
