import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { renderHook, act, waitFor } from "@testing-library/react";
import { useLogout } from "./useLogout";
import * as authApi from "../api";
import type { LogoutResponse } from "../types";
import { ERROR_CODE_INTERNAL_ERROR } from "../types";

// APIモジュールをモック
vi.mock("../api", async () => {
  const actual = await vi.importActual<typeof authApi>("../api");
  return {
    ...actual,
    logout: vi.fn(),
  };
});

const mockLogout = vi.mocked(authApi.logout);

describe("useLogout", () => {
  const successResponse: LogoutResponse = {
    message: "ログアウトしました",
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.resetAllMocks();
  });

  describe("初期状態", () => {
    it("正しい初期状態を持つこと", () => {
      const { result } = renderHook(() => useLogout());

      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeNull();
      expect(result.current.isSuccess).toBe(false);
      expect(typeof result.current.logout).toBe("function");
      expect(typeof result.current.reset).toBe("function");
    });
  });

  describe("ログアウト成功", () => {
    it("ログアウト成功時にisSuccessがtrueになること", async () => {
      mockLogout.mockResolvedValueOnce(successResponse);

      const { result } = renderHook(() => useLogout());

      await act(async () => {
        await result.current.logout();
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
        expect(result.current.isLoading).toBe(false);
        expect(result.current.error).toBeNull();
      });

      expect(mockLogout).toHaveBeenCalledTimes(1);
    });

    it("ログアウト成功時にonSuccessコールバックが呼ばれること", async () => {
      mockLogout.mockResolvedValueOnce(successResponse);
      const onSuccess = vi.fn();

      const { result } = renderHook(() => useLogout());

      await act(async () => {
        await result.current.logout(onSuccess);
      });

      await waitFor(() => {
        expect(result.current.isSuccess).toBe(true);
      });

      expect(onSuccess).toHaveBeenCalledTimes(1);
    });

    it("ログアウト中はisLoadingがtrueになること", async () => {
      let resolvePromise: (value: LogoutResponse) => void;
      const pendingPromise = new Promise<LogoutResponse>((resolve) => {
        resolvePromise = resolve;
      });
      mockLogout.mockReturnValueOnce(pendingPromise);

      const { result } = renderHook(() => useLogout());

      act(() => {
        result.current.logout();
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

  describe("ログアウト失敗", () => {
    it("サーバーエラー時にエラーが設定されること", async () => {
      const serverError = new authApi.ApiError(
        ERROR_CODE_INTERNAL_ERROR,
        "Internal server error",
        500
      );
      mockLogout.mockRejectedValueOnce(serverError);

      const { result } = renderHook(() => useLogout());

      await act(async () => {
        await result.current.logout();
      });

      await waitFor(() => {
        expect(result.current.error).not.toBeNull();
        expect(result.current.error?.code).toBe(ERROR_CODE_INTERNAL_ERROR);
        expect(result.current.isLoading).toBe(false);
        expect(result.current.isSuccess).toBe(false);
      });
    });

    it("失敗時にonSuccessコールバックが呼ばれないこと", async () => {
      const serverError = new authApi.ApiError(
        ERROR_CODE_INTERNAL_ERROR,
        "Internal server error",
        500
      );
      mockLogout.mockRejectedValueOnce(serverError);
      const onSuccess = vi.fn();

      const { result } = renderHook(() => useLogout());

      await act(async () => {
        await result.current.logout(onSuccess);
      });

      await waitFor(() => {
        expect(result.current.error).not.toBeNull();
      });

      expect(onSuccess).not.toHaveBeenCalled();
    });

    it("予期しないエラー時にINTERNAL_ERRORが設定されること", async () => {
      mockLogout.mockRejectedValueOnce(new Error("Network error"));

      const { result } = renderHook(() => useLogout());

      await act(async () => {
        await result.current.logout();
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
      const serverError = new authApi.ApiError(
        ERROR_CODE_INTERNAL_ERROR,
        "サーバーエラー",
        500
      );
      mockLogout.mockRejectedValueOnce(serverError);

      const { result } = renderHook(() => useLogout());

      await act(async () => {
        await result.current.logout();
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
      mockLogout.mockResolvedValueOnce(successResponse);

      const { result } = renderHook(() => useLogout());

      await act(async () => {
        await result.current.logout();
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
