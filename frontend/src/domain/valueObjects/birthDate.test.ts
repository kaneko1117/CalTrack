import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import { newBirthDate } from "./birthDate";

describe("newBirthDate", () => {
  // テストの安定性のため、現在時刻をモック
  beforeEach(() => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date("2024-01-15T12:00:00.000Z"));
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  describe("正常系", () => {
    const cases = [
      { name: "昨日の日付", input: new Date("2024-01-14") },
      { name: "100年前", input: new Date("1924-01-15") },
      { name: "149年前", input: new Date("1875-01-16") },
      { name: "通常の生年月日", input: new Date("1990-05-15") },
    ];

    cases.forEach(({ name, input }) => {
      it(name, () => {
        const result = newBirthDate(input);
        expect(result.ok).toBe(true);
        if (result.ok) {
          expect(result.value.value).toEqual(input);
        }
      });
    });
  });

  describe("異常系", () => {
    const cases = [
      {
        name: "今日(現在時刻と同じ)",
        input: new Date("2024-01-15T12:00:00.000Z"),
        expectedCode: "BIRTH_DATE_MUST_BE_PAST",
        expectedMessage: "生年月日は過去の日付を入力してください",
      },
      {
        name: "未来の日付",
        input: new Date("2025-01-01"),
        expectedCode: "BIRTH_DATE_MUST_BE_PAST",
        expectedMessage: "生年月日は過去の日付を入力してください",
      },
      {
        name: "151年前",
        input: new Date("1873-01-14"),
        expectedCode: "BIRTH_DATE_TOO_OLD",
        expectedMessage: "生年月日は150年以内の日付を入力してください",
      },
      {
        name: "200年前",
        input: new Date("1824-01-15"),
        expectedCode: "BIRTH_DATE_TOO_OLD",
        expectedMessage: "生年月日は150年以内の日付を入力してください",
      },
    ];

    cases.forEach(({ name, input, expectedCode, expectedMessage }) => {
      it(name, () => {
        const result = newBirthDate(input);
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
        date1: new Date("1990-05-15"),
        date2: new Date("1990-05-15"),
        expected: true,
      },
      {
        name: "異なる値でfalseを返す",
        date1: new Date("1990-05-15"),
        date2: new Date("1995-10-20"),
        expected: false,
      },
    ];

    cases.forEach(({ name, date1, date2, expected }) => {
      it(name, () => {
        const r1 = newBirthDate(date1);
        const r2 = newBirthDate(date2);
        if (r1.ok && r2.ok) {
          expect(r1.value.equals(r2.value)).toBe(expected);
        }
      });
    });
  });
});
