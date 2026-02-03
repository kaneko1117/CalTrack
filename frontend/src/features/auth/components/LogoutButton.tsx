/**
 * LogoutButton - ログアウトボタンコンポーネント
 */
import { Button } from "@/components/ui/button";
import { useRequestMutation } from "@/features/common/hooks";

type LogoutResponse = {
  message: string;
};

type LogoutButtonProps = {
  onSuccess?: () => void;
};

export function LogoutButton({ onSuccess }: LogoutButtonProps) {
  const { trigger, isMutating } = useRequestMutation<LogoutResponse>(
    "/api/v1/auth/logout",
    "POST",
    { onSuccess }
  );

  return (
    <Button onClick={() => trigger()} disabled={isMutating}>
      {isMutating ? "ログアウト中..." : "ログアウト"}
    </Button>
  );
}
