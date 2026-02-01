import { describe, it, expect } from "vitest";
import { newItemName } from "./itemName";

describe("newItemName", () => {
  describe("正常系", () => {
    const cases = [
      { name: "通常の食品名", input: "りんご" },
      { name: "長い食品名", input: "特製手作りチョコレートケーキ" },
      { name: "英語の食品名", input: "Apple Pie" },
      { name: "数字を含む食品名", input: "カロリーメイト2本" },
    ];

    cases.forEach(({ name, input }) => {
      it(name, () => {
        const result = newItemName(input);
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
        name: "空文字",
        input: "",
        expectedCode: "ITEM_NAME_REQUIRED",
        expectedMessage: "食品名を入力してください",
      },
      {
        name: "空白のみ",
        input: "   ",
        expectedCode: "ITEM_NAME_REQUIRED",
        expectedMessage: "食品名を入力してください",
      },
      {
        name: "タブ文字のみ",
        input: "\t\t",
        expectedCode: "ITEM_NAME_REQUIRED",
        expectedMessage: "食品名を入力してください",
      },
    ];

    cases.forEach(({ name, input, expectedCode, expectedMessage }) => {
      it(name, () => {
        const result = newItemName(input);
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
        name1: "りんご",
        name2: "りんご",
        expected: true,
      },
      {
        name: "異なる値でfalseを返す",
        name1: "りんご",
        name2: "バナナ",
        expected: false,
      },
    ];

    cases.forEach(({ name, name1, name2, expected }) => {
      it(name, () => {
        const r1 = newItemName(name1);
        const r2 = newItemName(name2);
        if (r1.ok && r2.ok) {
          expect(r1.value.equals(r2.value)).toBe(expected);
        }
      });
    });
  });
});
