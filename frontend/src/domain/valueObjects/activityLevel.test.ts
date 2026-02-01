import { describe, it, expect } from "vitest";
import { newActivityLevel } from "./activityLevel";

describe("newActivityLevel", () => {
  describe("正常系", () => {
    const cases = [
      { name: "sedentary", input: "sedentary", expected: "sedentary" },
      { name: "light", input: "light", expected: "light" },
      { name: "moderate", input: "moderate", expected: "moderate" },
      { name: "active", input: "active", expected: "active" },
      { name: "veryActive", input: "veryActive", expected: "veryActive" },
    ];

    cases.forEach(({ name, input, expected }) => {
      it(name, () => {
        const result = newActivityLevel(input);
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
        expectedCode: "ACTIVITY_LEVEL_INVALID",
        expectedMessage: "活動レベルを選択してください",
      },
      {
        name: "invalid",
        input: "invalid",
        expectedCode: "ACTIVITY_LEVEL_INVALID",
        expectedMessage: "活動レベルを選択してください",
      },
      {
        name: "大文字のSEDENTARY",
        input: "SEDENTARY",
        expectedCode: "ACTIVITY_LEVEL_INVALID",
        expectedMessage: "活動レベルを選択してください",
      },
      {
        name: "スペースを含む値",
        input: "very active",
        expectedCode: "ACTIVITY_LEVEL_INVALID",
        expectedMessage: "活動レベルを選択してください",
      },
      {
        name: "日本語(活発)",
        input: "活発",
        expectedCode: "ACTIVITY_LEVEL_INVALID",
        expectedMessage: "活動レベルを選択してください",
      },
    ];

    cases.forEach(({ name, input, expectedCode, expectedMessage }) => {
      it(name, () => {
        const result = newActivityLevel(input);
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
        level1: "moderate",
        level2: "moderate",
        expected: true,
      },
      {
        name: "異なる値でfalseを返す",
        level1: "sedentary",
        level2: "active",
        expected: false,
      },
    ];

    cases.forEach(({ name, level1, level2, expected }) => {
      it(name, () => {
        const r1 = newActivityLevel(level1);
        const r2 = newActivityLevel(level2);
        if (r1.ok && r2.ok) {
          expect(r1.value.equals(r2.value)).toBe(expected);
        }
      });
    });
  });
});
