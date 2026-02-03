type ProgressRingProps = {
  progress: number; // 0-100
  size?: number;
  strokeWidth?: number;
};

/**
 * 円形プログレスリングコンポーネント
 */
export function ProgressRing({
  progress,
  size = 120,
  strokeWidth = 8,
}: ProgressRingProps) {
  const radius = (size - strokeWidth) / 2;
  const circumference = radius * 2 * Math.PI;
  const offset = circumference - (progress / 100) * circumference;

  return (
    <div className="relative inline-flex items-center justify-center">
      <svg width={size} height={size} className="-rotate-90">
        {/* 背景円 */}
        <circle
          cx={size / 2}
          cy={size / 2}
          r={radius}
          fill="none"
          stroke="currentColor"
          strokeWidth={strokeWidth}
          className="text-muted/20"
        />
        {/* プログレス円 */}
        <circle
          cx={size / 2}
          cy={size / 2}
          r={radius}
          fill="none"
          stroke="currentColor"
          strokeWidth={strokeWidth}
          strokeLinecap="round"
          strokeDasharray={circumference}
          strokeDashoffset={offset}
          className="text-primary transition-all duration-1000 ease-out"
        />
      </svg>
      {/* 中央のパーセンテージ表示 */}
      <div className="absolute inset-0 flex flex-col items-center justify-center">
        <span className="text-3xl font-extrabold tabular-nums tracking-tight bg-gradient-to-br from-primary to-primary/70 bg-clip-text text-transparent drop-shadow-sm">
          {Math.round(progress)}
        </span>
        <span className="text-xs font-medium text-muted-foreground -mt-0.5">
          %
        </span>
      </div>
    </div>
  );
}
