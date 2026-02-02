/**
 * Header - 共通ヘッダーコンポーネント
 */
import { Link } from "react-router-dom";

export function Header() {
  return (
    <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur">
      <div className="container flex h-14 items-center justify-between px-4 mx-auto max-w-4xl">
        <Link to="/dashboard" className="flex items-center gap-2">
          <span className="text-xl font-bold text-primary">CalTrack</span>
        </Link>
      </div>
    </header>
  );
}
