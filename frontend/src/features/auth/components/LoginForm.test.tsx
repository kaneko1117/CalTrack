/**
 * LoginForm コンポーネントのテスト
 */
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { BrowserRouter } from "react-router-dom";
import { LoginForm } from "./LoginForm";
import type { LoginResponse } from "./LoginForm";
import * as api from "@/lib/api";

// lib/api の post関数をモック
vi.mock("@/lib/api", async () => {
  const actual = await vi.importActual<typeof api>("@/lib/api");
  return {
    ...actual,
    post: vi.fn(),
  };
});

const mockPost = vi.mocked(api.post);

/**
 * BrowserRouterでラップしてレンダリング
 */
function renderWithRouter(component: React.ReactElement) {
  return render(<BrowserRouter>{component}</BrowserRouter>);
}

describe("LoginForm", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("レンダリング", () => {
    it("すべてのフォームフィールドが表示される", () => {
      renderWithRouter(<LoginForm />);

      expect(screen.getByLabelText("メールアドレス")).toBeInTheDocument();
      expect(screen.getByLabelText("パスワード")).toBeInTheDocument();
      expect(
        screen.getByRole("button", { name: "ログイン" })
      ).toBeInTheDocument();
    });

    it("タイトルと説明が表示される", () => {
      renderWithRouter(<LoginForm />);

      expect(
        screen.getByRole("heading", { name: "ログイン" })
      ).toBeInTheDocument();
      expect(
        screen.getByText("アカウントにログインしてください")
      ).toBeInTheDocument();
    });

    it("新規登録リンクが表示される", () => {
      renderWithRouter(<LoginForm />);

      const registerLink = screen.getByRole("link", { name: "新規登録" });
      expect(registerLink).toBeInTheDocument();
      expect(registerLink).toHaveAttribute("href", "/register");
    });

    it("初期状態でログインボタンが無効化される", () => {
      renderWithRouter(<LoginForm />);

      const submitButton = screen.getByRole("button", { name: "ログイン" });
      expect(submitButton).toBeDisabled();
    });
  });

  describe("バリデーション", () => {
    it("不正なメール形式でエラーが表示される", async () => {
      const user = userEvent.setup();
      renderWithRouter(<LoginForm />);

      await user.type(screen.getByLabelText("メールアドレス"), "invalid-email");

      await waitFor(() => {
        expect(
          screen.getByText("有効なメールアドレスを入力してください")
        ).toBeInTheDocument();
      });
    });

    it("有効なメールを入力するとエラーが消える", async () => {
      const user = userEvent.setup();
      renderWithRouter(<LoginForm />);

      // 不正なメール
      await user.type(screen.getByLabelText("メールアドレス"), "invalid");
      expect(
        screen.getByText("有効なメールアドレスを入力してください")
      ).toBeInTheDocument();

      // 有効なメールに修正
      await user.clear(screen.getByLabelText("メールアドレス"));
      await user.type(
        screen.getByLabelText("メールアドレス"),
        "test@example.com"
      );

      await waitFor(() => {
        expect(
          screen.queryByText("有効なメールアドレスを入力してください")
        ).not.toBeInTheDocument();
      });
    });
  });

  describe("フォーム送信", () => {
    it("有効なデータでフォームを送信するとAPI呼び出しが行われる", async () => {
      const mockResponse: LoginResponse = {
        userId: "user-123",
        email: "test@example.com",
        nickname: "TestUser",
      };
      mockPost.mockResolvedValueOnce(mockResponse);

      const user = userEvent.setup();
      renderWithRouter(<LoginForm />);

      // フォームに入力
      await user.type(
        screen.getByLabelText("メールアドレス"),
        "test@example.com"
      );
      await user.type(screen.getByLabelText("パスワード"), "password123");

      // ボタンが有効化されることを確認
      await waitFor(() => {
        expect(
          screen.getByRole("button", { name: "ログイン" })
        ).not.toBeDisabled();
      });

      // 送信
      fireEvent.click(screen.getByRole("button", { name: "ログイン" }));

      await waitFor(() => {
        expect(mockPost).toHaveBeenCalledWith("/api/v1/auth/login", {
          email: "test@example.com",
          password: "password123",
        });
      });
    });

    it("onSuccessコールバックがAPI成功時に呼ばれる", async () => {
      const mockResponse: LoginResponse = {
        userId: "user-123",
        email: "test@example.com",
        nickname: "TestUser",
      };
      mockPost.mockResolvedValueOnce(mockResponse);

      const onSuccess = vi.fn();
      const user = userEvent.setup();
      renderWithRouter(<LoginForm onSuccess={onSuccess} />);

      // フォームに入力
      await user.type(
        screen.getByLabelText("メールアドレス"),
        "test@example.com"
      );
      await user.type(screen.getByLabelText("パスワード"), "password123");

      // ボタンが有効化されることを確認
      await waitFor(() => {
        expect(
          screen.getByRole("button", { name: "ログイン" })
        ).not.toBeDisabled();
      });

      // 送信
      fireEvent.click(screen.getByRole("button", { name: "ログイン" }));

      await waitFor(() => {
        expect(onSuccess).toHaveBeenCalledWith(mockResponse);
      });
    });
  });

  describe("エラー表示", () => {
    it("INVALID_CREDENTIALSエラーが表示される", async () => {
      const error: api.ApiErrorResponse = {
        code: "INVALID_CREDENTIALS",
        message: "Invalid credentials",
      };
      mockPost.mockRejectedValueOnce(error);

      const user = userEvent.setup();
      renderWithRouter(<LoginForm />);

      // フォームに入力
      await user.type(
        screen.getByLabelText("メールアドレス"),
        "test@example.com"
      );
      await user.type(screen.getByLabelText("パスワード"), "wrongpassword");

      // ボタンが有効化されることを確認
      await waitFor(() => {
        expect(
          screen.getByRole("button", { name: "ログイン" })
        ).not.toBeDisabled();
      });

      // 送信
      fireEvent.click(screen.getByRole("button", { name: "ログイン" }));

      await waitFor(() => {
        expect(
          screen.getByText("メールアドレスまたはパスワードが間違っています")
        ).toBeInTheDocument();
      });
    });

    it("VALIDATION_ERRORが表示される", async () => {
      const error: api.ApiErrorResponse = {
        code: "VALIDATION_ERROR",
        message: "Validation failed",
      };
      mockPost.mockRejectedValueOnce(error);

      const user = userEvent.setup();
      renderWithRouter(<LoginForm />);

      // フォームに入力
      await user.type(
        screen.getByLabelText("メールアドレス"),
        "test@example.com"
      );
      await user.type(screen.getByLabelText("パスワード"), "password123");

      // ボタンが有効化されることを確認
      await waitFor(() => {
        expect(
          screen.getByRole("button", { name: "ログイン" })
        ).not.toBeDisabled();
      });

      // 送信
      fireEvent.click(screen.getByRole("button", { name: "ログイン" }));

      await waitFor(() => {
        expect(screen.getByText("入力内容に誤りがあります")).toBeInTheDocument();
      });
    });

    it("INTERNAL_ERRORが表示される", async () => {
      const error: api.ApiErrorResponse = {
        code: "INTERNAL_ERROR",
        message: "Internal error",
      };
      mockPost.mockRejectedValueOnce(error);

      const user = userEvent.setup();
      renderWithRouter(<LoginForm />);

      // フォームに入力
      await user.type(
        screen.getByLabelText("メールアドレス"),
        "test@example.com"
      );
      await user.type(screen.getByLabelText("パスワード"), "password123");

      // ボタンが有効化されることを確認
      await waitFor(() => {
        expect(
          screen.getByRole("button", { name: "ログイン" })
        ).not.toBeDisabled();
      });

      // 送信
      fireEvent.click(screen.getByRole("button", { name: "ログイン" }));

      await waitFor(() => {
        expect(
          screen.getByText("予期しないエラーが発生しました")
        ).toBeInTheDocument();
      });
    });
  });

  describe("エラーリセット", () => {
    it("APIエラーがフィールド入力時にクリアされる", async () => {
      const error: api.ApiErrorResponse = {
        code: "INVALID_CREDENTIALS",
        message: "Invalid credentials",
      };
      mockPost.mockRejectedValueOnce(error);

      const user = userEvent.setup();
      renderWithRouter(<LoginForm />);

      // フォームに入力
      await user.type(
        screen.getByLabelText("メールアドレス"),
        "test@example.com"
      );
      await user.type(screen.getByLabelText("パスワード"), "wrongpassword");

      // ボタンが有効化されることを確認
      await waitFor(() => {
        expect(
          screen.getByRole("button", { name: "ログイン" })
        ).not.toBeDisabled();
      });

      // 送信
      fireEvent.click(screen.getByRole("button", { name: "ログイン" }));

      // エラーが表示されることを確認
      await waitFor(() => {
        expect(
          screen.getByText("メールアドレスまたはパスワードが間違っています")
        ).toBeInTheDocument();
      });

      // フィールドに入力するとエラーがクリアされる
      await user.type(screen.getByLabelText("パスワード"), "a");

      await waitFor(() => {
        expect(
          screen.queryByText("メールアドレスまたはパスワードが間違っています")
        ).not.toBeInTheDocument();
      });
    });
  });
});
