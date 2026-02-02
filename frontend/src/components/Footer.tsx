/**
 * Footer - 共通フッターコンポーネント
 */
export function Footer() {
  return (
    <footer className="border-t bg-muted/30">
      <div className="container py-4 px-4 mx-auto max-w-4xl">
        <p className="text-sm text-center text-muted-foreground">
          &copy; {new Date().getFullYear()} CalTrack
        </p>
      </div>
    </footer>
  );
}
