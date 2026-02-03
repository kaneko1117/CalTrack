/**
 * LogoutButton - ログアウトボタンコンポーネント
 */
import { Button } from "@/components/ui/button";
import { useApi } from "@/features/common/hooks";
import { post } from "@/lib/api";

type LogoutResponse = {
  message: string;
};

type LogoutButtonProps = {
  onSuccess?: () => void;
};

const logout = () => post<LogoutResponse>("/api/v1/auth/logout");

export function LogoutButton({ onSuccess }: LogoutButtonProps) {
  const { execute, isPending } = useApi(logout, { onSuccess });

  return (
    <Button onClick={execute} disabled={isPending}>
      {isPending ? "ログアウト中..." : "ログアウト"}
    </Button>
  );
}
