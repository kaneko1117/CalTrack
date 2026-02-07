/**
 * App - アプリケーションのルートコンポーネント
 *
 * React Routerを使用してルーティングを管理
 */
import { RouterProvider } from "react-router-dom";
import { SWRConfig } from "swr";
import { router } from "./routes";
import { fetcher } from "@/lib/swr";

/**
 * App - メインアプリケーションコンポーネント
 *
 * RouterProviderでルーティング機能を提供
 */
function App() {
  return (
    <SWRConfig value={{ fetcher }}>
      <RouterProvider router={router} />
    </SWRConfig>
  );
}

export default App;
