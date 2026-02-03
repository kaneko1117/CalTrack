import { describe, it, expect } from "vitest";
import { newHeight } from "./height";

describe("newHeight", () => {
  describe("正常系", () => {
    const cases = [
      { name: "最小値(1cm)", input: "1", expected: 1 },
      { name: "最大値(300cm)", input: "300", expected: 300 },
      { name: "通常の値", input: "170", expected: 170 },
      { name: "小数点を含む値", input: "170.5", expected: 170.5 },
    ];

    cases.forEach(({ name, input, expected }) => {
      it(name, () => {
        const result = newHeight(input);
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
        expectedCode: "HEIGHT_REQUIRED",
        expectedMessage: "身長を入力してください",
      },
      {
        name: "空白のみ",
        input: "   ",
        expectedCode: "HEIGHT_REQUIRED",
        expectedMessage: "身長を入力してください",
      },
      {
        name: "数値以外の文字列",
        input: "abc",
        expectedCode: "HEIGHT_INVALID",
        expectedMessage: "身長は有効な数値を入力してください",
      },
      {
        name: "0cm",
        input: "0",
        expectedCode: "HEIGHT_MUST_BE_POSITIVE",
        expectedMessage: "身長は0より大きい値を入力してください",
      },
      {
        name: "負の値(-1cm)",
        input: "-1",
        expectedCode: "HEIGHT_MUST_BE_POSITIVE",
        expectedMessage: "身長は0より大きい値を入力してください",
      },
      {
        name: "301cm超過",
        input: "301",
        expectedCode: "HEIGHT_TOO_TALL",
        expectedMessage: "身長は300cm以内で入力してください",
      },
      {
        name: "500cm超過",
        input: "500",
        expectedCode: "HEIGHT_TOO_TALL",
        expectedMessage: "身長は300cm以内で入力してください",
      },
    ];

    cases.forEach(({ name, input, expectedCode, expectedMessage }) => {
      it(name, () => {
        const result = newHeight(input);
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
        height1: "170",
        height2: "170",
        expected: true,
      },
      {
        name: "異なる値でfalseを返す",
        height1: "170",
        height2: "180",
        expected: false,
      },
    ];

    cases.forEach(({ name, height1, height2, expected }) => {
      it(name, () => {
        const r1 = newHeight(height1);
        const r2 = newHeight(height2);
        if (r1.ok && r2.ok) {
          expect(r1.value.equals(r2.value)).toBe(expected);
        }
      });
    });
  });
});
