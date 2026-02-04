/// <reference types="vite/client" />

/**
 * Vite環境変数の型定義
 */
interface ImportMetaEnv {
  /** API URL（バックエンドのベースURL） */
  readonly VITE_API_URL: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
