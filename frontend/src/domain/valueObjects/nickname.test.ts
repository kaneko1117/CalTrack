import { describe, it, expect } from "vitest";
import { newNickname } from "./nickname";

describe("newNickname", () => {
  describe("正常系", () => {
    const cases = [
      { name: "通常のニックネーム", input: "たろう" },
      { name: "50文字ちょうど", input: "a".repeat(50) },
      { name: "1文字", input: "A" },
      { name: "日本語と英語の混在", input: "太郎Taro" },
    ];

    cases.forEach(({ name, input }) => {
      it(name, () => {
        const result = newNickname(input);
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
        expectedCode: "NICKNAME_REQUIRED",
        expectedMessage: "ニックネームを入力してください",
      },
      {
        name: "空白のみ",
        input: "   ",
        expectedCode: "NICKNAME_REQUIRED",
        expectedMessage: "ニックネームを入力してください",
      },
      {
        name: "51文字超過",
        input: "a".repeat(51),
        expectedCode: "NICKNAME_TOO_LONG",
        expectedMessage: "ニックネームは50文字以内で入力してください",
      },
      {
        name: "100文字超過",
        input: "a".repeat(100),
        expectedCode: "NICKNAME_TOO_LONG",
        expectedMessage: "ニックネームは50文字以内で入力してください",
      },
    ];

    cases.forEach(({ name, input, expectedCode, expectedMessage }) => {
      it(name, () => {
        const result = newNickname(input);
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
        nickname1: "たろう",
        nickname2: "たろう",
        expected: true,
      },
      {
        name: "異なる値でfalseを返す",
        nickname1: "たろう",
        nickname2: "はなこ",
        expected: false,
      },
    ];

    cases.forEach(({ name, nickname1, nickname2, expected }) => {
      it(name, () => {
        const r1 = newNickname(nickname1);
        const r2 = newNickname(nickname2);
        if (r1.ok && r2.ok) {
          expect(r1.value.equals(r2.value)).toBe(expected);
        }
      });
    });
  });
});
