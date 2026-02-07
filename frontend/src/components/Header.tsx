/**
 * Header - 共通ヘッダーコンポーネント
 */
import { Link, useNavigate } from "react-router-dom";
import { useRequestMutation } from "@/features/common/hooks";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

// ハンバーガーメニューアイコン
function MenuIcon() {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="24"
      height="24"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
    >
      <line x1="3" y1="12" x2="21" y2="12" />
      <line x1="3" y1="6" x2="21" y2="6" />
      <line x1="3" y1="18" x2="21" y2="18" />
    </svg>
  );
}

// ログアウトアイコン
function LogoutIcon() {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="16"
      height="16"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      className="mr-2"
    >
      <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
      <polyline points="16 17 21 12 16 7" />
      <line x1="21" y1="12" x2="9" y2="12" />
    </svg>
  );
}

// 設定アイコン
function SettingsIcon() {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="16"
      height="16"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      className="mr-2"
    >
      <circle cx="12" cy="12" r="3" />
      <path d="M12 1v6m0 6v6m0-6h6m-6 0H6m12.22-5.22l-4.24 4.24m0 0l-4.24 4.24m4.24-4.24l-4.24-4.24m0 0L9.78 18.22" />
    </svg>
  );
}

type LogoutResponse = { message: string };

export function Header() {
  const navigate = useNavigate();
  const { trigger: logout, isMutating: isLoggingOut } =
    useRequestMutation<LogoutResponse>("/api/v1/auth/logout", "POST", {
      onSuccess: () => navigate("/"),
    });

  const handleSettingsClick = () => {
    navigate("/settings");
  };

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur">
      <div className="container flex h-14 items-center justify-between px-4 mx-auto max-w-4xl">
        <Link to="/dashboard" className="flex items-center gap-2">
          <span className="text-xl font-bold text-primary">CalTrack</span>
        </Link>

        {/* ハンバーガーメニュー */}
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <button
              type="button"
              className="rounded-md p-2 hover:bg-secondary transition-colors"
              aria-label="メニューを開く"
            >
              <MenuIcon />
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            {/* 設定 */}
            <DropdownMenuItem onClick={handleSettingsClick}>
              <SettingsIcon />
              設定
            </DropdownMenuItem>

            <DropdownMenuSeparator />

            {/* ログアウト */}
            <DropdownMenuItem
              onClick={() => logout()}
              disabled={isLoggingOut}
            >
              <LogoutIcon />
              {isLoggingOut ? "ログアウト中..." : "ログアウト"}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </header>
  );
}
