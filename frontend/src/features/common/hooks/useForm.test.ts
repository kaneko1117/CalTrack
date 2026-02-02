import { describe, it, expect, vi } from "vitest";
import { renderHook, act, waitFor } from "@testing-library/react";
import { useForm } from "./useForm";
import { ok, err } from "@/domain/shared/result";
import { domainError } from "@/domain/shared/errors";

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

describe("useForm", () => {
  describe("初期状態", () => {
    it("初期状態が正しく設定される", () => {
      const onSubmit = vi.fn();
      const { result } = renderHook(() =>
        useForm(
          { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          onSubmit
        )
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
      const onSubmit = vi.fn();
      const { result } = renderHook(() =>
        useForm(
          { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          onSubmit
        )
      );

      act(() => {
        result.current.handleChange("email")("test@example.com");
      });

      expect(result.current.formState.email).toBe("test@example.com");
      expect(result.current.errors.email).toBeNull();
    });

    it("バリデーション失敗時、値をセットしエラーメッセージをセットする", () => {
      const onSubmit = vi.fn();
      const { result } = renderHook(() =>
        useForm(
          { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          onSubmit
        )
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
      const onSubmit = vi.fn();
      const { result } = renderHook(() =>
        useForm(
          { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          onSubmit
        )
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
      const onSubmit = vi.fn();
      const { result } = renderHook(() =>
        useForm(
          { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          onSubmit
        )
      );

      act(() => {
        result.current.handleChange("email")("test@example.com");
        result.current.handleChange("password")("password123");
      });

      expect(result.current.isValid).toBe(true);
    });

    it("エラーがある場合falseを返す", () => {
      const onSubmit = vi.fn();
      const { result } = renderHook(() =>
        useForm(
          { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          onSubmit
        )
      );

      act(() => {
        result.current.handleChange("email")("invalid");
        result.current.handleChange("password")("password123");
      });

      expect(result.current.isValid).toBe(false);
    });

    it("未入力フィールドがある場合falseを返す", () => {
      const onSubmit = vi.fn();
      const { result } = renderHook(() =>
        useForm(
          { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          onSubmit
        )
      );

      act(() => {
        result.current.handleChange("email")("test@example.com");
      });

      expect(result.current.isValid).toBe(false);
    });
  });

  describe("handleSubmit", () => {
    it("isValidがtrueの場合、onSubmitが呼ばれる", async () => {
      const mockResult = { id: 1 };
      const onSubmit = vi.fn().mockResolvedValue(mockResult);
      const onSuccess = vi.fn();
      const { result } = renderHook(() =>
        useForm(
          { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          onSubmit,
          onSuccess
        )
      );

      act(() => {
        result.current.handleChange("email")("test@example.com");
        result.current.handleChange("password")("password123");
      });

      const mockEvent = {
        preventDefault: vi.fn(),
      } as unknown as React.FormEvent;

      act(() => {
        result.current.handleSubmit(mockEvent);
      });

      expect(mockEvent.preventDefault).toHaveBeenCalled();

      await waitFor(() => {
        expect(onSubmit).toHaveBeenCalledWith({
          email: "test@example.com",
          password: "password123",
        });
        expect(onSuccess).toHaveBeenCalledWith(mockResult);
      });
    });

    it("isValidがfalseの場合、onSubmitが呼ばれない", async () => {
      const onSubmit = vi.fn();
      const { result } = renderHook(() =>
        useForm(
          { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          onSubmit
        )
      );

      const mockEvent = {
        preventDefault: vi.fn(),
      } as unknown as React.FormEvent;

      act(() => {
        result.current.handleSubmit(mockEvent);
      });

      expect(mockEvent.preventDefault).toHaveBeenCalled();
      expect(onSubmit).not.toHaveBeenCalled();
    });

    it("APIエラー時、apiErrorがセットされる", async () => {
      const apiError = { code: "INVALID_CREDENTIALS", message: "認証エラー" };
      const onSubmit = vi.fn().mockRejectedValue(apiError);
      const { result } = renderHook(() =>
        useForm(
          { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          onSubmit
        )
      );

      act(() => {
        result.current.handleChange("email")("test@example.com");
        result.current.handleChange("password")("password123");
      });

      const mockEvent = {
        preventDefault: vi.fn(),
      } as unknown as React.FormEvent;

      act(() => {
        result.current.handleSubmit(mockEvent);
      });

      await waitFor(() => {
        expect(result.current.apiError).toEqual(apiError);
      });
    });
  });

  describe("reset", () => {
    it("フォームとエラーを初期状態に戻す", async () => {
      const apiError = { code: "ERROR", message: "エラー" };
      const onSubmit = vi.fn().mockRejectedValue(apiError);
      const { result } = renderHook(() =>
        useForm(
          { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          onSubmit
        )
      );

      act(() => {
        result.current.handleChange("email")("test@example.com");
        result.current.handleChange("password")("short");
      });

      // APIエラーをセット
      const mockEvent = {
        preventDefault: vi.fn(),
      } as unknown as React.FormEvent;
      act(() => {
        result.current.handleChange("password")("password123");
      });
      act(() => {
        result.current.handleSubmit(mockEvent);
      });

      await waitFor(() => {
        expect(result.current.apiError).toEqual(apiError);
      });

      act(() => {
        result.current.reset();
      });

      expect(result.current.formState).toEqual({ email: "", password: "" });
      expect(result.current.errors).toEqual({ email: null, password: null });
      expect(result.current.apiError).toBeNull();
    });
  });

  describe("handleChange時のapiErrorクリア", () => {
    it("入力時にapiErrorがクリアされる", async () => {
      const apiError = { code: "ERROR", message: "エラー" };
      const onSubmit = vi.fn().mockRejectedValue(apiError);
      const { result } = renderHook(() =>
        useForm(
          { email: mockEmailFactory, password: mockPasswordFactory },
          initialFormState,
          initialErrors,
          onSubmit
        )
      );

      act(() => {
        result.current.handleChange("email")("test@example.com");
        result.current.handleChange("password")("password123");
      });

      const mockEvent = {
        preventDefault: vi.fn(),
      } as unknown as React.FormEvent;

      act(() => {
        result.current.handleSubmit(mockEvent);
      });

      await waitFor(() => {
        expect(result.current.apiError).toEqual(apiError);
      });

      act(() => {
        result.current.handleChange("email")("new@example.com");
      });

      expect(result.current.apiError).toBeNull();
    });
  });
});
