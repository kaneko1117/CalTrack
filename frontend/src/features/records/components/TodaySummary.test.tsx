/**
 * TodaySummary コンポーネントのテスト
 */
import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { render, screen } from "@testing-library/react";
import { TodaySummary, TodaySummaryProps } from "./TodaySummary";
import type { TodayRecordsResponse } from "../hooks/useTodayRecords";
import type { ApiErrorResponse } from "@/lib/api";

// useCountUpをモックして即座に最終値を返すようにする
vi.mock("@/features/common/hooks", () => ({
  useCountUp: ({ end }: { end: number }) => end,
}));

describe("TodaySummary", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  // アニメーションをモックしているため、completeAnimationsは不要になったが
  // 既存の呼び出しのために空関数として残す
  const completeAnimations = () => {};
  // difference = targetCalories - totalCalories
  // 正の値 = 残りがある, 負の値 = 超過
  const mockData: TodayRecordsResponse = {
    date: "2024-01-15",
    totalCalories: 1200,
    targetCalories: 2000,
    difference: 800, // 2000 - 1200 = 800（残り）
    records: [],
  };

  const defaultProps: TodaySummaryProps = {
    data: null,
    isLoading: false,
    error: null,
  };

  describe("ローディング状態", () => {
    it("ローディング中はスケルトンが表示される", () => {
      render(<TodaySummary {...defaultProps} isLoading={true} data={null} />);

      // Skeletonコンポーネントが表示される（6つ - 3カード x 2スケルトン）
      const skeletons = document.querySelectorAll(".animate-pulse");
      expect(skeletons.length).toBeGreaterThanOrEqual(3);
    });

    it("データがある状態でローディング中の場合、スケルトンは表示されない", () => {
      render(<TodaySummary {...defaultProps} isLoading={true} data={mockData} />);
      completeAnimations();

      // データが表示される
      expect(screen.getByText("目標")).toBeInTheDocument();
      expect(screen.getByText("2,000")).toBeInTheDocument();
    });
  });

  describe("エラー状態", () => {
    it("エラー時にエラーメッセージが表示される", () => {
      const mockError: ApiErrorResponse = {
        code: "INTERNAL_ERROR",
        message: "Internal Server Error",
      };

      render(<TodaySummary {...defaultProps} error={mockError} />);

      expect(screen.getByText("データの取得に失敗しました")).toBeInTheDocument();
    });
  });

  describe("データなし状態", () => {
    it("データがない場合は何も表示されない", () => {
      const { container } = render(<TodaySummary {...defaultProps} />);

      expect(container.firstChild).toBeNull();
    });
  });

  describe("正常なデータ表示", () => {
    it("目標カロリーが表示される", () => {
      render(<TodaySummary {...defaultProps} data={mockData} />);
      completeAnimations();

      expect(screen.getByText("目標")).toBeInTheDocument();
      expect(screen.getByText("2,000")).toBeInTheDocument();
    });

    it("摂取カロリーが表示される", () => {
      render(<TodaySummary {...defaultProps} data={mockData} />);
      completeAnimations();

      expect(screen.getByText("摂取")).toBeInTheDocument();
      expect(screen.getByText("1,200")).toBeInTheDocument();
    });

    it("残りカロリーが表示される（目標未達の場合）", () => {
      render(<TodaySummary {...defaultProps} data={mockData} />);
      completeAnimations();

      expect(screen.getByText("残り")).toBeInTheDocument();
      expect(screen.getByText("800")).toBeInTheDocument();
    });

    it("kcal単位が各カードに表示される", () => {
      render(<TodaySummary {...defaultProps} data={mockData} />);
      completeAnimations();

      const kcalLabels = screen.getAllByText("kcal");
      expect(kcalLabels).toHaveLength(3);
    });
  });

  describe("目標超過の表示", () => {
    it("目標超過時は「超過」と表示される", () => {
      const overTargetData: TodayRecordsResponse = {
        ...mockData,
        totalCalories: 2500,
        difference: -500, // 2000 - 2500 = -500（超過）
      };

      render(<TodaySummary {...defaultProps} data={overTargetData} />);
      completeAnimations();

      expect(screen.getByText("超過")).toBeInTheDocument();
      expect(screen.getByText("500")).toBeInTheDocument();
    });

    it("目標超過時は超過カロリーが赤色で表示される", () => {
      const overTargetData: TodayRecordsResponse = {
        ...mockData,
        totalCalories: 2500,
        difference: -500, // 2000 - 2500 = -500（超過）
      };

      render(<TodaySummary {...defaultProps} data={overTargetData} />);
      completeAnimations();

      // gradient クラスを持つ要素を探す（from-red-500）
      const overElement = screen.getByText("500").closest("p");
      expect(overElement).toHaveClass("from-red-500");
    });

    it("目標未達時は残りカロリーが緑色で表示される", () => {
      render(<TodaySummary {...defaultProps} data={mockData} />);
      completeAnimations();

      // gradient クラスを持つ要素を探す（from-emerald-500）
      const remainingElement = screen.getByText("800").closest("p");
      expect(remainingElement).toHaveClass("from-emerald-500");
    });
  });

  describe("数値フォーマット", () => {
    it("大きな数値がカンマ区切りで表示される", () => {
      const largeData: TodayRecordsResponse = {
        ...mockData,
        totalCalories: 12345,
        targetCalories: 25000,
      };

      render(<TodaySummary {...defaultProps} data={largeData} />);
      completeAnimations();

      expect(screen.getByText("12,345")).toBeInTheDocument();
      expect(screen.getByText("25,000")).toBeInTheDocument();
    });
  });

  describe("境界値テスト", () => {
    it("カロリーが0の場合でも表示される", () => {
      const zeroData: TodayRecordsResponse = {
        ...mockData,
        totalCalories: 0,
        targetCalories: 2000,
        difference: 2000,
      };

      render(<TodaySummary {...defaultProps} data={zeroData} />);
      completeAnimations();

      // 摂取カロリーが0（ProgressRingの達成率0%と、摂取カロリー0の両方で表示される）
      expect(screen.getAllByText("0")).toHaveLength(2);
      // 目標カロリーと残りカロリーが同じ2000なのでgetAllByTextを使用
      expect(screen.getAllByText("2,000")).toHaveLength(2);
    });

    it("目標と摂取が同じ場合、残りは0と表示される", () => {
      const equalData: TodayRecordsResponse = {
        ...mockData,
        totalCalories: 2000,
        targetCalories: 2000,
        difference: 0,
      };

      render(<TodaySummary {...defaultProps} data={equalData} />);
      completeAnimations();

      expect(screen.getByText("残り")).toBeInTheDocument();
      // 目標と摂取が2000、残りが0
      expect(screen.getAllByText("2,000")).toHaveLength(2);
      expect(screen.getByText("0")).toBeInTheDocument();
    });
  });
});
