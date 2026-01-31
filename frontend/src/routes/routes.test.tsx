/**
 * ルーティングのテスト
 *
 * 各ルートが正しくレンダリングされることを確認
 */
import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { createMemoryRouter, RouterProvider } from "react-router-dom";
import { routes } from "./index";

// fetchをモック化
global.fetch = vi.fn(() =>
  Promise.resolve({
    json: () => Promise.resolve({ status: "healthy", message: "OK" }),
  } as Response)
);

describe("ルーティング", () => {
  describe("/register ルート", () => {
    it("RegisterPageが表示される", async () => {
      const router = createMemoryRouter(routes, {
        initialEntries: ["/register"],
      });

      render(<RouterProvider router={router} />);

      // RegisterFormのタイトルが表示されることを確認
      expect(await screen.findByText("新規登録")).toBeInTheDocument();
    });

    it("登録フォームの主要フィールドが表示される", async () => {
      const router = createMemoryRouter(routes, {
        initialEntries: ["/register"],
      });

      render(<RouterProvider router={router} />);

      // 主要なフォーム要素が表示されることを確認
      expect(await screen.findByLabelText("メールアドレス")).toBeInTheDocument();
      expect(screen.getByLabelText("パスワード")).toBeInTheDocument();
      expect(screen.getByLabelText("ニックネーム")).toBeInTheDocument();
      expect(screen.getByRole("button", { name: "登録する" })).toBeInTheDocument();
    });
  });

  describe("/ ルート", () => {
    it("HomePageが表示される", async () => {
      const router = createMemoryRouter(routes, {
        initialEntries: ["/"],
      });

      render(<RouterProvider router={router} />);

      // HomePageのタイトルが表示されることを確認
      expect(await screen.findByText("CalTrack")).toBeInTheDocument();
    });
  });
});
