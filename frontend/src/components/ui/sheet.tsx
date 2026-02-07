/**
 * Sheet - 軽量ドロワーコンポーネント
 * Radix UIなしのピュアReact実装
 */
import { createContext, useContext, useState, useCallback, ReactNode } from "react";
import { cn } from "@/lib/utils";

// Context型定義
type SheetContextValue = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
};

const SheetContext = createContext<SheetContextValue | undefined>(undefined);

// Contextフック
function useSheet() {
  const context = useContext(SheetContext);
  if (!context) {
    throw new Error("Sheet components must be used within Sheet");
  }
  return context;
}

// ルートコンポーネント
type SheetProps = {
  children: ReactNode;
  open?: boolean;
  onOpenChange?: (open: boolean) => void;
};

export function Sheet({ children, open: controlledOpen, onOpenChange }: SheetProps) {
  const [uncontrolledOpen, setUncontrolledOpen] = useState(false);

  const open = controlledOpen !== undefined ? controlledOpen : uncontrolledOpen;
  const handleOpenChange = useCallback(
    (newOpen: boolean) => {
      if (onOpenChange) {
        onOpenChange(newOpen);
      } else {
        setUncontrolledOpen(newOpen);
      }
    },
    [onOpenChange]
  );

  return (
    <SheetContext.Provider value={{ open, onOpenChange: handleOpenChange }}>
      {children}
    </SheetContext.Provider>
  );
}

// トリガーボタン
type SheetTriggerProps = {
  children: ReactNode;
  asChild?: boolean;
  className?: string;
  onClick?: () => void;
};

export function SheetTrigger({ children, asChild, className, onClick }: SheetTriggerProps) {
  const { onOpenChange } = useSheet();

  const handleClick = () => {
    onOpenChange(true);
    onClick?.();
  };

  if (asChild) {
    // asChildの場合、childrenの最初の要素のonClickを上書き
    // React.cloneElementを使用して型安全に実装
    const child = children as React.ReactElement<{ onClick?: (e: React.MouseEvent) => void }>;

    if (child && typeof child === "object" && "props" in child) {
      const originalOnClick = child.props.onClick;

      return (
        <>
          {typeof child.type === "string" || typeof child.type === "function"
            ? {
                ...child,
                props: {
                  ...(child.props as Record<string, unknown>),
                  onClick: (e: React.MouseEvent) => {
                    originalOnClick?.(e);
                    handleClick();
                  },
                },
              }
            : children}
        </>
      );
    }

    return <>{children}</>;
  }

  return (
    <button type="button" onClick={handleClick} className={className}>
      {children}
    </button>
  );
}

// コンテンツ
type SheetContentProps = {
  children: ReactNode;
  side?: "left" | "right";
  className?: string;
};

export function SheetContent({ children, side = "right", className }: SheetContentProps) {
  const { open, onOpenChange } = useSheet();

  if (!open) return null;

  const sideClass = side === "right" ? "right-0 animate-slide-in-right" : "left-0 animate-slide-in-left";

  return (
    <>
      {/* オーバーレイ */}
      <div
        className="fixed inset-0 z-50 bg-black/50 animate-fade-in"
        onClick={() => onOpenChange(false)}
        aria-hidden="true"
      />

      {/* コンテンツ */}
      <div
        className={cn(
          "fixed top-0 z-50 h-full w-80 bg-background p-6 shadow-lg",
          sideClass,
          className
        )}
        role="dialog"
        aria-modal="true"
      >
        {/* 閉じるボタン */}
        <button
          type="button"
          onClick={() => onOpenChange(false)}
          className="absolute right-4 top-4 rounded-sm opacity-70 ring-offset-background transition-opacity hover:opacity-100 focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2"
          aria-label="Close"
        >
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
            <line x1="18" y1="6" x2="6" y2="18" />
            <line x1="6" y1="6" x2="18" y2="18" />
          </svg>
        </button>

        {children}
      </div>
    </>
  );
}

// ヘッダー
type SheetHeaderProps = {
  children: ReactNode;
  className?: string;
};

export function SheetHeader({ children, className }: SheetHeaderProps) {
  return (
    <div className={cn("flex flex-col space-y-2 text-center sm:text-left", className)}>
      {children}
    </div>
  );
}

// タイトル
type SheetTitleProps = {
  children: ReactNode;
  className?: string;
};

export function SheetTitle({ children, className }: SheetTitleProps) {
  return (
    <h2 className={cn("text-lg font-semibold text-foreground", className)}>
      {children}
    </h2>
  );
}
