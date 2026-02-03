import { describe, it, expect, vi } from "vitest";
import { renderHook } from "@testing-library/react";
import { useRequestGet, useRequestMutation } from "./useRequest";

vi.mock("@/lib/swr", () => ({
  fetcher: vi.fn(),
  mutate: vi.fn(),
}));

describe("useRequest", () => {
  describe("useRequestGet", () => {
    it("urlがnullの場合はフェッチしない", () => {
      const { result } = renderHook(() => useRequestGet(null));
      expect(result.current.data).toBeUndefined();
      expect(result.current.isLoading).toBe(false);
    });
  });

  describe("useRequestMutation", () => {
    it("初期状態が正しく設定される", () => {
      const { result } = renderHook(() =>
        useRequestMutation<{ id: number }>("/api/test", "POST")
      );
      expect(result.current.data).toBeUndefined();
      expect(result.current.isMutating).toBe(false);
    });
  });
});
