import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useForm } from "./useForm";
import { ok, err } from "@/domain/shared/result";
import { domainError } from "@/domain/shared/errors";
import type { ApiErrorResponse } from "@/lib/api";

// SWR mutationをモック
const mockTrigger = vi.fn();
const mockReset = vi.fn();
let mockError: ApiErrorResponse | undefined = undefined;

vi.mock("./useRequest", () => ({
  useRequestMutation: () => ({
    trigger: mockTrigger,
    isMutating: false,
    error: mockError,
    data: undefined,
    reset: mockReset,
  }),
}));

const mockEmailFactory = vi.fn((value: string) => {
  if (!value || value.trim() === "") {
    return err(
      domainError("EMAIL_REQUIRED", "メールアドレスを入力してください")
    );
  }
  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) {
    return err(
      domainError("EMAIL_INVALID_FORMAT", "有効なメールアドレスを入力してください")
    );
  }
  return ok({ value });
});

const mockPasswordFactory = vi.fn((value: string) => {
  if (!value || value.trim() === "") {
    return err(
      domainError("PASSWORD_REQUIRED", "パスワードを入力してください")
    );
  }
  if (value.length < 8) {
    return err(
      domainError("PASSWORD_TOO_SHORT", "パスワードは8文字以上で入力してください")
    );
  }
  return ok({ value });
});

const initialFormState = { email: "", password: "" };
const initialErrors = { email: null, password: null };

type TestField = "email" | "password";
type TestResponse = { id: number };
type TestRequest = { email: string; password: string };

describe("useForm", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockError = undefined;
  });

  describe("初期状態", () => {
    it("初期状態が正しく設定される", () => {
      const { result } = renderHook(() =>
        useForm<TestField, TestResponse, TestRequest>({
          config: { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          url: "/api/test",
          transformData: (data) => ({ email: data.email, password: data.password }),
        })
      );

      expect(result.current.formState).toEqual({ email: "", password: "" });
      expect(result.current.errors).toEqual({ email: null, password: null });
      expect(result.current.apiError).toBeNull();
      expect(result.current.isValid).toBe(false);
      expect(result.current.isPending).toBe(false);
    });
  });

  describe("handleChange", () => {
    it("バリデーション成功時、値をセットしエラーをnullにする", () => {
      const { result } = renderHook(() =>
        useForm<TestField, TestResponse, TestRequest>({
          config: { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          url: "/api/test",
          transformData: (data) => ({ email: data.email, password: data.password }),
        })
      );

      act(() => {
        result.current.handleChange("email")("test@example.com");
      });

      expect(result.current.formState.email).toBe("test@example.com");
      expect(result.current.errors.email).toBeNull();
    });

    it("バリデーション失敗時、値をセットしエラーメッセージをセットする", () => {
      const { result } = renderHook(() =>
        useForm<TestField, TestResponse, TestRequest>({
          config: { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          url: "/api/test",
          transformData: (data) => ({ email: data.email, password: data.password }),
        })
      );

      act(() => {
        result.current.handleChange("email")("invalid-email");
      });

      expect(result.current.formState.email).toBe("invalid-email");
      expect(result.current.errors.email).toBe(
        "有効なメールアドレスを入力してください"
      );
    });

    it("指定したフィールドのみ更新される", () => {
      const { result } = renderHook(() =>
        useForm<TestField, TestResponse, TestRequest>({
          config: { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          url: "/api/test",
          transformData: (data) => ({ email: data.email, password: data.password }),
        })
      );

      act(() => {
        result.current.handleChange("email")("test@example.com");
      });

      expect(result.current.formState.email).toBe("test@example.com");
      expect(result.current.formState.password).toBe("");
    });
  });

  describe("isValid", () => {
    it("全フィールドが有効かつ入力済みでtrueを返す", () => {
      const { result } = renderHook(() =>
        useForm<TestField, TestResponse, TestRequest>({
          config: { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          url: "/api/test",
          transformData: (data) => ({ email: data.email, password: data.password }),
        })
      );

      act(() => {
        result.current.handleChange("email")("test@example.com");
        result.current.handleChange("password")("password123");
      });

      expect(result.current.isValid).toBe(true);
    });

    it("エラーがある場合falseを返す", () => {
      const { result } = renderHook(() =>
        useForm<TestField, TestResponse, TestRequest>({
          config: { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          url: "/api/test",
          transformData: (data) => ({ email: data.email, password: data.password }),
        })
      );

      act(() => {
        result.current.handleChange("email")("invalid");
        result.current.handleChange("password")("password123");
      });

      expect(result.current.isValid).toBe(false);
    });

    it("未入力フィールドがある場合falseを返す", () => {
      const { result } = renderHook(() =>
        useForm<TestField, TestResponse, TestRequest>({
          config: { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          url: "/api/test",
          transformData: (data) => ({ email: data.email, password: data.password }),
        })
      );

      act(() => {
        result.current.handleChange("email")("test@example.com");
      });

      expect(result.current.isValid).toBe(false);
    });
  });

  describe("handleSubmit", () => {
    it("isValidがtrueの場合、triggerが呼ばれる", async () => {
      mockTrigger.mockResolvedValue({ id: 1 });
      const { result } = renderHook(() =>
        useForm<TestField, TestResponse, TestRequest>({
          config: { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          url: "/api/test",
          transformData: (data) => ({ email: data.email, password: data.password }),
        })
      );

      act(() => {
        result.current.handleChange("email")("test@example.com");
        result.current.handleChange("password")("password123");
      });

      const mockEvent = {
        preventDefault: vi.fn(),
      } as unknown as React.FormEvent;

      await act(async () => {
        await result.current.handleSubmit(mockEvent);
      });

      expect(mockEvent.preventDefault).toHaveBeenCalled();
      expect(mockTrigger).toHaveBeenCalledWith({
        email: "test@example.com",
        password: "password123",
      });
    });

    it("isValidがfalseの場合、triggerが呼ばれない", async () => {
      const { result } = renderHook(() =>
        useForm<TestField, TestResponse, TestRequest>({
          config: { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          url: "/api/test",
          transformData: (data) => ({ email: data.email, password: data.password }),
        })
      );

      const mockEvent = {
        preventDefault: vi.fn(),
      } as unknown as React.FormEvent;

      await act(async () => {
        await result.current.handleSubmit(mockEvent);
      });

      expect(mockEvent.preventDefault).toHaveBeenCalled();
      expect(mockTrigger).not.toHaveBeenCalled();
    });

    it("APIエラー時、エラーがスローされてもクラッシュしない", async () => {
      const apiError = { code: "INVALID_CREDENTIALS", message: "認証エラー" };
      mockTrigger.mockRejectedValue(apiError);
      const { result } = renderHook(() =>
        useForm<TestField, TestResponse, TestRequest>({
          config: { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          url: "/api/test",
          transformData: (data) => ({ email: data.email, password: data.password }),
        })
      );

      act(() => {
        result.current.handleChange("email")("test@example.com");
        result.current.handleChange("password")("password123");
      });

      const mockEvent = {
        preventDefault: vi.fn(),
      } as unknown as React.FormEvent;

      // エラーがスローされてもクラッシュしないことを確認
      await act(async () => {
        await result.current.handleSubmit(mockEvent);
      });

      expect(mockTrigger).toHaveBeenCalled();
    });
  });

  describe("reset", () => {
    it("フォームとエラーを初期状態に戻す", async () => {
      const { result } = renderHook(() =>
        useForm<TestField, TestResponse, TestRequest>({
          config: { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          url: "/api/test",
          transformData: (data) => ({ email: data.email, password: data.password }),
        })
      );

      act(() => {
        result.current.handleChange("email")("test@example.com");
        result.current.handleChange("password")("password123");
      });

      act(() => {
        result.current.reset();
      });

      expect(result.current.formState).toEqual({ email: "", password: "" });
      expect(result.current.errors).toEqual({ email: null, password: null });
      expect(mockReset).toHaveBeenCalled();
    });
  });

  describe("handleChange時のapiErrorクリア", () => {
    it("入力時にapiErrorがある場合、resetが呼ばれる", async () => {
      // エラーがある状態でモック
      mockError = { code: "INVALID_CREDENTIALS", message: "認証エラー" };

      const { result } = renderHook(() =>
        useForm<TestField, TestResponse, TestRequest>({
          config: { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          url: "/api/test",
          transformData: (data) => ({ email: data.email, password: data.password }),
        })
      );

      act(() => {
        result.current.handleChange("email")("new@example.com");
      });

      expect(mockReset).toHaveBeenCalled();
    });
  });
});
