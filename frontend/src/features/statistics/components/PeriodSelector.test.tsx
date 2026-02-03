/**
 * PeriodSelector コンポーネントのテスト
 */
import { describe, it, expect, vi } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { PeriodSelector } from "./PeriodSelector";

describe("PeriodSelector", () => {
  describe("レンダリング", () => {
    it("週間と月間のボタンが表示される", () => {
      const onChange = vi.fn();
      render(<PeriodSelector value="week" onChange={onChange} />);

      expect(screen.getByText("週間")).toBeInTheDocument();
      expect(screen.getByText("月間")).toBeInTheDocument();
    });

    it("ボタンが2つ表示される", () => {
      const onChange = vi.fn();
      render(<PeriodSelector value="week" onChange={onChange} />);

      const buttons = screen.getAllByRole("button");
      expect(buttons).toHaveLength(2);
    });
  });

  describe("選択状態", () => {
    it("weekが選択されている場合、週間ボタンがアクティブスタイルを持つ", () => {
      const onChange = vi.fn();
      render(<PeriodSelector value="week" onChange={onChange} />);

      const weekButton = screen.getByText("週間");
      expect(weekButton).toHaveClass("bg-background");
    });

    it("monthが選択されている場合、月間ボタンがアクティブスタイルを持つ", () => {
      const onChange = vi.fn();
      render(<PeriodSelector value="month" onChange={onChange} />);

      const monthButton = screen.getByText("月間");
      expect(monthButton).toHaveClass("bg-background");
    });

    it("非選択のボタンはアクティブスタイルを持たない", () => {
      const onChange = vi.fn();
      render(<PeriodSelector value="week" onChange={onChange} />);

      const monthButton = screen.getByText("月間");
      expect(monthButton).not.toHaveClass("bg-background");
      expect(monthButton).toHaveClass("text-muted-foreground");
    });
  });

  describe("クリックイベント", () => {
    it("週間ボタンをクリックするとonChangeがweekで呼ばれる", () => {
      const onChange = vi.fn();
      render(<PeriodSelector value="month" onChange={onChange} />);

      fireEvent.click(screen.getByText("週間"));
      expect(onChange).toHaveBeenCalledWith("week");
    });

    it("月間ボタンをクリックするとonChangeがmonthで呼ばれる", () => {
      const onChange = vi.fn();
      render(<PeriodSelector value="week" onChange={onChange} />);

      fireEvent.click(screen.getByText("月間"));
      expect(onChange).toHaveBeenCalledWith("month");
    });

    it("現在選択中のボタンをクリックしてもonChangeが呼ばれる", () => {
      const onChange = vi.fn();
      render(<PeriodSelector value="week" onChange={onChange} />);

      fireEvent.click(screen.getByText("週間"));
      expect(onChange).toHaveBeenCalledWith("week");
    });
  });

  describe("アクセシビリティ", () => {
    it("ボタンにtype=buttonが設定されている", () => {
      const onChange = vi.fn();
      render(<PeriodSelector value="week" onChange={onChange} />);

      const buttons = screen.getAllByRole("button");
      buttons.forEach((button) => {
        expect(button).toHaveAttribute("type", "button");
      });
    });
  });
});
