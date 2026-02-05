/**
 * NutritionAdviceCard テスト
 */
import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { NutritionAdviceCard } from "./NutritionAdviceCard";

describe("NutritionAdviceCard", () => {
  it("アドバイスが正しく表示される", () => {
    render(
      <NutritionAdviceCard
        advice="タンパク質が不足しています。鶏肉や卵を食べましょう。"
        isLoading={false}
        error={null}
      />
    );

    expect(screen.getByText("今日のアドバイス")).toBeInTheDocument();
    expect(
      screen.getByText("タンパク質が不足しています。鶏肉や卵を食べましょう。")
    ).toBeInTheDocument();
  });

  it("ローディング中はスケルトンが表示される", () => {
    const { container } = render(
      <NutritionAdviceCard advice={null} isLoading={true} error={null} />
    );

    // Skeletonコンポーネントが存在することを確認
    const skeletons = container.querySelectorAll('[class*="animate-pulse"]');
    expect(skeletons.length).toBeGreaterThan(0);
  });

  it("エラー時はエラーメッセージが表示される", () => {
    render(
      <NutritionAdviceCard
        advice={null}
        isLoading={false}
        error={{ code: "INTERNAL_ERROR", message: "Server error" }}
      />
    );

    expect(
      screen.getByText("アドバイスの取得に失敗しました")
    ).toBeInTheDocument();
  });

  it("アドバイスがない場合は案内メッセージが表示される", () => {
    render(
      <NutritionAdviceCard advice={null} isLoading={false} error={null} />
    );

    expect(
      screen.getByText("食事を記録するとアドバイスが表示されます")
    ).toBeInTheDocument();
  });
});
