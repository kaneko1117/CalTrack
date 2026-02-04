/**
 * StatisticsCard コンポーネントのテスト
 */
import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { StatisticsCard } from "./StatisticsCard";
import type { StatisticsResponse } from "../hooks/useStatistics";

// useCountUpをモックして即座に最終値を返すようにする
vi.mock("@/features/common/hooks", () => ({
  useCountUp: ({ end }: { end: number }) => end,
}));

describe("StatisticsCard", () => {
  const mockData: StatisticsResponse = {
    period: "week",
    targetCalories: 2000,
    averageCalories: 1800,
    totalDays: 7,
    achievedDays: 5,
    overDays: 2,
    dailyStatistics: [],
  };

  describe("目標カロリー表示", () => {
    it("目標カロリーが表示される", () => {
      render(<StatisticsCard data={mockData} />);

      expect(screen.getByText("目標")).toBeInTheDocument();
      expect(screen.getByText("2,000")).toBeInTheDocument();
    });

    it("目標カロリーの単位が表示される", () => {
      render(<StatisticsCard data={mockData} />);

      expect(screen.getAllByText("kcal/日")).toHaveLength(2);
    });
  });

  describe("平均摂取カロリー表示", () => {
    it("平均摂取カロリーが表示される", () => {
      render(<StatisticsCard data={mockData} />);

      expect(screen.getByText("平均摂取")).toBeInTheDocument();
      expect(screen.getByText("1,800")).toBeInTheDocument();
    });

    it("平均摂取カロリーが青色グラデーションで表示される", () => {
      render(<StatisticsCard data={mockData} />);

      const avgElement = screen.getByText("1,800");
      expect(avgElement).toHaveClass("from-blue-600");
    });
  });

  describe("達成日数表示", () => {
    it("達成日数が表示される", () => {
      render(<StatisticsCard data={mockData} />);

      expect(screen.getByText("達成日数")).toBeInTheDocument();
      expect(screen.getByText("5")).toBeInTheDocument();
    });

    it("達成率が表示される", () => {
      render(<StatisticsCard data={mockData} />);

      // 5/7 = 71%
      expect(screen.getByText("日 (71%)")).toBeInTheDocument();
    });

    it("達成日数が緑色グラデーションで表示される", () => {
      render(<StatisticsCard data={mockData} />);

      const achievedElement = screen.getByText("5");
      expect(achievedElement).toHaveClass("from-emerald-500");
    });
  });

  describe("超過日数表示", () => {
    it("超過日数が表示される", () => {
      render(<StatisticsCard data={mockData} />);

      expect(screen.getByText("超過日数")).toBeInTheDocument();
      expect(screen.getByText("2")).toBeInTheDocument();
    });

    it("超過日数がある場合は赤色グラデーションで表示される", () => {
      render(<StatisticsCard data={mockData} />);

      const overElement = screen.getByText("2");
      expect(overElement).toHaveClass("from-red-500");
    });

    it("超過日数が0の場合はグレーグラデーションで表示される", () => {
      const noOverData: StatisticsResponse = {
        ...mockData,
        overDays: 0,
      };
      render(<StatisticsCard data={noOverData} />);

      const overElement = screen.getByText("0");
      expect(overElement).toHaveClass("from-gray-400");
    });
  });

  describe("達成率計算", () => {
    it("totalDaysが0の場合、達成率は0%と表示される", () => {
      const zeroDaysData: StatisticsResponse = {
        ...mockData,
        totalDays: 0,
        achievedDays: 0,
      };
      render(<StatisticsCard data={zeroDaysData} />);

      expect(screen.getByText("日 (0%)")).toBeInTheDocument();
    });

    it("全日達成の場合、100%と表示される", () => {
      const fullAchievedData: StatisticsResponse = {
        ...mockData,
        totalDays: 7,
        achievedDays: 7,
        overDays: 0,
      };
      render(<StatisticsCard data={fullAchievedData} />);

      expect(screen.getByText("日 (100%)")).toBeInTheDocument();
    });
  });

  describe("数値フォーマット", () => {
    it("大きな数値がカンマ区切りで表示される", () => {
      const largeData: StatisticsResponse = {
        ...mockData,
        targetCalories: 12345,
        averageCalories: 10000,
      };
      render(<StatisticsCard data={largeData} />);

      expect(screen.getByText("12,345")).toBeInTheDocument();
      expect(screen.getByText("10,000")).toBeInTheDocument();
    });
  });

  describe("レイアウト", () => {
    it("4つのカードが表示される", () => {
      const { container } = render(<StatisticsCard data={mockData} />);

      // Card コンポーネントの子要素（CardContent）を持つ要素が4つ
      const cards = container.querySelectorAll(".flex.flex-col.items-center");
      expect(cards).toHaveLength(4);
    });
  });
});
