/**
 * RegisterPage - 新規登録ページコンポーネント
 * ルーティング対応のためのページラッパー
 */
import { useNavigate } from "react-router-dom";
import { RegisterForm } from "./RegisterForm";

/** RegisterPageコンポーネントのProps */
export type RegisterPageProps = {
  /** 登録成功時の遷移先URL */
  redirectTo?: string;
};

/**
 * RegisterPage - 新規登録ページ
 * フルページレイアウトでRegisterFormを表示
 */
export function RegisterPage({ redirectTo = "/" }: RegisterPageProps) {
  const navigate = useNavigate();

  /**
   * 登録成功時のハンドラ
   * 指定されたパスへ遷移する
   */
  const handleSuccess = () => {
    navigate(redirectTo);
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <RegisterForm onSuccess={handleSuccess} />
    </div>
  );
}
