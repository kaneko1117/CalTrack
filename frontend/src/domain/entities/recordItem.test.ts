import { describe, it, expect } from "vitest";
import { newRecordItem } from "./recordItem";

describe("newRecordItem", () => {
  describe("正常系", () => {
    const cases = [
      {
        name: "通常の食品アイテム",
        input: { name: "りんご", calories: "100" },
        expected: 100,
      },
      {
        name: "カロリー最小値",
        input: { name: "水", calories: "1" },
        expected: 1,
      },
      {
        name: "高カロリー食品",
        input: { name: "チョコレートケーキ", calories: "500" },
        expected: 500,
      },
    ];

    cases.forEach(({ name, input, expected }) => {
      it(name, () => {
        const result = newRecordItem(input);
        expect(result.ok).toBe(true);
        if (result.ok) {
          expect(result.value.name.value).toBe(input.name);
          expect(result.value.calories.value).toBe(expected);
        }
      });
    });
  });

  describe("異常系 - 単一フィールドエラー", () => {
    const cases = [
      {
        name: "食品名が空",
        input: { name: "", calories: "100" },
        expectedErrors: {
          name: { code: "ITEM_NAME_REQUIRED" },
        },
      },
      {
        name: "カロリーが0",
        input: { name: "りんご", calories: "0" },
        expectedErrors: {
          calories: { code: "CALORIES_MUST_BE_POSITIVE" },
        },
      },
      {
        name: "カロリーが負の値",
        input: { name: "りんご", calories: "-10" },
        expectedErrors: {
          calories: { code: "CALORIES_MUST_BE_POSITIVE" },
        },
      },
    ];

    cases.forEach(({ name, input, expectedErrors }) => {
      it(name, () => {
        const result = newRecordItem(input);
        expect(result.ok).toBe(false);
        if (!result.ok) {
          if (expectedErrors.name) {
            expect(result.error.name).toBeDefined();
            expect(result.error.name?.code).toBe(expectedErrors.name.code);
          }
          if (expectedErrors.calories) {
            expect(result.error.calories).toBeDefined();
            expect(result.error.calories?.code).toBe(expectedErrors.calories.code);
          }
        }
      });
    });
  });

  describe("異常系 - 複数フィールドエラー", () => {
    it("食品名が空かつカロリーが0", () => {
      const result = newRecordItem({ name: "", calories: "0" });
      expect(result.ok).toBe(false);
      if (!result.ok) {
        expect(result.error.name).toBeDefined();
        expect(result.error.name?.code).toBe("ITEM_NAME_REQUIRED");
        expect(result.error.calories).toBeDefined();
        expect(result.error.calories?.code).toBe("CALORIES_MUST_BE_POSITIVE");
      }
    });

    it("食品名が空白のみかつカロリーが負の値", () => {
      const result = newRecordItem({ name: "   ", calories: "-5" });
      expect(result.ok).toBe(false);
      if (!result.ok) {
        expect(result.error.name).toBeDefined();
        expect(result.error.name?.code).toBe("ITEM_NAME_REQUIRED");
        expect(result.error.calories).toBeDefined();
        expect(result.error.calories?.code).toBe("CALORIES_MUST_BE_POSITIVE");
      }
    });
  });
});
