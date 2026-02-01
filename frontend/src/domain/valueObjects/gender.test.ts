import { describe, it, expect } from "vitest";
import { newGender } from "./gender";

describe("newGender", () => {
  describe("正常系", () => {
    const cases = [
      { name: "male", input: "male", expected: "male" },
      { name: "female", input: "female", expected: "female" },
      { name: "other", input: "other", expected: "other" },
    ];

    cases.forEach(({ name, input, expected }) => {
      it(name, () => {
        const result = newGender(input);
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
        expectedCode: "GENDER_INVALID",
        expectedMessage: "性別を選択してください",
      },
      {
        name: "unknown",
        input: "unknown",
        expectedCode: "GENDER_INVALID",
        expectedMessage: "性別を選択してください",
      },
      {
        name: "大文字のMALE",
        input: "MALE",
        expectedCode: "GENDER_INVALID",
        expectedMessage: "性別を選択してください",
      },
      {
        name: "日本語(男性)",
        input: "男性",
        expectedCode: "GENDER_INVALID",
        expectedMessage: "性別を選択してください",
      },
    ];

    cases.forEach(({ name, input, expectedCode, expectedMessage }) => {
      it(name, () => {
        const result = newGender(input);
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
        gender1: "male",
        gender2: "male",
        expected: true,
      },
      {
        name: "異なる値でfalseを返す",
        gender1: "male",
        gender2: "female",
        expected: false,
      },
    ];

    cases.forEach(({ name, gender1, gender2, expected }) => {
      it(name, () => {
        const r1 = newGender(gender1);
        const r2 = newGender(gender2);
        if (r1.ok && r2.ok) {
          expect(r1.value.equals(r2.value)).toBe(expected);
        }
      });
    });
  });
});
