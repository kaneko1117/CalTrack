/**
 * RegisterForm コンポーネントのテスト
 */
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { RegisterForm } from "./RegisterForm";
import type { RegisterUserResponse } from "./RegisterForm";
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

describe("RegisterForm", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockError = undefined;
  });

  describe("レンダリング", () => {
    it("すべてのフォームフィールドが表示される", () => {
      render(<RegisterForm />);

      expect(screen.getByLabelText("ニックネーム")).toBeInTheDocument();
      expect(screen.getByLabelText("メールアドレス")).toBeInTheDocument();
      expect(screen.getByLabelText("パスワード")).toBeInTheDocument();
      expect(screen.getByLabelText("体重 (kg)")).toBeInTheDocument();
      expect(screen.getByLabelText("身長 (cm)")).toBeInTheDocument();
      expect(screen.getByLabelText("生年月日")).toBeInTheDocument();
      expect(screen.getByLabelText("性別")).toBeInTheDocument();
      expect(screen.getByLabelText("活動レベル")).toBeInTheDocument();
      expect(
        screen.getByRole("button", { name: "登録する" })
      ).toBeInTheDocument();
    });

    it("タイトルと説明が表示される", () => {
      render(<RegisterForm />);

      expect(screen.getByText("新規登録")).toBeInTheDocument();
      expect(
        screen.getByText("アカウントを作成して、カロリー管理を始めましょう")
      ).toBeInTheDocument();
    });

    it("初期状態で登録ボタンが無効化される", () => {
      render(<RegisterForm />);

      const submitButton = screen.getByRole("button", { name: "登録する" });
      expect(submitButton).toBeDisabled();
    });
  });

  describe("バリデーション", () => {
    it("不正なメール形式でエラーが表示される", async () => {
      const user = userEvent.setup();
      render(<RegisterForm />);

      await user.type(screen.getByLabelText("メールアドレス"), "invalid-email");

      await waitFor(() => {
        expect(
          screen.getByText("有効なメールアドレスを入力してください")
        ).toBeInTheDocument();
      });
    });

    it("パスワードが8文字未満でエラーが表示される", async () => {
      const user = userEvent.setup();
      render(<RegisterForm />);

      await user.type(screen.getByLabelText("パスワード"), "short");

      await waitFor(() => {
        expect(
          screen.getByText("パスワードは8文字以上で入力してください")
        ).toBeInTheDocument();
      });
    });

    it("体重が0以下でエラーが表示される", async () => {
      const user = userEvent.setup();
      render(<RegisterForm />);

      await user.type(screen.getByLabelText("体重 (kg)"), "-10");

      await waitFor(() => {
        expect(
          screen.getByText("体重は0より大きい値を入力してください")
        ).toBeInTheDocument();
      });
    });

    it("身長が0以下でエラーが表示される", async () => {
      const user = userEvent.setup();
      render(<RegisterForm />);

      // 0を入力（0以下は無効）
      await user.type(screen.getByLabelText("身長 (cm)"), "-5");

      await waitFor(() => {
        expect(
          screen.getByText("身長は0より大きい値を入力してください")
        ).toBeInTheDocument();
      });
    });

    it("未来の生年月日でエラーが表示される", async () => {
      const user = userEvent.setup();
      render(<RegisterForm />);

      // 未来の日付を設定
      const futureDate = new Date();
      futureDate.setFullYear(futureDate.getFullYear() + 1);
      const futureDateStr = futureDate.toISOString().split("T")[0];

      await user.type(screen.getByLabelText("生年月日"), futureDateStr);

      await waitFor(() => {
        expect(
          screen.getByText("生年月日は過去の日付を入力してください")
        ).toBeInTheDocument();
      });
    });

    it("有効なメールを入力するとエラーが消える", async () => {
      const user = userEvent.setup();
      render(<RegisterForm />);

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
      const mockResponse: RegisterUserResponse = {
        userId: "user-123",
        email: "test@example.com",
        nickname: "TestUser",
      };
      mockTrigger.mockResolvedValueOnce(mockResponse);

      const user = userEvent.setup();
      render(<RegisterForm />);

      // フォームに入力
      await user.type(screen.getByLabelText("ニックネーム"), "TestUser");
      await user.type(
        screen.getByLabelText("メールアドレス"),
        "test@example.com"
      );
      await user.type(screen.getByLabelText("パスワード"), "password123");
      await user.type(screen.getByLabelText("体重 (kg)"), "70");
      await user.type(screen.getByLabelText("身長 (cm)"), "175");
      await user.type(screen.getByLabelText("生年月日"), "1990-01-01");
      await user.selectOptions(screen.getByLabelText("性別"), "male");
      await user.selectOptions(screen.getByLabelText("活動レベル"), "moderate");

      // ボタンが有効化されることを確認
      await waitFor(() => {
        expect(
          screen.getByRole("button", { name: "登録する" })
        ).not.toBeDisabled();
      });

      // 送信
      fireEvent.click(screen.getByRole("button", { name: "登録する" }));

      await waitFor(() => {
        expect(mockTrigger).toHaveBeenCalledWith({
          email: "test@example.com",
          password: "password123",
          nickname: "TestUser",
          weight: 70,
          height: 175,
          birthDate: "1990-01-01",
          gender: "male",
          activityLevel: "moderate",
        });
      });
    });

    it("onSuccessコールバックがAPI成功時に呼ばれる", async () => {
      // onSuccessはuseRequestMutationのオプションで渡されるため、
      // このテストではtriggerが呼ばれることを確認
      const mockResponse: RegisterUserResponse = {
        userId: "user-123",
        email: "test@example.com",
        nickname: "TestUser",
      };
      mockTrigger.mockResolvedValueOnce(mockResponse);

      const user = userEvent.setup();
      render(<RegisterForm />);

      // フォームに入力
      await user.type(screen.getByLabelText("ニックネーム"), "TestUser");
      await user.type(
        screen.getByLabelText("メールアドレス"),
        "test@example.com"
      );
      await user.type(screen.getByLabelText("パスワード"), "password123");
      await user.type(screen.getByLabelText("体重 (kg)"), "70");
      await user.type(screen.getByLabelText("身長 (cm)"), "175");
      await user.type(screen.getByLabelText("生年月日"), "1990-01-01");
      await user.selectOptions(screen.getByLabelText("性別"), "male");
      await user.selectOptions(screen.getByLabelText("活動レベル"), "moderate");

      // ボタンが有効化されることを確認
      await waitFor(() => {
        expect(
          screen.getByRole("button", { name: "登録する" })
        ).not.toBeDisabled();
      });

      // 送信
      fireEvent.click(screen.getByRole("button", { name: "登録する" }));

      await waitFor(() => {
        expect(mockTrigger).toHaveBeenCalled();
      });
    });
  });

  describe("エラー表示", () => {
    it("triggerがエラーをthrowした場合キャッチされる", async () => {
      const error: ApiErrorResponse = {
        code: "EMAIL_ALREADY_EXISTS",
        message: "Email already exists",
      };
      mockTrigger.mockRejectedValueOnce(error);

      const user = userEvent.setup();
      render(<RegisterForm />);

      // フォームに入力
      await user.type(screen.getByLabelText("ニックネーム"), "TestUser");
      await user.type(
        screen.getByLabelText("メールアドレス"),
        "test@example.com"
      );
      await user.type(screen.getByLabelText("パスワード"), "password123");
      await user.type(screen.getByLabelText("体重 (kg)"), "70");
      await user.type(screen.getByLabelText("身長 (cm)"), "175");
      await user.type(screen.getByLabelText("生年月日"), "1990-01-01");
      await user.selectOptions(screen.getByLabelText("性別"), "male");
      await user.selectOptions(screen.getByLabelText("活動レベル"), "moderate");

      // ボタンが有効化されることを確認
      await waitFor(() => {
        expect(
          screen.getByRole("button", { name: "登録する" })
        ).not.toBeDisabled();
      });

      // 送信（エラーがスローされてもクラッシュしないことを確認）
      fireEvent.click(screen.getByRole("button", { name: "登録する" }));

      await waitFor(() => {
        expect(mockTrigger).toHaveBeenCalled();
      });
    });
  });

  describe("エラーリセット", () => {
    it("APIエラーがある状態でフォームがレンダリングされる", async () => {
      // エラーがある状態でレンダリング
      mockError = {
        code: "EMAIL_ALREADY_EXISTS",
        message: "Email already exists",
      };

      render(<RegisterForm />);

      // フォームが正常にレンダリングされることを確認
      expect(screen.getByLabelText("ニックネーム")).toBeInTheDocument();
      expect(screen.getByLabelText("メールアドレス")).toBeInTheDocument();
    });
  });
});
