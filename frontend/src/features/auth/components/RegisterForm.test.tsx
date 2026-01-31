/**
 * RegisterForm コンポーネントのテスト
 */
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { RegisterForm } from "./RegisterForm";
import * as hooks from "../hooks";
import { ApiError } from "../api";

// useRegisterUserフックをモック
vi.mock("../hooks", async () => {
  const actual = await vi.importActual<typeof hooks>("../hooks");
  return {
    ...actual,
    useRegisterUser: vi.fn(),
  };
});

const mockUseRegisterUser = vi.mocked(hooks.useRegisterUser);

describe("RegisterForm", () => {
  const mockRegister = vi.fn();
  const mockReset = vi.fn();

  const defaultHookReturn: hooks.UseRegisterUserReturn = {
    register: mockRegister,
    isLoading: false,
    error: null,
    isSuccess: false,
    reset: mockReset,
  };

  beforeEach(() => {
    vi.clearAllMocks();
    mockUseRegisterUser.mockReturnValue(defaultHookReturn);
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
      expect(screen.getByRole("button", { name: "登録する" })).toBeInTheDocument();
    });

    it("タイトルと説明が表示される", () => {
      render(<RegisterForm />);

      expect(screen.getByText("新規登録")).toBeInTheDocument();
      expect(
        screen.getByText("アカウントを作成して、カロリー管理を始めましょう")
      ).toBeInTheDocument();
    });
  });

  describe("バリデーション", () => {
    it("空のフォームを送信するとバリデーションエラーが表示される", async () => {
      render(<RegisterForm />);

      fireEvent.click(screen.getByRole("button", { name: "登録する" }));

      await waitFor(() => {
        expect(screen.getByText("ニックネームを入力してください")).toBeInTheDocument();
        expect(screen.getByText("メールアドレスを入力してください")).toBeInTheDocument();
        expect(screen.getByText("パスワードを入力してください")).toBeInTheDocument();
        expect(screen.getByText("体重を入力してください")).toBeInTheDocument();
        expect(screen.getByText("身長を入力してください")).toBeInTheDocument();
        expect(screen.getByText("生年月日を入力してください")).toBeInTheDocument();
        expect(screen.getByText("性別を選択してください")).toBeInTheDocument();
        expect(screen.getByText("活動レベルを選択してください")).toBeInTheDocument();
      });

      // register関数は呼ばれない
      expect(mockRegister).not.toHaveBeenCalled();
    });

    it("不正なメール形式でエラーが表示される", async () => {
      render(<RegisterForm />);

      // fireEvent.changeで直接値を設定
      fireEvent.change(screen.getByLabelText("メールアドレス"), {
        target: { value: "invalid-email" },
      });
      // フォームを直接サブミット（HTML5バリデーションをバイパス）
      const form = screen.getByRole("button", { name: "登録する" }).closest("form");
      fireEvent.submit(form!);

      await waitFor(() => {
        expect(
          screen.getByText("正しいメールアドレス形式で入力してください")
        ).toBeInTheDocument();
      });
    });

    it("パスワードが8文字未満でエラーが表示される", async () => {
      const user = userEvent.setup();
      render(<RegisterForm />);

      await user.type(screen.getByLabelText("パスワード"), "short");
      fireEvent.click(screen.getByRole("button", { name: "登録する" }));

      await waitFor(() => {
        expect(
          screen.getByText("パスワードは8文字以上で入力してください")
        ).toBeInTheDocument();
      });
    });

    it("体重が0以下でエラーが表示される", async () => {
      render(<RegisterForm />);

      // fireEvent.changeで直接値を設定
      fireEvent.change(screen.getByLabelText("体重 (kg)"), {
        target: { value: "-10" },
      });
      // フォームを直接サブミット（HTML5バリデーションをバイパス）
      const form = screen.getByRole("button", { name: "登録する" }).closest("form");
      fireEvent.submit(form!);

      await waitFor(() => {
        expect(screen.getByText("正しい体重を入力してください")).toBeInTheDocument();
      });
    });

    it("身長が0以下でエラーが表示される", async () => {
      const user = userEvent.setup();
      render(<RegisterForm />);

      await user.type(screen.getByLabelText("身長 (cm)"), "0");
      fireEvent.click(screen.getByRole("button", { name: "登録する" }));

      await waitFor(() => {
        expect(screen.getByText("正しい身長を入力してください")).toBeInTheDocument();
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
      fireEvent.click(screen.getByRole("button", { name: "登録する" }));

      await waitFor(() => {
        expect(screen.getByText("過去の日付を入力してください")).toBeInTheDocument();
      });
    });
  });

  describe("フォーム送信", () => {
    it("有効なデータでフォームを送信するとregister関数が呼ばれる", async () => {
      const user = userEvent.setup();
      render(<RegisterForm />);

      // フォームに入力
      await user.type(screen.getByLabelText("ニックネーム"), "TestUser");
      await user.type(screen.getByLabelText("メールアドレス"), "test@example.com");
      await user.type(screen.getByLabelText("パスワード"), "password123");
      await user.type(screen.getByLabelText("体重 (kg)"), "70");
      await user.type(screen.getByLabelText("身長 (cm)"), "175");
      await user.type(screen.getByLabelText("生年月日"), "1990-01-01");
      await user.selectOptions(screen.getByLabelText("性別"), "male");
      await user.selectOptions(screen.getByLabelText("活動レベル"), "moderate");

      // 送信
      fireEvent.click(screen.getByRole("button", { name: "登録する" }));

      await waitFor(() => {
        expect(mockRegister).toHaveBeenCalledWith(
          {
            email: "test@example.com",
            password: "password123",
            nickname: "TestUser",
            weight: 70,
            height: 175,
            birthDate: "1990-01-01",
            gender: "male",
            activityLevel: "moderate",
          },
          undefined // onSuccessコールバック（未指定時はundefined）
        );
      });
    });

    it("onSuccessコールバックがregister関数に渡される", async () => {
      const onSuccess = vi.fn();
      const user = userEvent.setup();
      render(<RegisterForm onSuccess={onSuccess} />);

      // フォームに入力
      await user.type(screen.getByLabelText("ニックネーム"), "TestUser");
      await user.type(screen.getByLabelText("メールアドレス"), "test@example.com");
      await user.type(screen.getByLabelText("パスワード"), "password123");
      await user.type(screen.getByLabelText("体重 (kg)"), "70");
      await user.type(screen.getByLabelText("身長 (cm)"), "175");
      await user.type(screen.getByLabelText("生年月日"), "1990-01-01");
      await user.selectOptions(screen.getByLabelText("性別"), "male");
      await user.selectOptions(screen.getByLabelText("活動レベル"), "moderate");

      // 送信
      fireEvent.click(screen.getByRole("button", { name: "登録する" }));

      await waitFor(() => {
        expect(mockRegister).toHaveBeenCalledWith(
          expect.any(Object),
          onSuccess // onSuccessコールバックが渡される
        );
      });
    });
  });

  describe("ローディング状態", () => {
    it("isLoadingがtrueの時、ボタンが無効化される", () => {
      mockUseRegisterUser.mockReturnValue({
        ...defaultHookReturn,
        isLoading: true,
      });

      render(<RegisterForm />);

      const submitButton = screen.getByRole("button", { name: "登録中..." });
      expect(submitButton).toBeDisabled();
    });

    it("isLoadingがtrueの時、フォームフィールドが無効化される", () => {
      mockUseRegisterUser.mockReturnValue({
        ...defaultHookReturn,
        isLoading: true,
      });

      render(<RegisterForm />);

      expect(screen.getByLabelText("ニックネーム")).toBeDisabled();
      expect(screen.getByLabelText("メールアドレス")).toBeDisabled();
      expect(screen.getByLabelText("パスワード")).toBeDisabled();
      expect(screen.getByLabelText("体重 (kg)")).toBeDisabled();
      expect(screen.getByLabelText("身長 (cm)")).toBeDisabled();
      expect(screen.getByLabelText("生年月日")).toBeDisabled();
      expect(screen.getByLabelText("性別")).toBeDisabled();
      expect(screen.getByLabelText("活動レベル")).toBeDisabled();
    });
  });

  describe("エラー表示", () => {
    it("EMAIL_ALREADY_EXISTSエラーが表示される", () => {
      const error = new ApiError(
        "EMAIL_ALREADY_EXISTS",
        "Email already exists",
        409
      );
      mockUseRegisterUser.mockReturnValue({
        ...defaultHookReturn,
        error,
      });

      render(<RegisterForm />);

      expect(
        screen.getByText("このメールアドレスは既に登録されています")
      ).toBeInTheDocument();
    });

    it("VALIDATION_ERRORとdetailsが表示される", () => {
      const error = new ApiError("VALIDATION_ERROR", "Validation failed", 400, [
        "email: 不正な形式です",
        "password: 短すぎます",
      ]);
      mockUseRegisterUser.mockReturnValue({
        ...defaultHookReturn,
        error,
      });

      render(<RegisterForm />);

      expect(screen.getByText("入力内容に誤りがあります")).toBeInTheDocument();
      expect(screen.getByText("email: 不正な形式です")).toBeInTheDocument();
      expect(screen.getByText("password: 短すぎます")).toBeInTheDocument();
    });

    it("INTERNAL_ERRORが表示される", () => {
      const error = new ApiError("INTERNAL_ERROR", "Internal error", 500);
      mockUseRegisterUser.mockReturnValue({
        ...defaultHookReturn,
        error,
      });

      render(<RegisterForm />);

      expect(
        screen.getByText("予期しないエラーが発生しました")
      ).toBeInTheDocument();
    });
  });

  describe("成功状態", () => {
    it("isSuccessがtrueの時、成功メッセージが表示される", () => {
      mockUseRegisterUser.mockReturnValue({
        ...defaultHookReturn,
        isSuccess: true,
      });

      render(<RegisterForm />);

      expect(screen.getByText("登録が完了しました")).toBeInTheDocument();
    });

    it("register関数成功時にonSuccessコールバックが呼ばれる", async () => {
      // register関数が成功時にonSuccessを呼び出すようにモック
      const onSuccess = vi.fn();
      mockUseRegisterUser.mockImplementation(() => ({
        ...defaultHookReturn,
        register: vi.fn().mockImplementation(async (_request, callback) => {
          // 成功をシミュレート
          callback?.();
        }),
        isSuccess: true,
      }));

      const user = userEvent.setup();
      render(<RegisterForm onSuccess={onSuccess} />);

      // フォームに入力
      await user.type(screen.getByLabelText("ニックネーム"), "TestUser");
      await user.type(screen.getByLabelText("メールアドレス"), "test@example.com");
      await user.type(screen.getByLabelText("パスワード"), "password123");
      await user.type(screen.getByLabelText("体重 (kg)"), "70");
      await user.type(screen.getByLabelText("身長 (cm)"), "175");
      await user.type(screen.getByLabelText("生年月日"), "1990-01-01");
      await user.selectOptions(screen.getByLabelText("性別"), "male");
      await user.selectOptions(screen.getByLabelText("活動レベル"), "moderate");

      // 送信
      fireEvent.click(screen.getByRole("button", { name: "登録する" }));

      await waitFor(() => {
        expect(onSuccess).toHaveBeenCalled();
      });
    });
  });

  describe("エラーリセット", () => {
    it("フィールド入力時にエラーがリセットされる", async () => {
      const error = new ApiError(
        "EMAIL_ALREADY_EXISTS",
        "Email already exists",
        409
      );
      mockUseRegisterUser.mockReturnValue({
        ...defaultHookReturn,
        error,
      });

      const user = userEvent.setup();
      render(<RegisterForm />);

      // エラーが表示されていることを確認
      expect(
        screen.getByText("このメールアドレスは既に登録されています")
      ).toBeInTheDocument();

      // フィールドに入力
      await user.type(screen.getByLabelText("ニックネーム"), "a");

      // reset関数が呼ばれることを確認
      expect(mockReset).toHaveBeenCalled();
    });
  });
});
