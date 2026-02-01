import { describe, it, expect } from "vitest";
import { newRecord } from "./record";
import { newRecordItem, RecordItem } from "./recordItem";

describe("newRecord", () => {
  // テスト用のRecordItemを作成するヘルパー関数
  const createRecordItem = (name: string, calories: number): RecordItem => {
    const result = newRecordItem({ name, calories });
    if (!result.ok) {
      throw new Error("Failed to create RecordItem for test");
    }
    return result.value;
  };

  describe("正常系", () => {
    const now = new Date("2024-01-15T12:00:00");

    const cases = [
      {
        name: "アイテムなしのレコード",
        input: {
          eatenAt: new Date("2024-01-15T08:00:00"),
          items: [] as readonly RecordItem[],
        },
        expectedTotalCalories: 0,
      },
      {
        name: "単一アイテムのレコード",
        input: {
          eatenAt: new Date("2024-01-15T08:00:00"),
          items: [createRecordItem("りんご", 100)] as readonly RecordItem[],
        },
        expectedTotalCalories: 100,
      },
      {
        name: "複数アイテムのレコード",
        input: {
          eatenAt: new Date("2024-01-15T08:00:00"),
          items: [
            createRecordItem("りんご", 100),
            createRecordItem("バナナ", 80),
            createRecordItem("ヨーグルト", 60),
          ] as readonly RecordItem[],
        },
        expectedTotalCalories: 240,
      },
    ];

    cases.forEach(({ name, input, expectedTotalCalories }) => {
      it(name, () => {
        const result = newRecord(input, now);
        expect(result.ok).toBe(true);
        if (result.ok) {
          expect(result.value.eatenAt.value.getTime()).toBe(input.eatenAt.getTime());
          expect(result.value.items.length).toBe(input.items.length);
          expect(result.value.totalCalories()).toBe(expectedTotalCalories);
        }
      });
    });
  });

  describe("異常系", () => {
    const now = new Date("2024-01-15T12:00:00");

    const cases = [
      {
        name: "食事日時が未来",
        input: {
          eatenAt: new Date("2024-01-15T13:00:00"),
          items: [createRecordItem("りんご", 100)] as readonly RecordItem[],
        },
        expectedCode: "EATEN_AT_MUST_NOT_BE_FUTURE",
      },
      {
        name: "食事日時が1日後",
        input: {
          eatenAt: new Date("2024-01-16T12:00:00"),
          items: [] as readonly RecordItem[],
        },
        expectedCode: "EATEN_AT_MUST_NOT_BE_FUTURE",
      },
    ];

    cases.forEach(({ name, input, expectedCode }) => {
      it(name, () => {
        const result = newRecord(input, now);
        expect(result.ok).toBe(false);
        if (!result.ok) {
          expect(result.error.eatenAt).toBeDefined();
          expect(result.error.eatenAt?.code).toBe(expectedCode);
        }
      });
    });
  });

  describe("totalCalories", () => {
    const now = new Date("2024-01-15T12:00:00");

    it("空のアイテムリストで0を返す", () => {
      const result = newRecord(
        {
          eatenAt: new Date("2024-01-15T08:00:00"),
          items: [],
        },
        now
      );
      expect(result.ok).toBe(true);
      if (result.ok) {
        expect(result.value.totalCalories()).toBe(0);
      }
    });

    it("複数アイテムの合計カロリーを正しく計算する", () => {
      const items = [
        createRecordItem("ご飯", 250),
        createRecordItem("味噌汁", 30),
        createRecordItem("焼き魚", 150),
        createRecordItem("サラダ", 20),
      ];
      const result = newRecord(
        {
          eatenAt: new Date("2024-01-15T08:00:00"),
          items,
        },
        now
      );
      expect(result.ok).toBe(true);
      if (result.ok) {
        expect(result.value.totalCalories()).toBe(450);
      }
    });

    it("totalCaloriesを複数回呼び出しても同じ結果を返す", () => {
      const items = [
        createRecordItem("りんご", 100),
        createRecordItem("バナナ", 80),
      ];
      const result = newRecord(
        {
          eatenAt: new Date("2024-01-15T08:00:00"),
          items,
        },
        now
      );
      expect(result.ok).toBe(true);
      if (result.ok) {
        expect(result.value.totalCalories()).toBe(180);
        expect(result.value.totalCalories()).toBe(180);
        expect(result.value.totalCalories()).toBe(180);
      }
    });
  });

  describe("mealType連携", () => {
    const now = new Date("2024-01-15T23:59:59");

    it("朝の時間帯でbreakfastを返す", () => {
      const result = newRecord(
        {
          eatenAt: new Date("2024-01-15T08:00:00"),
          items: [],
        },
        now
      );
      expect(result.ok).toBe(true);
      if (result.ok) {
        expect(result.value.eatenAt.mealType()).toBe("breakfast");
      }
    });

    it("昼の時間帯でlunchを返す", () => {
      const result = newRecord(
        {
          eatenAt: new Date("2024-01-15T12:00:00"),
          items: [],
        },
        now
      );
      expect(result.ok).toBe(true);
      if (result.ok) {
        expect(result.value.eatenAt.mealType()).toBe("lunch");
      }
    });

    it("夜の時間帯でdinnerを返す", () => {
      const result = newRecord(
        {
          eatenAt: new Date("2024-01-15T19:00:00"),
          items: [],
        },
        now
      );
      expect(result.ok).toBe(true);
      if (result.ok) {
        expect(result.value.eatenAt.mealType()).toBe("dinner");
      }
    });
  });
});
