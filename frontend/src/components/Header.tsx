/**
 * Header - 共通ヘッダーコンポーネント
 */
import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { LogoutButton } from "@/features/auth";
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetTrigger } from "@/components/ui/sheet";
import { Button } from "@/components/ui/button";

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

// 設定アイコン
function SettingsIcon() {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width="20"
      height="20"
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

export function Header() {
  const navigate = useNavigate();
  const [menuOpen, setMenuOpen] = useState(false);

  const handleSettingsClick = () => {
    navigate("/settings");
    setMenuOpen(false);
  };

  const handleLogoutSuccess = () => {
    setMenuOpen(false);
    navigate("/");
  };

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur">
      <div className="container flex h-14 items-center justify-between px-4 mx-auto max-w-4xl">
        <Link to="/dashboard" className="flex items-center gap-2">
          <span className="text-xl font-bold text-primary">CalTrack</span>
        </Link>

        {/* ハンバーガーメニュー */}
        <Sheet open={menuOpen} onOpenChange={setMenuOpen}>
          <SheetTrigger asChild>
            <button
              type="button"
              className="rounded-md p-2 hover:bg-secondary transition-colors"
              aria-label="メニューを開く"
            >
              <MenuIcon />
            </button>
          </SheetTrigger>
          <SheetContent>
            <SheetHeader>
              <SheetTitle>メニュー</SheetTitle>
            </SheetHeader>

            <div className="flex flex-col gap-4 mt-8">
              {/* 設定リンク */}
              <Button
                variant="outline"
                onClick={handleSettingsClick}
                className="w-full justify-start"
              >
                <SettingsIcon />
                設定
              </Button>

              {/* ログアウトボタン */}
              <div className="mt-auto pt-8">
                <LogoutButton onSuccess={handleLogoutSuccess} />
              </div>
            </div>
          </SheetContent>
        </Sheet>
      </div>
    </header>
  );
}
