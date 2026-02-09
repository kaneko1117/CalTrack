/**
 * useTodayPfc フックのテスト
 */
import { describe, it, expect, vi } from "vitest";
import { renderHook } from "@testing-library/react";
import { useTodayPfc, ERROR_MESSAGE_PFC_FETCH_FAILED } from "./useTodayPfc";

vi.mock("@/features/common/hooks", () => ({
  useRequestGet: vi.fn(() => {
    return { data: undefined, error: null, isLoading: true, isValidating: false, mutate: vi.fn() };
  }),
}));

describe("useTodayPfc", () => {
  it("初期状態ではローディング中となる", () => {
    const { result } = renderHook(() => useTodayPfc());
    expect(result.current.isLoading).toBe(true);
    expect(result.current.data).toBeUndefined();
    expect(result.current.error).toBeNull();
  });

  it("refetch関数が返される", () => {
    const { result } = renderHook(() => useTodayPfc());
    expect(typeof result.current.refetch).toBe("function");
  });

  it("ERROR_MESSAGE_PFC_FETCH_FAILEDが正しい", () => {
    expect(ERROR_MESSAGE_PFC_FETCH_FAILED).toBe("PFCデータの取得に失敗しました");
  });
});
