import { describe, it, expect } from "vitest";
import { newPfcNutrient, PFC_NUTRIENT_OPTIONS } from "./pfcNutrient";

describe("PfcNutrient", () => {
  describe("正常系", () => {
    it("protein を正しく作成できる", () => {
      const result = newPfcNutrient("protein");
      expect(result.ok).toBe(true);
      if (result.ok) {
        expect(result.value.value).toBe("protein");
      }
    });

    it("fat を正しく作成できる", () => {
      const result = newPfcNutrient("fat");
      expect(result.ok).toBe(true);
      if (result.ok) {
        expect(result.value.value).toBe("fat");
      }
    });

    it("carbs を正しく作成できる", () => {
      const result = newPfcNutrient("carbs");
      expect(result.ok).toBe(true);
      if (result.ok) {
        expect(result.value.value).toBe("carbs");
      }
    });
  });

  describe("getLabel", () => {
    it('protein は "タンパク質" を返す', () => {
      const result = newPfcNutrient("protein");
      if (result.ok) {
        expect(result.value.getLabel()).toBe("タンパク質");
      }
    });

    it('fat は "脂質" を返す', () => {
      const result = newPfcNutrient("fat");
      if (result.ok) {
        expect(result.value.getLabel()).toBe("脂質");
      }
    });

    it('carbs は "炭水化物" を返す', () => {
      const result = newPfcNutrient("carbs");
      if (result.ok) {
        expect(result.value.getLabel()).toBe("炭水化物");
      }
    });
  });

  describe("getShortLabel", () => {
    it('protein は "P" を返す', () => {
      const result = newPfcNutrient("protein");
      if (result.ok) {
        expect(result.value.getShortLabel()).toBe("P");
      }
    });

    it('fat は "F" を返す', () => {
      const result = newPfcNutrient("fat");
      if (result.ok) {
        expect(result.value.getShortLabel()).toBe("F");
      }
    });

    it('carbs は "C" を返す', () => {
      const result = newPfcNutrient("carbs");
      if (result.ok) {
        expect(result.value.getShortLabel()).toBe("C");
      }
    });
  });

  describe("異常系", () => {
    it("空文字はエラーを返す", () => {
      const result = newPfcNutrient("");
      expect(result.ok).toBe(false);
      if (!result.ok) {
        expect(result.error.code).toBe("PFC_NUTRIENT_INVALID");
      }
    });

    it("unknown はエラーを返す", () => {
      const result = newPfcNutrient("unknown");
      expect(result.ok).toBe(false);
      if (!result.ok) {
        expect(result.error.code).toBe("PFC_NUTRIENT_INVALID");
      }
    });

    it("大文字 PROTEIN はエラーを返す", () => {
      const result = newPfcNutrient("PROTEIN");
      expect(result.ok).toBe(false);
      if (!result.ok) {
        expect(result.error.code).toBe("PFC_NUTRIENT_INVALID");
      }
    });

    it("日本語 タンパク質 はエラーを返す", () => {
      const result = newPfcNutrient("タンパク質");
      expect(result.ok).toBe(false);
      if (!result.ok) {
        expect(result.error.code).toBe("PFC_NUTRIENT_INVALID");
      }
    });
  });

  describe("equals", () => {
    it("同じ値同士は true を返す", () => {
      const protein1 = newPfcNutrient("protein");
      const protein2 = newPfcNutrient("protein");
      if (protein1.ok && protein2.ok) {
        expect(protein1.value.equals(protein2.value)).toBe(true);
      }
    });

    it("異なる値同士は false を返す", () => {
      const protein = newPfcNutrient("protein");
      const fat = newPfcNutrient("fat");
      if (protein.ok && fat.ok) {
        expect(protein.value.equals(fat.value)).toBe(false);
      }
    });
  });

  describe("PFC_NUTRIENT_OPTIONS", () => {
    it("3つの選択肢が存在する", () => {
      expect(PFC_NUTRIENT_OPTIONS).toHaveLength(3);
      expect(PFC_NUTRIENT_OPTIONS.map((opt) => opt.value)).toEqual([
        "protein",
        "fat",
        "carbs",
      ]);
    });
  });
});
