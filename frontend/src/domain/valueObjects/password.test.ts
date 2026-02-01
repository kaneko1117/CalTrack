import { describe, it, expect } from "vitest";
import { newPassword } from "./password";

describe("newPassword", () => {
  describe("正常系", () => {
    const cases = [
      { name: "8文字ちょうど", input: "12345678" },
      { name: "8文字以上", input: "password123" },
      { name: "特殊文字を含む", input: "p@ssw0rd!" },
    ];

    cases.forEach(({ name, input }) => {
      it(name, () => {
        const result = newPassword(input);
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
        expectedCode: "PASSWORD_REQUIRED",
        expectedMessage: "パスワードを入力してください",
      },
      {
        name: "空白のみ",
        input: "   ",
        expectedCode: "PASSWORD_REQUIRED",
        expectedMessage: "パスワードを入力してください",
      },
      {
        name: "7文字",
        input: "1234567",
        expectedCode: "PASSWORD_TOO_SHORT",
        expectedMessage: "パスワードは8文字以上で入力してください",
      },
      {
        name: "1文字",
        input: "a",
        expectedCode: "PASSWORD_TOO_SHORT",
        expectedMessage: "パスワードは8文字以上で入力してください",
      },
    ];

    cases.forEach(({ name, input, expectedCode, expectedMessage }) => {
      it(name, () => {
        const result = newPassword(input);
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
        password1: "password123",
        password2: "password123",
        expected: true,
      },
      {
        name: "異なる値でfalseを返す",
        password1: "password123",
        password2: "password456",
        expected: false,
      },
    ];

    cases.forEach(({ name, password1, password2, expected }) => {
      it(name, () => {
        const r1 = newPassword(password1);
        const r2 = newPassword(password2);
        if (r1.ok && r2.ok) {
          expect(r1.value.equals(r2.value)).toBe(expected);
        }
      });
    });
  });
});
