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

    expect(screen.getByText("AIによるアドバイス")).toBeInTheDocument();
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

  it("ローディング中は前の情報が表示されない", () => {
    const { container } = render(
      <NutritionAdviceCard
        advice="前回のアドバイス"
        isLoading={true}
        error={null}
      />
    );

    // スケルトンが表示される
    const skeletons = container.querySelectorAll('[class*="animate-pulse"]');
    expect(skeletons.length).toBeGreaterThan(0);

    // 前の情報は表示されない
    expect(screen.queryByText("前回のアドバイス")).not.toBeInTheDocument();
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
