import { describe, it, expect } from "vitest";
import { newEmail } from "./email";

describe("newEmail", () => {
  describe("正常系", () => {
    const cases = [
      { name: "通常のメールアドレス", input: "test@example.com" },
      { name: "サブドメイン付き", input: "test@sub.example.com" },
      { name: "最大長(254文字)", input: "a".repeat(242) + "@example.com" },
    ];

    cases.forEach(({ name, input }) => {
      it(name, () => {
        const result = newEmail(input);
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
        expectedCode: "EMAIL_REQUIRED",
        expectedMessage: "メールアドレスを入力してください",
      },
      {
        name: "空白のみ",
        input: "   ",
        expectedCode: "EMAIL_REQUIRED",
        expectedMessage: "メールアドレスを入力してください",
      },
      {
        name: "255文字超過",
        input: "a".repeat(244) + "@example.com",
        expectedCode: "EMAIL_TOO_LONG",
        expectedMessage: "メールアドレスは254文字以内で入力してください",
      },
      {
        name: "@なし",
        input: "testexample.com",
        expectedCode: "EMAIL_INVALID_FORMAT",
        expectedMessage: "有効なメールアドレスを入力してください",
      },
      {
        name: "ドットなしのドメイン",
        input: "test@example",
        expectedCode: "EMAIL_INVALID_FORMAT",
        expectedMessage: "有効なメールアドレスを入力してください",
      },
      {
        name: "ローカル部なし",
        input: "@example.com",
        expectedCode: "EMAIL_INVALID_FORMAT",
        expectedMessage: "有効なメールアドレスを入力してください",
      },
    ];

    cases.forEach(({ name, input, expectedCode, expectedMessage }) => {
      it(name, () => {
        const result = newEmail(input);
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
        email1: "test@example.com",
        email2: "test@example.com",
        expected: true,
      },
      {
        name: "異なる値でfalseを返す",
        email1: "test1@example.com",
        email2: "test2@example.com",
        expected: false,
      },
    ];

    cases.forEach(({ name, email1, email2, expected }) => {
      it(name, () => {
        const r1 = newEmail(email1);
        const r2 = newEmail(email2);
        if (r1.ok && r2.ok) {
          expect(r1.value.equals(r2.value)).toBe(expected);
        }
      });
    });
  });
});
