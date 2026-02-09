/**
 * PfcProgressCard テスト
 */
import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { PfcProgressCard } from "./PfcProgressCard";
import type { TodayPfcResponse } from "../hooks/useTodayPfc";

vi.mock("@/features/common/hooks", () => ({
  useCountUp: ({ end }: { end: number }) => end,
}));

const mockData: TodayPfcResponse = {
  date: "2026-02-09T00:00:00Z",
  current: { protein: 35.0, fat: 25.0, carbs: 100.0 },
  target: { protein: 120.0, fat: 65.0, carbs: 300.0 },
};

describe("PfcProgressCard", () => {
  it("PFCデータが正しく表示される", () => {
    render(<PfcProgressCard data={mockData} isLoading={false} error={null} />);
    expect(screen.getByText("今日のPFCバランス")).toBeInTheDocument();
    expect(screen.getByText("タンパク質")).toBeInTheDocument();
    expect(screen.getByText("脂質")).toBeInTheDocument();
    expect(screen.getByText("炭水化物")).toBeInTheDocument();
  });

  it("タンパク質の数値とパーセンテージが正しく表示される", () => {
    render(<PfcProgressCard data={mockData} isLoading={false} error={null} />);
    expect(screen.getByText("35g / 120g (29%)")).toBeInTheDocument();
  });

  it("脂質の数値とパーセンテージが正しく表示される", () => {
    render(<PfcProgressCard data={mockData} isLoading={false} error={null} />);
    expect(screen.getByText("25g / 65g (38%)")).toBeInTheDocument();
  });

  it("炭水化物の数値とパーセンテージが正しく表示される", () => {
    render(<PfcProgressCard data={mockData} isLoading={false} error={null} />);
    expect(screen.getByText("100g / 300g (33%)")).toBeInTheDocument();
  });

  it("進捗バーにaria-labelが設定される", () => {
    render(<PfcProgressCard data={mockData} isLoading={false} error={null} />);
    expect(screen.getByLabelText("タンパク質の進捗")).toBeInTheDocument();
    expect(screen.getByLabelText("脂質の進捗")).toBeInTheDocument();
    expect(screen.getByLabelText("炭水化物の進捗")).toBeInTheDocument();
  });

  it("ローディング中はデータが表示されない", () => {
    render(<PfcProgressCard data={mockData} isLoading={true} error={null} />);
    expect(screen.queryByText("タンパク質")).not.toBeInTheDocument();
  });

  it("エラー時はエラーメッセージが表示される", () => {
    render(
      <PfcProgressCard
        data={null}
        isLoading={false}
        error={{ code: "INTERNAL_ERROR", message: "Server error" }}
      />
    );
    expect(screen.getByText("PFCデータの取得に失敗しました")).toBeInTheDocument();
  });

  it("データがない場合は案内メッセージが表示される", () => {
    render(<PfcProgressCard data={null} isLoading={false} error={null} />);
    expect(screen.getByText("食事を記録するとPFCバランスが表示されます")).toBeInTheDocument();
  });

  it("100%超は実際のパーセンテージが表示される", () => {
    const overData: TodayPfcResponse = {
      date: "2026-02-09T00:00:00Z",
      current: { protein: 200.0, fat: 100.0, carbs: 500.0 },
      target: { protein: 120.0, fat: 65.0, carbs: 300.0 },
    };
    render(<PfcProgressCard data={overData} isLoading={false} error={null} />);
    expect(screen.getByText("200g / 120g (167%)")).toBeInTheDocument();
    expect(screen.getByText("100g / 65g (154%)")).toBeInTheDocument();
    expect(screen.getByText("500g / 300g (167%)")).toBeInTheDocument();
  });

  it("80-100%でoptimalステータスのスタイルが適用される", () => {
    const optimalData: TodayPfcResponse = {
      date: "2026-02-09T00:00:00Z",
      current: { protein: 100.0, fat: 55.0, carbs: 250.0 },
      target: { protein: 120.0, fat: 65.0, carbs: 300.0 },
    };
    render(<PfcProgressCard data={optimalData} isLoading={false} error={null} />);
    // protein: 83%, fat: 85%, carbs: 83% → 全部optimal
    expect(screen.getByText("100g / 120g (83%)")).toBeInTheDocument();
    expect(screen.getByText("55g / 65g (85%)")).toBeInTheDocument();
    expect(screen.getByText("250g / 300g (83%)")).toBeInTheDocument();

    // optimalステータスでtext-green-600が適用されること
    const proteinText = screen.getByText("100g / 120g (83%)");
    expect(proteinText.className).toContain("text-green-600");
  });

  it("100%超でoverステータスのスタイルが適用される", () => {
    const overData: TodayPfcResponse = {
      date: "2026-02-09T00:00:00Z",
      current: { protein: 200.0, fat: 100.0, carbs: 500.0 },
      target: { protein: 120.0, fat: 65.0, carbs: 300.0 },
    };
    render(<PfcProgressCard data={overData} isLoading={false} error={null} />);

    // overステータスでtext-red-600が適用されること
    const proteinText = screen.getByText("200g / 120g (167%)");
    expect(proteinText.className).toContain("text-red-600");
  });
});
