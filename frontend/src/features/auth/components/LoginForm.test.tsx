/**
 * LoginForm コンポーネントのテスト
 */
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { BrowserRouter } from "react-router-dom";
import { LoginForm } from "./LoginForm";
import * as hooks from "../hooks";
import { ApiError } from "../api";

// useLoginフックをモック
vi.mock("../hooks", async () => {
  const actual = await vi.importActual<typeof hooks>("../hooks");
  return {
    ...actual,
    useLogin: vi.fn(),
  };
});

const mockUseLogin = vi.mocked(hooks.useLogin);

/**
 * BrowserRouterでラップしてレンダリング
 */
function renderWithRouter(component: React.ReactElement) {
  return render(<BrowserRouter>{component}</BrowserRouter>);
}

describe("LoginForm", () => {
  const mockLogin = vi.fn();
  const mockReset = vi.fn();

  const defaultHookReturn: hooks.UseLoginReturn = {
    login: mockLogin,
    isLoading: false,
    error: null,
    isSuccess: false,
    reset: mockReset,
  };

  beforeEach(() => {
    vi.clearAllMocks();
    mockUseLogin.mockReturnValue(defaultHookReturn);
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

      // h3要素のタイトルを取得（ボタンの「ログイン」と区別するためroleを使用）
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
  });

  describe("バリデーション", () => {
    it("空のフォームを送信するとバリデーションエラーが表示される", async () => {
      renderWithRouter(<LoginForm />);

      fireEvent.click(screen.getByRole("button", { name: "ログイン" }));

      await waitFor(() => {
        expect(
          screen.getByText("メールアドレスを入力してください")
        ).toBeInTheDocument();
        expect(
          screen.getByText("パスワードを入力してください")
        ).toBeInTheDocument();
      });

      // login関数は呼ばれない
      expect(mockLogin).not.toHaveBeenCalled();
    });

    it("不正なメール形式でエラーが表示される", async () => {
      renderWithRouter(<LoginForm />);

      // fireEvent.changeで直接値を設定
      fireEvent.change(screen.getByLabelText("メールアドレス"), {
        target: { value: "invalid-email" },
      });
      // フォームを直接サブミット（HTML5バリデーションをバイパス）
      const form = screen
        .getByRole("button", { name: "ログイン" })
        .closest("form");
      fireEvent.submit(form!);

      await waitFor(() => {
        expect(
          screen.getByText("正しいメールアドレス形式で入力してください")
        ).toBeInTheDocument();
      });
    });

    it("メールアドレスのみ入力でパスワードエラーが表示される", async () => {
      const user = userEvent.setup();
      renderWithRouter(<LoginForm />);

      await user.type(
        screen.getByLabelText("メールアドレス"),
        "test@example.com"
      );
      fireEvent.click(screen.getByRole("button", { name: "ログイン" }));

      await waitFor(() => {
        expect(
          screen.getByText("パスワードを入力してください")
        ).toBeInTheDocument();
        expect(
          screen.queryByText("メールアドレスを入力してください")
        ).not.toBeInTheDocument();
      });
    });
  });

  describe("フォーム送信", () => {
    it("有効なデータでフォームを送信するとlogin関数が呼ばれる", async () => {
      const user = userEvent.setup();
      renderWithRouter(<LoginForm />);

      // フォームに入力
      await user.type(
        screen.getByLabelText("メールアドレス"),
        "test@example.com"
      );
      await user.type(screen.getByLabelText("パスワード"), "password123");

      // 送信
      fireEvent.click(screen.getByRole("button", { name: "ログイン" }));

      await waitFor(() => {
        expect(mockLogin).toHaveBeenCalledWith(
          {
            email: "test@example.com",
            password: "password123",
          },
          undefined // onSuccessコールバック（未指定時はundefined）
        );
      });
    });

    it("onSuccessコールバックがlogin関数に渡される", async () => {
      const onSuccess = vi.fn();
      const user = userEvent.setup();
      renderWithRouter(<LoginForm onSuccess={onSuccess} />);

      // フォームに入力
      await user.type(
        screen.getByLabelText("メールアドレス"),
        "test@example.com"
      );
      await user.type(screen.getByLabelText("パスワード"), "password123");

      // 送信
      fireEvent.click(screen.getByRole("button", { name: "ログイン" }));

      await waitFor(() => {
        expect(mockLogin).toHaveBeenCalledWith(
          expect.any(Object),
          onSuccess // onSuccessコールバックが渡される
        );
      });
    });
  });

  describe("ローディング状態", () => {
    it("isLoadingがtrueの時、ボタンが無効化される", () => {
      mockUseLogin.mockReturnValue({
        ...defaultHookReturn,
        isLoading: true,
      });

      renderWithRouter(<LoginForm />);

      const submitButton = screen.getByRole("button", { name: "ログイン中..." });
      expect(submitButton).toBeDisabled();
    });

    it("isLoadingがtrueの時、フォームフィールドが無効化される", () => {
      mockUseLogin.mockReturnValue({
        ...defaultHookReturn,
        isLoading: true,
      });

      renderWithRouter(<LoginForm />);

      expect(screen.getByLabelText("メールアドレス")).toBeDisabled();
      expect(screen.getByLabelText("パスワード")).toBeDisabled();
    });
  });

  describe("エラー表示", () => {
    it("INVALID_CREDENTIALSエラーが表示される", () => {
      const error = new ApiError(
        "INVALID_CREDENTIALS",
        "Invalid credentials",
        401
      );
      mockUseLogin.mockReturnValue({
        ...defaultHookReturn,
        error,
      });

      renderWithRouter(<LoginForm />);

      expect(
        screen.getByText("メールアドレスまたはパスワードが間違っています")
      ).toBeInTheDocument();
    });

    it("VALIDATION_ERRORとdetailsが表示される", () => {
      const error = new ApiError("VALIDATION_ERROR", "Validation failed", 400, [
        "email: 不正な形式です",
        "password: 必須項目です",
      ]);
      mockUseLogin.mockReturnValue({
        ...defaultHookReturn,
        error,
      });

      renderWithRouter(<LoginForm />);

      expect(screen.getByText("入力内容に誤りがあります")).toBeInTheDocument();
      expect(screen.getByText("email: 不正な形式です")).toBeInTheDocument();
      expect(screen.getByText("password: 必須項目です")).toBeInTheDocument();
    });

    it("INTERNAL_ERRORが表示される", () => {
      const error = new ApiError("INTERNAL_ERROR", "Internal error", 500);
      mockUseLogin.mockReturnValue({
        ...defaultHookReturn,
        error,
      });

      renderWithRouter(<LoginForm />);

      expect(
        screen.getByText("予期しないエラーが発生しました")
      ).toBeInTheDocument();
    });
  });

  describe("エラーリセット", () => {
    it("フィールド入力時にエラーがリセットされる", async () => {
      const error = new ApiError(
        "INVALID_CREDENTIALS",
        "Invalid credentials",
        401
      );
      mockUseLogin.mockReturnValue({
        ...defaultHookReturn,
        error,
      });

      const user = userEvent.setup();
      renderWithRouter(<LoginForm />);

      // エラーが表示されていることを確認
      expect(
        screen.getByText("メールアドレスまたはパスワードが間違っています")
      ).toBeInTheDocument();

      // フィールドに入力
      await user.type(screen.getByLabelText("メールアドレス"), "a");

      // reset関数が呼ばれることを確認
      expect(mockReset).toHaveBeenCalled();
    });

    it("バリデーションエラーがフィールド入力時にクリアされる", async () => {
      const user = userEvent.setup();
      renderWithRouter(<LoginForm />);

      // 空のまま送信してバリデーションエラーを発生させる
      fireEvent.click(screen.getByRole("button", { name: "ログイン" }));

      await waitFor(() => {
        expect(
          screen.getByText("メールアドレスを入力してください")
        ).toBeInTheDocument();
      });

      // フィールドに入力
      await user.type(
        screen.getByLabelText("メールアドレス"),
        "test@example.com"
      );

      // バリデーションエラーがクリアされる
      expect(
        screen.queryByText("メールアドレスを入力してください")
      ).not.toBeInTheDocument();
    });
  });

  describe("ログイン成功", () => {
    it("login関数成功時にonSuccessコールバックが呼ばれる", async () => {
      const onSuccess = vi.fn();
      const mockResponse = {
        userId: "user-123",
        email: "test@example.com",
        nickname: "TestUser",
      };

      mockUseLogin.mockImplementation(() => ({
        ...defaultHookReturn,
        login: vi.fn().mockImplementation(async (_request, callback) => {
          // 成功をシミュレート
          callback?.(mockResponse);
        }),
        isSuccess: true,
      }));

      const user = userEvent.setup();
      renderWithRouter(<LoginForm onSuccess={onSuccess} />);

      // フォームに入力
      await user.type(
        screen.getByLabelText("メールアドレス"),
        "test@example.com"
      );
      await user.type(screen.getByLabelText("パスワード"), "password123");

      // 送信
      fireEvent.click(screen.getByRole("button", { name: "ログイン" }));

      await waitFor(() => {
        expect(onSuccess).toHaveBeenCalledWith(mockResponse);
      });
    });
  });
});
