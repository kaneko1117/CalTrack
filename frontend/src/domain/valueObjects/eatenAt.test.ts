import { describe, it, expect } from "vitest";
import { newEatenAt, MEAL_TYPE_LABELS } from "./eatenAt";

describe("newEatenAt", () => {
  describe("正常系", () => {
    const now = new Date("2024-01-15T12:00:00");

    const cases = [
      {
        name: "現在時刻と同じ",
        input: "2024-01-15T12:00:00",
        now,
      },
      {
        name: "1時間前",
        input: "2024-01-15T11:00:00",
        now,
      },
      {
        name: "1日前",
        input: "2024-01-14T12:00:00",
        now,
      },
      {
        name: "1ヶ月前",
        input: "2023-12-15T12:00:00",
        now,
      },
      {
        name: "datetime-local形式（秒なし）",
        input: "2024-01-15T11:00",
        now,
      },
    ];

    cases.forEach(({ name, input, now: nowDate }) => {
      it(name, () => {
        const result = newEatenAt(input, nowDate);
        expect(result.ok).toBe(true);
        if (result.ok) {
          expect(result.value.value.getTime()).toBe(new Date(input).getTime());
        }
      });
    });
  });

  describe("異常系", () => {
    const now = new Date("2024-01-15T12:00:00");

    const cases = [
      {
        name: "空文字（必須エラー）",
        input: "",
        now,
        expectedCode: "EATEN_AT_REQUIRED",
        expectedMessage: "食事日時を入力してください",
      },
      {
        name: "無効な形式（形式エラー）",
        input: "invalid-date",
        now,
        expectedCode: "EATEN_AT_INVALID",
        expectedMessage: "有効な日時を入力してください",
      },
      {
        name: "1秒後（未来）",
        input: "2024-01-15T12:00:01",
        now,
        expectedCode: "EATEN_AT_MUST_NOT_BE_FUTURE",
        expectedMessage: "食事日時は現在より過去を指定してください",
      },
      {
        name: "1時間後（未来）",
        input: "2024-01-15T13:00:00",
        now,
        expectedCode: "EATEN_AT_MUST_NOT_BE_FUTURE",
        expectedMessage: "食事日時は現在より過去を指定してください",
      },
      {
        name: "1日後（未来）",
        input: "2024-01-16T12:00:00",
        now,
        expectedCode: "EATEN_AT_MUST_NOT_BE_FUTURE",
        expectedMessage: "食事日時は現在より過去を指定してください",
      },
    ];

    cases.forEach(({ name, input, now: nowDate, expectedCode, expectedMessage }) => {
      it(name, () => {
        const result = newEatenAt(input, nowDate);
        expect(result.ok).toBe(false);
        if (!result.ok) {
          expect(result.error.code).toBe(expectedCode);
          expect(result.error.message).toBe(expectedMessage);
        }
      });
    });
  });

  describe("mealType判定", () => {
    const now = new Date("2024-01-15T23:59:59");

    const cases = [
      { name: "5時は朝食", input: "2024-01-15T05:00:00", expected: "breakfast" },
      { name: "10時は朝食", input: "2024-01-15T10:59:00", expected: "breakfast" },
      { name: "11時は昼食", input: "2024-01-15T11:00:00", expected: "lunch" },
      { name: "13時は昼食", input: "2024-01-15T13:59:00", expected: "lunch" },
      { name: "14時は間食", input: "2024-01-15T14:00:00", expected: "snack" },
      { name: "16時は間食", input: "2024-01-15T16:59:00", expected: "snack" },
      { name: "17時は夕食", input: "2024-01-15T17:00:00", expected: "dinner" },
      { name: "20時は夕食", input: "2024-01-15T20:59:00", expected: "dinner" },
      { name: "21時は夜食", input: "2024-01-15T21:00:00", expected: "lateNight" },
      { name: "0時は夜食", input: "2024-01-15T00:00:00", expected: "lateNight" },
      { name: "4時は夜食", input: "2024-01-15T04:59:00", expected: "lateNight" },
    ];

    cases.forEach(({ name, input, expected }) => {
      it(name, () => {
        const result = newEatenAt(input, now);
        expect(result.ok).toBe(true);
        if (result.ok) {
          expect(result.value.mealType()).toBe(expected);
        }
      });
    });
  });

  describe("MEAL_TYPE_LABELS", () => {
    it("全てのラベルが定義されている", () => {
      expect(MEAL_TYPE_LABELS.breakfast).toBe("朝食");
      expect(MEAL_TYPE_LABELS.lunch).toBe("昼食");
      expect(MEAL_TYPE_LABELS.snack).toBe("間食");
      expect(MEAL_TYPE_LABELS.dinner).toBe("夕食");
      expect(MEAL_TYPE_LABELS.lateNight).toBe("夜食");
    });
  });

  describe("equals", () => {
    const now = new Date("2024-01-15T12:00:00");

    const cases = [
      {
        name: "同じ日時でtrueを返す",
        date1: "2024-01-15T10:00:00",
        date2: "2024-01-15T10:00:00",
        expected: true,
      },
      {
        name: "異なる日時でfalseを返す",
        date1: "2024-01-15T10:00:00",
        date2: "2024-01-15T11:00:00",
        expected: false,
      },
    ];

    cases.forEach(({ name, date1, date2, expected }) => {
      it(name, () => {
        const r1 = newEatenAt(date1, now);
        const r2 = newEatenAt(date2, now);
        if (r1.ok && r2.ok) {
          expect(r1.value.equals(r2.value)).toBe(expected);
        }
      });
    });
  });

  describe("toDateTimeLocal", () => {
    const now = new Date("2024-01-15T23:59:59");

    it("datetime-local形式を返す", () => {
      const result = newEatenAt("2024-01-15T10:30:00", now);
      expect(result.ok).toBe(true);
      if (result.ok) {
        // toDateTimeLocalの結果を検証（タイムゾーンに依存しない形で検証）
        const dateTimeLocal = result.value.toDateTimeLocal();
        expect(dateTimeLocal).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}$/);
      }
    });
  });

  describe("toISOString", () => {
    const now = new Date("2024-01-15T23:59:59");

    it("ISO 8601形式を返す", () => {
      const result = newEatenAt("2024-01-15T10:30:00", now);
      expect(result.ok).toBe(true);
      if (result.ok) {
        const isoString = result.value.toISOString();
        expect(isoString).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.\d{3}Z$/);
      }
    });
  });
});
