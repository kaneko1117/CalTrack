import { describe, it, expect } from "vitest";
import { newCalories } from "./calories";

describe("newCalories", () => {
  describe("正常系", () => {
    const cases = [
      { name: "最小値(1)", input: 1 },
      { name: "通常の値(100)", input: 100 },
      { name: "大きな値(2000)", input: 2000 },
      { name: "整数の上限に近い値", input: 10000 },
    ];

    cases.forEach(({ name, input }) => {
      it(name, () => {
        const result = newCalories(input);
        expect(result.ok).toBe(true);
        if (result.ok) {
          expect(result.value.value).toBe(input);
        }
      });
    });
  });

  describe("異常系", () => {
    const cases = [
      {
        name: "0",
        input: 0,
        expectedCode: "CALORIES_MUST_BE_POSITIVE",
        expectedMessage: "カロリーは1以上の整数で入力してください",
      },
      {
        name: "負の値",
        input: -1,
        expectedCode: "CALORIES_MUST_BE_POSITIVE",
        expectedMessage: "カロリーは1以上の整数で入力してください",
      },
      {
        name: "大きな負の値",
        input: -100,
        expectedCode: "CALORIES_MUST_BE_POSITIVE",
        expectedMessage: "カロリーは1以上の整数で入力してください",
      },
    ];

    cases.forEach(({ name, input, expectedCode, expectedMessage }) => {
      it(name, () => {
        const result = newCalories(input);
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
        calories1: 100,
        calories2: 100,
        expected: true,
      },
      {
        name: "異なる値でfalseを返す",
        calories1: 100,
        calories2: 200,
        expected: false,
      },
    ];

    cases.forEach(({ name, calories1, calories2, expected }) => {
      it(name, () => {
        const r1 = newCalories(calories1);
        const r2 = newCalories(calories2);
        if (r1.ok && r2.ok) {
          expect(r1.value.equals(r2.value)).toBe(expected);
        }
      });
    });
  });
});
