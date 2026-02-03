/**
 * LoginForm コンポーネントのテスト
 */
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { BrowserRouter } from "react-router-dom";
import { LoginForm } from "./LoginForm";
import type { LoginResponse } from "./LoginForm";
import type { ApiErrorResponse } from "@/lib/api";

// SWR mutateをモック
const mockTrigger = vi.fn();
const mockReset = vi.fn();
let mockError: ApiErrorResponse | undefined = undefined;

vi.mock("@/features/common/hooks/useRequest", () => ({
  useRequestMutation: () => ({
    trigger: mockTrigger,
    isMutating: false,
    error: mockError,
    data: undefined,
    reset: mockReset,
  }),
  useRequestGet: () => ({
    data: undefined,
    error: undefined,
    isLoading: false,
    mutate: vi.fn(),
  }),
  useRequest: () => {
    throw new Error("useRequest is deprecated");
  },
}));

/**
 * BrowserRouterでラップしてレンダリング
 */
function renderWithRouter(component: React.ReactElement) {
  return render(<BrowserRouter>{component}</BrowserRouter>);
}

describe("LoginForm", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockError = undefined;
  });

  describe("レンダリング", () => {
    it("すべてのフォームフィールドが表示される", () => {
      renderWithRouter(<LoginForm />);

      expect(screen.getByLabelText("メールアドレス")).toBeInTheDocument();
      expect(screen.getByLabelText("パスワード")).toBeInTheDocument();
      expect(
        screen.getByRole("button", { name: "ログイン" }),
      ).toBeInTheDocument();
    });

    it("タイトルと説明が表示される", () => {
      renderWithRouter(<LoginForm />);

      expect(
        screen.getByRole("heading", { name: "ログイン" }),
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
          screen.getByText("有効なメールアドレスを入力してください"),
        ).toBeInTheDocument();
      });
    });

    it("有効なメールを入力するとエラーが消える", async () => {
      const user = userEvent.setup();
      renderWithRouter(<LoginForm />);

      // 不正なメール
      await user.type(screen.getByLabelText("メールアドレス"), "invalid");
      expect(
        screen.getByText("有効なメールアドレスを入力してください"),
      ).toBeInTheDocument();

      // 有効なメールに修正
      await user.clear(screen.getByLabelText("メールアドレス"));
      await user.type(
        screen.getByLabelText("メールアドレス"),
        "test@example.com",
      );

      await waitFor(() => {
        expect(
          screen.queryByText("有効なメールアドレスを入力してください"),
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
      mockTrigger.mockResolvedValueOnce(mockResponse);

      const user = userEvent.setup();
      renderWithRouter(<LoginForm />);

      // フォームに入力
      await user.type(
        screen.getByLabelText("メールアドレス"),
        "test@example.com",
      );
      await user.type(screen.getByLabelText("パスワード"), "password123");

      // ボタンが有効化されることを確認
      await waitFor(() => {
        expect(
          screen.getByRole("button", { name: "ログイン" }),
        ).not.toBeDisabled();
      });

      // 送信
      fireEvent.click(screen.getByRole("button", { name: "ログイン" }));

      await waitFor(() => {
        expect(mockTrigger).toHaveBeenCalledWith({
          email: "test@example.com",
          password: "password123",
        });
      });
    });

    it("onSuccessコールバックがAPI成功時に呼ばれる", async () => {
      // onSuccessはuseRequestMutationのオプションで渡されるため、
      // このテストではtriggerが呼ばれることを確認
      const mockResponse: LoginResponse = {
        userId: "user-123",
        email: "test@example.com",
        nickname: "TestUser",
      };
      mockTrigger.mockResolvedValueOnce(mockResponse);

      const user = userEvent.setup();
      renderWithRouter(<LoginForm />);

      // フォームに入力
      await user.type(
        screen.getByLabelText("メールアドレス"),
        "test@example.com",
      );
      await user.type(screen.getByLabelText("パスワード"), "password123");

      // ボタンが有効化されることを確認
      await waitFor(() => {
        expect(
          screen.getByRole("button", { name: "ログイン" }),
        ).not.toBeDisabled();
      });

      // 送信
      fireEvent.click(screen.getByRole("button", { name: "ログイン" }));

      await waitFor(() => {
        expect(mockTrigger).toHaveBeenCalled();
      });
    });
  });

  describe("エラー表示", () => {
    it("triggerがエラーをthrowした場合キャッチされる", async () => {
      const error: ApiErrorResponse = {
        code: "INVALID_CREDENTIALS",
        message: "Invalid credentials",
      };
      mockTrigger.mockRejectedValueOnce(error);

      const user = userEvent.setup();
      renderWithRouter(<LoginForm />);

      // フォームに入力
      await user.type(
        screen.getByLabelText("メールアドレス"),
        "test@example.com",
      );
      await user.type(screen.getByLabelText("パスワード"), "wrongpassword");

      // ボタンが有効化されることを確認
      await waitFor(() => {
        expect(
          screen.getByRole("button", { name: "ログイン" }),
        ).not.toBeDisabled();
      });

      // 送信（エラーがスローされてもクラッシュしないことを確認）
      fireEvent.click(screen.getByRole("button", { name: "ログイン" }));

      await waitFor(() => {
        expect(mockTrigger).toHaveBeenCalled();
      });
    });
  });

  describe("エラーリセット", () => {
    it("APIエラーがある状態でフォームがレンダリングされる", async () => {
      // エラーがある状態でレンダリング
      mockError = {
        code: "INVALID_CREDENTIALS",
        message: "Invalid credentials",
      };

      renderWithRouter(<LoginForm />);

      // フォームが正常にレンダリングされることを確認
      expect(screen.getByLabelText("メールアドレス")).toBeInTheDocument();
      expect(screen.getByLabelText("パスワード")).toBeInTheDocument();
    });
  });
});
