import { describe, it, expect } from "vitest";
import { newWeight } from "./weight";

describe("newWeight", () => {
  describe("正常系", () => {
    const cases = [
      { name: "最小値(1kg)", input: "1", expected: 1 },
      { name: "最大値(500kg)", input: "500", expected: 500 },
      { name: "通常の値", input: "65", expected: 65 },
      { name: "小数点を含む値", input: "65.5", expected: 65.5 },
    ];

    cases.forEach(({ name, input, expected }) => {
      it(name, () => {
        const result = newWeight(input);
        expect(result.ok).toBe(true);
        if (result.ok) {
          expect(result.value.value).toBe(expected);
        }
      });
    });
  });

  describe("異常系", () => {
    const cases = [
      {
        name: "空文字",
        input: "",
        expectedCode: "WEIGHT_REQUIRED",
        expectedMessage: "体重を入力してください",
      },
      {
        name: "空白のみ",
        input: "   ",
        expectedCode: "WEIGHT_REQUIRED",
        expectedMessage: "体重を入力してください",
      },
      {
        name: "数値以外の文字列",
        input: "abc",
        expectedCode: "WEIGHT_INVALID",
        expectedMessage: "体重は有効な数値を入力してください",
      },
      {
        name: "0kg",
        input: "0",
        expectedCode: "WEIGHT_MUST_BE_POSITIVE",
        expectedMessage: "体重は0より大きい値を入力してください",
      },
      {
        name: "負の値(-1kg)",
        input: "-1",
        expectedCode: "WEIGHT_MUST_BE_POSITIVE",
        expectedMessage: "体重は0より大きい値を入力してください",
      },
      {
        name: "501kg超過",
        input: "501",
        expectedCode: "WEIGHT_TOO_HEAVY",
        expectedMessage: "体重は500kg以内で入力してください",
      },
      {
        name: "1000kg超過",
        input: "1000",
        expectedCode: "WEIGHT_TOO_HEAVY",
        expectedMessage: "体重は500kg以内で入力してください",
      },
    ];

    cases.forEach(({ name, input, expectedCode, expectedMessage }) => {
      it(name, () => {
        const result = newWeight(input);
        expect(result.ok).toBe(false);
        if (!result.ok) {
          expect(result.error.code).toBe(expectedCode);
          expect(result.error.message).toBe(expectedMessage);
        }
      });
    });
  });

  describe("equals", () => {
    const cases = [
      {
        name: "同じ値でtrueを返す",
        weight1: "65",
        weight2: "65",
        expected: true,
      },
      {
        name: "異なる値でfalseを返す",
        weight1: "65",
        weight2: "70",
        expected: false,
      },
    ];

    cases.forEach(({ name, weight1, weight2, expected }) => {
      it(name, () => {
        const r1 = newWeight(weight1);
        const r2 = newWeight(weight2);
        if (r1.ok && r2.ok) {
          expect(r1.value.equals(r2.value)).toBe(expected);
        }
      });
    });
  });
});
