import { describe, it, expect, vi, afterEach } from "vitest";
import { renderHook } from "@testing-library/react";
import { useCountUp } from "./useCountUp";

describe("useCountUp", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("初期値が0であること", () => {
    vi.spyOn(window, "requestAnimationFrame").mockReturnValue(1);

    const { result } = renderHook(() => useCountUp({ end: 100 }));
    expect(result.current).toBe(0);
  });

  it("startOnMountがfalseの場合、アニメーションが開始されないこと", () => {
    const rafSpy = vi.spyOn(window, "requestAnimationFrame").mockReturnValue(1);

    const { result } = renderHook(() =>
      useCountUp({ end: 100, startOnMount: false })
    );

    // requestAnimationFrameが呼ばれていない
    expect(rafSpy).not.toHaveBeenCalled();
    expect(result.current).toBe(0);
  });

  it("endが0の場合、0を返すこと", () => {
    vi.spyOn(window, "requestAnimationFrame").mockReturnValue(1);

    const { result } = renderHook(() => useCountUp({ end: 0 }));
    expect(result.current).toBe(0);
  });

  it("startOnMountがtrueの場合、requestAnimationFrameが呼ばれること", () => {
    const rafSpy = vi.spyOn(window, "requestAnimationFrame").mockReturnValue(1);

    renderHook(() => useCountUp({ end: 100, startOnMount: true }));

    expect(rafSpy).toHaveBeenCalled();
  });

  it("アンマウント時にcancelAnimationFrameが呼ばれること", () => {
    vi.spyOn(window, "requestAnimationFrame").mockReturnValue(123);
    const cancelSpy = vi.spyOn(window, "cancelAnimationFrame");

    const { unmount } = renderHook(() => useCountUp({ end: 100 }));
    unmount();

    expect(cancelSpy).toHaveBeenCalledWith(123);
  });

  it("durationを指定できること", () => {
    vi.spyOn(window, "requestAnimationFrame").mockReturnValue(1);

    const { result } = renderHook(() => useCountUp({ end: 100, duration: 500 }));
    // 初期状態は0
    expect(result.current).toBe(0);
  });

  it("endが変わると再アニメーションされること", () => {
    const rafSpy = vi.spyOn(window, "requestAnimationFrame").mockReturnValue(1);

    const { rerender } = renderHook(
      ({ end }) => useCountUp({ end, duration: 1000 }),
      { initialProps: { end: 100 } }
    );

    const initialCallCount = rafSpy.mock.calls.length;

    // endを変更
    rerender({ end: 200 });

    // 新しいアニメーションが開始される
    expect(rafSpy.mock.calls.length).toBeGreaterThan(initialCallCount);
  });
});
