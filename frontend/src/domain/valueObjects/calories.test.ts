import { describe, it, expect } from "vitest";
import { newCalories, sumCalories } from "./calories";

describe("newCalories", () => {
  describe("正常系", () => {
    const cases = [
      { name: "最小値(1)", input: "1", expected: 1 },
      { name: "通常の値(100)", input: "100", expected: 100 },
      { name: "大きな値(2000)", input: "2000", expected: 2000 },
      { name: "整数の上限に近い値", input: "10000", expected: 10000 },
    ];

    cases.forEach(({ name, input, expected }) => {
      it(name, () => {
        const result = newCalories(input);
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
        expectedCode: "CALORIES_REQUIRED",
        expectedMessage: "カロリーを入力してください",
      },
      {
        name: "空白のみ",
        input: "   ",
        expectedCode: "CALORIES_REQUIRED",
        expectedMessage: "カロリーを入力してください",
      },
      {
        name: "数値以外の文字列",
        input: "abc",
        expectedCode: "CALORIES_INVALID",
        expectedMessage: "カロリーは有効な数値を入力してください",
      },
      {
        name: "0",
        input: "0",
        expectedCode: "CALORIES_MUST_BE_POSITIVE",
        expectedMessage: "カロリーは1以上の整数で入力してください",
      },
      {
        name: "負の値",
        input: "-1",
        expectedCode: "CALORIES_MUST_BE_POSITIVE",
        expectedMessage: "カロリーは1以上の整数で入力してください",
      },
      {
        name: "大きな負の値",
        input: "-100",
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
        calories1: "100",
        calories2: "100",
        expected: true,
      },
      {
        name: "異なる値でfalseを返す",
        calories1: "100",
        calories2: "200",
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

describe("sumCalories", () => {
  const cases = [
    {
      name: "空配列で0を返す",
      input: [],
      expected: 0,
    },
    {
      name: "単一要素の合計",
      input: [100],
      expected: 100,
    },
    {
      name: "複数要素の合計",
      input: [100, 200, 300],
      expected: 600,
    },
    {
      name: "NaN値は0として扱う",
      input: [100, NaN, 200],
      expected: 300,
    },
    {
      name: "全てNaNの場合は0を返す",
      input: [NaN, NaN, NaN],
      expected: 0,
    },
  ];

  cases.forEach(({ name, input, expected }) => {
    it(name, () => {
      const result = sumCalories(input);
      expect(result).toBe(expected);
    });
  });
});
