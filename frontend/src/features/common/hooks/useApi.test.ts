import { describe, it, expect, vi } from "vitest";
import { renderHook, act, waitFor } from "@testing-library/react";
import { useApi } from "./useApi";
import type { ApiErrorResponse } from "@/lib/api";

describe("useApi", () => {
  describe("初期状態", () => {
    it("初期状態が正しく設定される", () => {
      const apiFunction = vi.fn().mockResolvedValue({ id: 1 });
      const { result } = renderHook(() => useApi(apiFunction));

      expect(result.current.data).toBeNull();
      expect(result.current.error).toBeNull();
      expect(result.current.isPending).toBe(false);
      expect(result.current.isSuccess).toBe(false);
    });
  });

  describe("execute", () => {
    it("API成功時、dataがセットされisSuccessがtrueになる", async () => {
      const mockData = { id: 1, name: "Test" };
      const apiFunction = vi.fn().mockResolvedValue(mockData);
      const { result } = renderHook(() => useApi(apiFunction));

      act(() => {
        result.current.execute();
      });

      await waitFor(() => {
        expect(result.current.data).toEqual(mockData);
        expect(result.current.isSuccess).toBe(true);
        expect(result.current.error).toBeNull();
      });
    });

    it("API成功時、onSuccessコールバックが呼ばれる", async () => {
      const mockData = { id: 1 };
      const apiFunction = vi.fn().mockResolvedValue(mockData);
      const onSuccess = vi.fn();
      const { result } = renderHook(() => useApi(apiFunction, { onSuccess }));

      act(() => {
        result.current.execute();
      });

      await waitFor(() => {
        expect(onSuccess).toHaveBeenCalledWith(mockData);
      });
    });

    it("APIエラー時、errorがセットされる", async () => {
      const mockError: ApiErrorResponse = {
        code: "INTERNAL_ERROR",
        message: "サーバーエラー",
      };
      const apiFunction = vi.fn().mockRejectedValue(mockError);
      const { result } = renderHook(() => useApi(apiFunction));

      act(() => {
        result.current.execute();
      });

      await waitFor(() => {
        expect(result.current.error).toEqual(mockError);
        expect(result.current.isSuccess).toBe(false);
      });
    });

    it("APIエラー時、onErrorコールバックが呼ばれる", async () => {
      const mockError: ApiErrorResponse = {
        code: "INTERNAL_ERROR",
        message: "サーバーエラー",
      };
      const apiFunction = vi.fn().mockRejectedValue(mockError);
      const onError = vi.fn();
      const { result } = renderHook(() => useApi(apiFunction, { onError }));

      act(() => {
        result.current.execute();
      });

      await waitFor(() => {
        expect(onError).toHaveBeenCalledWith(mockError);
      });
    });
  });

  describe("reset", () => {
    it("状態を初期状態に戻す", async () => {
      const mockData = { id: 1 };
      const apiFunction = vi.fn().mockResolvedValue(mockData);
      const { result } = renderHook(() => useApi(apiFunction));

      act(() => {
        result.current.execute();
      });

      await waitFor(() => {
        expect(result.current.data).toEqual(mockData);
      });

      act(() => {
        result.current.reset();
      });

      expect(result.current.data).toBeNull();
      expect(result.current.error).toBeNull();
      expect(result.current.isSuccess).toBe(false);
    });
  });
});
