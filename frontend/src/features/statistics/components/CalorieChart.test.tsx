/**
 * CalorieChart コンポーネントのテスト
 */
import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { CalorieChart } from "./CalorieChart";
import type { DailyStatistic } from "../hooks/useStatistics";

// recharts のモック（ResizeObserverエラー回避）
vi.mock("recharts", async (importOriginal) => {
  const actual = await importOriginal<typeof import("recharts")>();
  return {
    ...actual,
    ResponsiveContainer: ({ children }: { children: React.ReactNode }) => (
      <div data-testid="responsive-container">{children}</div>
    ),
  };
});

describe("CalorieChart", () => {
  const mockData: DailyStatistic[] = [
    { date: "2024-01-15", totalCalories: 1800 },
    { date: "2024-01-16", totalCalories: 2100 },
    { date: "2024-01-17", totalCalories: 1900 },
  ];

  describe("ローディング状態", () => {
    it("ローディング中はスケルトンが表示される", () => {
      const { container } = render(
        <CalorieChart data={[]} targetCalories={2000} isLoading={true} />
      );

      // Skeletonコンポーネントが表示される
      const skeleton = container.querySelector(".animate-pulse");
      expect(skeleton).toBeInTheDocument();
    });

    it("ローディング中もタイトルは表示される", () => {
      render(<CalorieChart data={[]} targetCalories={2000} isLoading={true} />);

      expect(screen.getByText("カロリー推移")).toBeInTheDocument();
    });
  });

  describe("データなし状態", () => {
    it("データが空の場合、メッセージが表示される", () => {
      render(<CalorieChart data={[]} targetCalories={2000} isLoading={false} />);

      expect(screen.getByText("表示するデータがありません")).toBeInTheDocument();
    });

    it("データが空でもタイトルは表示される", () => {
      render(<CalorieChart data={[]} targetCalories={2000} isLoading={false} />);

      expect(screen.getByText("カロリー推移")).toBeInTheDocument();
    });
  });

  describe("データ表示", () => {
    it("タイトルが表示される", () => {
      render(<CalorieChart data={mockData} targetCalories={2000} />);

      expect(screen.getByText("カロリー推移")).toBeInTheDocument();
    });

    it("凡例が表示される", () => {
      render(<CalorieChart data={mockData} targetCalories={2000} />);

      expect(screen.getByText("摂取カロリー")).toBeInTheDocument();
      expect(screen.getByText("目標")).toBeInTheDocument();
    });
  });

  describe("デフォルト値", () => {
    it("isLoadingのデフォルトはfalse", () => {
      render(<CalorieChart data={mockData} targetCalories={2000} />);

      // スケルトンが表示されない
      const skeleton = document.querySelector(".animate-pulse");
      expect(skeleton).not.toBeInTheDocument();
    });
  });

  describe("境界値", () => {
    it("データが1件でも表示される", () => {
      const singleData: DailyStatistic[] = [
        { date: "2024-01-15", totalCalories: 1800 },
      ];
      render(<CalorieChart data={singleData} targetCalories={2000} />);

      expect(screen.getByText("カロリー推移")).toBeInTheDocument();
      expect(screen.getByText("摂取カロリー")).toBeInTheDocument();
    });

    it("targetCaloriesが0でも表示される", () => {
      render(<CalorieChart data={mockData} targetCalories={0} />);

      expect(screen.getByText("カロリー推移")).toBeInTheDocument();
    });

    it("大量のデータでもエラーにならない", () => {
      const largeData: DailyStatistic[] = Array.from({ length: 31 }, (_, i) => ({
        date: `2024-01-${String(i + 1).padStart(2, "0")}`,
        totalCalories: 1800 + (i % 5) * 100,
      }));
      render(<CalorieChart data={largeData} targetCalories={2000} />);

      expect(screen.getByText("カロリー推移")).toBeInTheDocument();
    });
  });
});
