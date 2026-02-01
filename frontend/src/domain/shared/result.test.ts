import { describe, it, expect } from "vitest";
import { ok, err } from "./result";

describe("Result", () => {
  describe("ok", () => {
    it("成功結果を生成できる", () => {
      const result = ok("test value");

      expect(result.ok).toBe(true);
      expect(result.value).toBe("test value");
    });

    it("数値で成功結果を生成できる", () => {
      const result = ok(42);

      expect(result.ok).toBe(true);
      expect(result.value).toBe(42);
    });

    it("オブジェクトで成功結果を生成できる", () => {
      const value = { id: 1, name: "test" };
      const result = ok(value);

      expect(result.ok).toBe(true);
      expect(result.value).toEqual(value);
    });

    it("結果オブジェクトは凍結されている", () => {
      const result = ok("test");

      expect(Object.isFrozen(result)).toBe(true);
    });
  });

  describe("err", () => {
    it("失敗結果を生成できる", () => {
      const error = { code: "TEST_ERROR", message: "テストエラー" };
      const result = err(error);

      expect(result.ok).toBe(false);
      expect(result.error).toEqual(error);
    });

    it("文字列でエラーを生成できる", () => {
      const result = err("エラーメッセージ");

      expect(result.ok).toBe(false);
      expect(result.error).toBe("エラーメッセージ");
    });

    it("結果オブジェクトは凍結されている", () => {
      const result = err("error");

      expect(Object.isFrozen(result)).toBe(true);
    });
  });
});
