import { describe, it, expect } from "vitest";
import { newRecord } from "./record";

describe("newRecord", () => {
  describe("正常系", () => {
    const now = new Date("2024-01-15T12:00:00");

    const cases = [
      {
        name: "単一アイテムのレコード",
        input: {
          eatenAt: "2024-01-15T08:00:00",
          items: [{ name: "りんご", calories: 100 }],
        },
        expectedTotalCalories: 100,
      },
      {
        name: "複数アイテムのレコード",
        input: {
          eatenAt: "2024-01-15T08:00:00",
          items: [
            { name: "りんご", calories: 100 },
            { name: "バナナ", calories: 80 },
            { name: "ヨーグルト", calories: 60 },
          ],
        },
        expectedTotalCalories: 240,
      },
    ];

    cases.forEach(({ name, input, expectedTotalCalories }) => {
      it(name, () => {
        const result = newRecord(input, now);
        expect(result.ok).toBe(true);
        if (result.ok) {
          expect(result.value.eatenAt.value.getTime()).toBe(new Date(input.eatenAt).getTime());
          expect(result.value.items.length).toBe(input.items.length);
          expect(result.value.totalCalories()).toBe(expectedTotalCalories);
        }
      });
    });
  });

  describe("異常系 - eatenAtエラー", () => {
    const now = new Date("2024-01-15T12:00:00");

    const cases = [
      {
        name: "食事日時が空",
        input: {
          eatenAt: "",
          items: [{ name: "りんご", calories: 100 }],
        },
        expectedMessage: "食事日時を入力してください",
      },
      {
        name: "食事日時が無効",
        input: {
          eatenAt: "invalid-date",
          items: [{ name: "りんご", calories: 100 }],
        },
        expectedMessage: "有効な日時を入力してください",
      },
      {
        name: "食事日時が未来",
        input: {
          eatenAt: "2024-01-15T13:00:00",
          items: [{ name: "りんご", calories: 100 }],
        },
        expectedMessage: "食事日時は現在より過去を指定してください",
      },
    ];

    cases.forEach(({ name, input, expectedMessage }) => {
      it(name, () => {
        const result = newRecord(input, now);
        expect(result.ok).toBe(false);
        if (!result.ok) {
          expect(result.error.eatenAt).toBe(expectedMessage);
        }
      });
    });
  });

  describe("異常系 - itemsエラー", () => {
    const now = new Date("2024-01-15T12:00:00");

    it("アイテムが0件", () => {
      const result = newRecord(
        {
          eatenAt: "2024-01-15T08:00:00",
          items: [],
        },
        now
      );
      expect(result.ok).toBe(false);
      if (!result.ok) {
        expect(result.error.items.length).toBe(1);
        expect(result.error.items[0].name).toBe("少なくとも1つの食品を追加してください");
      }
    });

    it("アイテムの食品名が空", () => {
      const result = newRecord(
        {
          eatenAt: "2024-01-15T08:00:00",
          items: [{ name: "", calories: 100 }],
        },
        now
      );
      expect(result.ok).toBe(false);
      if (!result.ok) {
        expect(result.error.items.length).toBe(1);
        expect(result.error.items[0].name).toBe("食品名を入力してください");
      }
    });

    it("アイテムのカロリーが0", () => {
      const result = newRecord(
        {
          eatenAt: "2024-01-15T08:00:00",
          items: [{ name: "りんご", calories: 0 }],
        },
        now
      );
      expect(result.ok).toBe(false);
      if (!result.ok) {
        expect(result.error.items.length).toBe(1);
        expect(result.error.items[0].calories).toBe("カロリーは1以上の整数で入力してください");
      }
    });

    it("複数アイテムでエラーが複数ある", () => {
      const result = newRecord(
        {
          eatenAt: "2024-01-15T08:00:00",
          items: [
            { name: "", calories: 100 },
            { name: "バナナ", calories: 0 },
            { name: "りんご", calories: 50 },
          ],
        },
        now
      );
      expect(result.ok).toBe(false);
      if (!result.ok) {
        expect(result.error.items.length).toBe(3);
        expect(result.error.items[0].name).toBe("食品名を入力してください");
        expect(result.error.items[0].calories).toBeNull();
        expect(result.error.items[1].name).toBeNull();
        expect(result.error.items[1].calories).toBe("カロリーは1以上の整数で入力してください");
        expect(result.error.items[2].name).toBeNull();
        expect(result.error.items[2].calories).toBeNull();
      }
    });
  });

  describe("異常系 - 複合エラー", () => {
    const now = new Date("2024-01-15T12:00:00");

    it("eatenAtとitemsの両方にエラー", () => {
      const result = newRecord(
        {
          eatenAt: "",
          items: [{ name: "", calories: 0 }],
        },
        now
      );
      expect(result.ok).toBe(false);
      if (!result.ok) {
        expect(result.error.eatenAt).toBe("食事日時を入力してください");
        expect(result.error.items.length).toBe(1);
        expect(result.error.items[0].name).toBe("食品名を入力してください");
        expect(result.error.items[0].calories).toBe("カロリーは1以上の整数で入力してください");
      }
    });
  });

  describe("totalCalories", () => {
    const now = new Date("2024-01-15T12:00:00");

    it("複数アイテムの合計カロリーを正しく計算する", () => {
      const items = [
        { name: "ご飯", calories: 250 },
        { name: "味噌汁", calories: 30 },
        { name: "焼き魚", calories: 150 },
        { name: "サラダ", calories: 20 },
      ];
      const result = newRecord(
        {
          eatenAt: "2024-01-15T08:00:00",
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
        { name: "りんご", calories: 100 },
        { name: "バナナ", calories: 80 },
      ];
      const result = newRecord(
        {
          eatenAt: "2024-01-15T08:00:00",
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
          eatenAt: "2024-01-15T08:00:00",
          items: [{ name: "りんご", calories: 100 }],
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
          eatenAt: "2024-01-15T12:00:00",
          items: [{ name: "りんご", calories: 100 }],
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
          eatenAt: "2024-01-15T19:00:00",
          items: [{ name: "りんご", calories: 100 }],
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
