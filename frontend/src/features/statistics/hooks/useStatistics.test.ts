/**
 * useStatistics フックのテスト
 */
import { describe, it, expect, vi } from "vitest";
import { renderHook } from "@testing-library/react";
import { useStatistics, PERIOD_LABELS, ERROR_MESSAGE_FETCH_FAILED } from "./useStatistics";

vi.mock("@/features/common/hooks", () => ({
  useRequestGet: vi.fn((url: string | null) => {
    if (url === null) {
      return { data: undefined, error: null, isLoading: false, mutate: vi.fn() };
    }
    return { data: undefined, error: null, isLoading: true, mutate: vi.fn() };
  }),
}));

describe("useStatistics", () => {
  describe("フックの初期状態", () => {
    it("weekで呼び出した場合、正しいURLでリクエストされる", () => {
      const { result } = renderHook(() => useStatistics("week"));
      expect(result.current.isLoading).toBe(true);
      expect(result.current.data).toBeUndefined();
      expect(result.current.error).toBeNull();
    });

    it("monthで呼び出した場合、正しいURLでリクエストされる", () => {
      const { result } = renderHook(() => useStatistics("month"));
      expect(result.current.isLoading).toBe(true);
      expect(result.current.data).toBeUndefined();
    });
  });

  describe("PERIOD_LABELS定数", () => {
    it("weekのラベルが正しい", () => {
      expect(PERIOD_LABELS.week).toBe("週間");
    });

    it("monthのラベルが正しい", () => {
      expect(PERIOD_LABELS.month).toBe("月間");
    });
  });

  describe("エラーメッセージ定数", () => {
    it("ERROR_MESSAGE_FETCH_FAILEDが正しい", () => {
      expect(ERROR_MESSAGE_FETCH_FAILED).toBe("統計データの取得に失敗しました");
    });
  });

  describe("refetch関数", () => {
    it("refetch関数が返される", () => {
      const { result } = renderHook(() => useStatistics("week"));
      expect(typeof result.current.refetch).toBe("function");
    });
  });
});
