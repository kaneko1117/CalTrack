import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { ProgressRing } from "./ProgressRing";

describe("ProgressRing", () => {
  describe("パーセンテージ表示", () => {
    it("パーセンテージが正しく表示されること", () => {
      render(<ProgressRing progress={50} />);
      expect(screen.getByText("50")).toBeInTheDocument();
      expect(screen.getByText("%")).toBeInTheDocument();
    });

    it("小数点以下が四捨五入されること", () => {
      render(<ProgressRing progress={33.7} />);
      expect(screen.getByText("34")).toBeInTheDocument();
    });

    it("0%が正しく表示されること", () => {
      render(<ProgressRing progress={0} />);
      expect(screen.getByText("0")).toBeInTheDocument();
    });

    it("100%が正しく表示されること", () => {
      render(<ProgressRing progress={100} />);
      expect(screen.getByText("100")).toBeInTheDocument();
    });
  });

  describe("SVG要素", () => {
    it("SVG要素がレンダリングされること", () => {
      const { container } = render(<ProgressRing progress={50} />);
      expect(container.querySelector("svg")).toBeInTheDocument();
    });

    it("2つのcircle要素が存在すること", () => {
      const { container } = render(<ProgressRing progress={50} />);
      const circles = container.querySelectorAll("circle");
      expect(circles).toHaveLength(2);
    });
  });

  describe("サイズオプション", () => {
    it("デフォルトサイズが120であること", () => {
      const { container } = render(<ProgressRing progress={50} />);
      const svg = container.querySelector("svg");
      expect(svg).toHaveAttribute("width", "120");
      expect(svg).toHaveAttribute("height", "120");
    });

    it("カスタムサイズが適用されること", () => {
      const { container } = render(<ProgressRing progress={50} size={200} />);
      const svg = container.querySelector("svg");
      expect(svg).toHaveAttribute("width", "200");
      expect(svg).toHaveAttribute("height", "200");
    });
  });

  describe("strokeWidthオプション", () => {
    it("デフォルトstrokeWidthが8であること", () => {
      const { container } = render(<ProgressRing progress={50} />);
      const circles = container.querySelectorAll("circle");
      circles.forEach((circle) => {
        expect(circle).toHaveAttribute("stroke-width", "8");
      });
    });

    it("カスタムstrokeWidthが適用されること", () => {
      const { container } = render(<ProgressRing progress={50} strokeWidth={12} />);
      const circles = container.querySelectorAll("circle");
      circles.forEach((circle) => {
        expect(circle).toHaveAttribute("stroke-width", "12");
      });
    });
  });

  describe("プログレス計算", () => {
    it("0%の場合、strokeDashoffsetが最大になること", () => {
      const { container } = render(<ProgressRing progress={0} />);
      const progressCircle = container.querySelectorAll("circle")[1];
      const radius = (120 - 8) / 2;
      const circumference = radius * 2 * Math.PI;
      expect(progressCircle).toHaveAttribute(
        "stroke-dashoffset",
        String(circumference)
      );
    });

    it("100%の場合、strokeDashoffsetが0になること", () => {
      const { container } = render(<ProgressRing progress={100} />);
      const progressCircle = container.querySelectorAll("circle")[1];
      expect(progressCircle).toHaveAttribute("stroke-dashoffset", "0");
    });

    it("50%の場合、strokeDashoffsetが半分になること", () => {
      const { container } = render(<ProgressRing progress={50} />);
      const progressCircle = container.querySelectorAll("circle")[1];
      const radius = (120 - 8) / 2;
      const circumference = radius * 2 * Math.PI;
      const expectedOffset = circumference - (50 / 100) * circumference;
      expect(progressCircle).toHaveAttribute(
        "stroke-dashoffset",
        String(expectedOffset)
      );
    });
  });

  describe("境界値", () => {
    it("負の値が渡された場合も動作すること", () => {
      const { container } = render(<ProgressRing progress={-10} />);
      expect(screen.getByText("-10")).toBeInTheDocument();
      expect(container.querySelector("svg")).toBeInTheDocument();
    });

    it("100を超える値が渡された場合も動作すること", () => {
      const { container } = render(<ProgressRing progress={150} />);
      expect(screen.getByText("150")).toBeInTheDocument();
      expect(container.querySelector("svg")).toBeInTheDocument();
    });
  });
});
