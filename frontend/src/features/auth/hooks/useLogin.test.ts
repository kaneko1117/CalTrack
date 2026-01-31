import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { renderHook, act, waitFor } from "@testing-library/react";
import { useLogin } from "./useLogin";
import * as authApi from "../api";
import type { LoginRequest, LoginResponse } from "../types";
import {
  ERROR_CODE_INTERNAL_ERROR,
  ERROR_CODE_INVALID_CREDENTIALS,
} from "../types";

// APIモジュールをモック
vi.mock("../api", async () => {
  const actual = await vi.importActual<typeof authApi>("../api");
  return {
    ...actual,
    login: vi.fn(),
  };
});

const mockLogin = vi.mocked(authApi.login);

describe("useLogin", () => {
  const validRequest: LoginRequest = {
    email: "test@example.com",
    password: "password123",
  };

  const successResponse: LoginResponse = {
    userId: "user-123",
    email: "test@example.com",
    nickname: "TestUser",
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.resetAllMocks();
  });

  describe("初期状態", () => {
    it("正しい初期状態を持つこと", () => {
      const { result } = renderHook(() => useLogin());

      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeNull();
      expect(result.current.isSuccess).toBe(false);
      expect(typeof result.current.login).toBe("function");
      expect(typeof result.current.reset).toBe("function");
    });
  });

  describe("ログイン成功", () => {
    it("ログイン成功時にisSuccessがtrueになること", async () => {
      mockLogin.mockResolvedValueOnce(successResponse);

      const { result } = renderHook(() => useLogin());

      await act(async () => {
        await result.current.login(validRequest);
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
        expect(result.current.isLoading).toBe(false);
        expect(result.current.error).toBeNull();
      });

      expect(mockLogin).toHaveBeenCalledWith(validRequest);
      expect(mockLogin).toHaveBeenCalledTimes(1);
    });

    it("ログイン成功時にonSuccessコールバックが呼ばれること", async () => {
      mockLogin.mockResolvedValueOnce(successResponse);
      const onSuccess = vi.fn();

      const { result } = renderHook(() => useLogin());

      await act(async () => {
        await result.current.login(validRequest, onSuccess);
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(onSuccess).toHaveBeenCalledWith(successResponse);
      expect(onSuccess).toHaveBeenCalledTimes(1);
    });

    it("ログイン中はisLoadingがtrueになること", async () => {
      let resolvePromise: (value: LoginResponse) => void;
      const pendingPromise = new Promise<LoginResponse>((resolve) => {
        resolvePromise = resolve;
      });
      mockLogin.mockReturnValueOnce(pendingPromise);

      const { result } = renderHook(() => useLogin());

      act(() => {
        result.current.login(validRequest);
      });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(true);
      });

      await act(async () => {
        resolvePromise!(successResponse);
      });

      await waitFor(() => {
        expect(result.current.isLoading).toBe(false);
      });
    });
  });

  describe("認証エラー", () => {
    it("認証エラー時にINVALID_CREDENTIALSエラーが設定されること", async () => {
      const authError = new authApi.ApiError(
        ERROR_CODE_INVALID_CREDENTIALS,
        "メールアドレスまたはパスワードが間違っています",
        401
      );
      mockLogin.mockRejectedValueOnce(authError);

      const { result } = renderHook(() => useLogin());

      await act(async () => {
        await result.current.login(validRequest);
      });

      await waitFor(() => {
        expect(result.current.error).not.toBeNull();
        expect(result.current.error?.code).toBe(ERROR_CODE_INVALID_CREDENTIALS);
        expect(result.current.isLoading).toBe(false);
        expect(result.current.isSuccess).toBe(false);
      });
    });

    it("認証エラー時にonSuccessコールバックが呼ばれないこと", async () => {
      const authError = new authApi.ApiError(
        ERROR_CODE_INVALID_CREDENTIALS,
        "メールアドレスまたはパスワードが間違っています",
        401
      );
      mockLogin.mockRejectedValueOnce(authError);
      const onSuccess = vi.fn();

      const { result } = renderHook(() => useLogin());

      await act(async () => {
        await result.current.login(validRequest, onSuccess);
      });

      await waitFor(() => {
        expect(result.current.error).not.toBeNull();
      });

      expect(onSuccess).not.toHaveBeenCalled();
    });
  });

  describe("サーバーエラー", () => {
    it("サーバーエラー時にエラーが設定されること", async () => {
      const serverError = new authApi.ApiError(
        ERROR_CODE_INTERNAL_ERROR,
        "Internal server error",
        500
      );
      mockLogin.mockRejectedValueOnce(serverError);

      const { result } = renderHook(() => useLogin());

      await act(async () => {
        await result.current.login(validRequest);
      });

      await waitFor(() => {
        expect(result.current.error).not.toBeNull();
        expect(result.current.error?.code).toBe(ERROR_CODE_INTERNAL_ERROR);
        expect(result.current.isLoading).toBe(false);
        expect(result.current.isSuccess).toBe(false);
      });
    });

    it("予期しないエラー時にINTERNAL_ERRORが設定されること", async () => {
      mockLogin.mockRejectedValueOnce(new Error("Network error"));

      const { result } = renderHook(() => useLogin());

      await act(async () => {
        await result.current.login(validRequest);
      });

      await waitFor(() => {
        expect(result.current.error).not.toBeNull();
        expect(result.current.error?.code).toBe(ERROR_CODE_INTERNAL_ERROR);
        expect(result.current.isLoading).toBe(false);
        expect(result.current.isSuccess).toBe(false);
      });
    });
  });

  describe("reset", () => {
    it("reset()で全てのステートがリセットされること", async () => {
      const authError = new authApi.ApiError(
        ERROR_CODE_INVALID_CREDENTIALS,
        "認証エラー",
        401
      );
      mockLogin.mockRejectedValueOnce(authError);

      const { result } = renderHook(() => useLogin());

      await act(async () => {
        await result.current.login(validRequest);
      });

      await waitFor(() => {
        expect(result.current.error).not.toBeNull();
      });

      act(() => {
        result.current.reset();
      });

      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeNull();
      expect(result.current.isSuccess).toBe(false);
    });

    it("成功状態からreset()でステートがリセットされること", async () => {
      mockLogin.mockResolvedValueOnce(successResponse);

      const { result } = renderHook(() => useLogin());

      await act(async () => {
        await result.current.login(validRequest);
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      act(() => {
        result.current.reset();
      });

      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeNull();
      expect(result.current.isSuccess).toBe(false);
    });
  });
});
