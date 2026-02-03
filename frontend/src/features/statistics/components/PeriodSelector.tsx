/**
 * PeriodSelector - 期間選択タブ
 */
import { type Period, PERIOD_LABELS } from "../hooks/useStatistics";

export type PeriodSelectorProps = {
  value: Period;
  onChange: (period: Period) => void;
};

const PERIODS: Period[] = ["week", "month"];

export function PeriodSelector({ value, onChange }: PeriodSelectorProps) {
  return (
    <div className="flex gap-1 p-1 bg-muted rounded-lg">
      {PERIODS.map((period) => (
        <button
          key={period}
          type="button"
          onClick={() => onChange(period)}
          className={`px-4 py-2 text-sm font-medium rounded-md transition-colors ${
            value === period
              ? "bg-background text-foreground shadow-sm"
              : "text-muted-foreground hover:text-foreground"
          }`}
        >
          {PERIOD_LABELS[period]}
        </button>
      ))}
    </div>
  );
}
