import { describe, it, expect } from "vitest";
import { domainError } from "./errors";

describe("domainError", () => {
  it("エラーコードとメッセージからDomainErrorを生成できる", () => {
    const error = domainError("TEST_ERROR", "テストエラーメッセージ");

    expect(error.code).toBe("TEST_ERROR");
    expect(error.message).toBe("テストエラーメッセージ");
  });

  it("異なるエラーコードで生成できる", () => {
    const error = domainError("VALIDATION_ERROR", "バリデーションエラー");

    expect(error.code).toBe("VALIDATION_ERROR");
    expect(error.message).toBe("バリデーションエラー");
  });

  it("生成されたエラーオブジェクトは凍結されている", () => {
    const error = domainError("TEST_ERROR", "テストエラー");

    expect(Object.isFrozen(error)).toBe(true);
  });
});
