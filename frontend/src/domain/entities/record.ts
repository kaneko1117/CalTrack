import { Result, ok, err } from "../shared/result";
import {
  EatenAt,
  newEatenAt,
} from "../valueObjects";
import { RecordItem, NewRecordItemInput, newRecordItem } from "./recordItem";

export type Record = Readonly<{
  eatenAt: EatenAt;
  items: readonly RecordItem[];
  totalCalories: () => number;
}>;

/** バリデーションエラー型（フォームエラー表示用） */
export type RecordValidationErrors = {
  eatenAt: string | null;
  items: Array<{ name: string | null; calories: string | null }>;
};

/** 新規Record生成の入力型 */
export type NewRecordInput = {
  eatenAt: string;
  items: readonly NewRecordItemInput[];
};

/** エラーメッセージ定数 */
const ERROR_MESSAGE_ITEMS_REQUIRED = "少なくとも1つの食品を追加してください";

/**
 * Record Entity を生成
 * @param input - 入力データ（eatenAt: string, items: NewRecordItemInput[]）
 * @param now - 現在日時（テスト用にDI可能）
 * @returns Result<Record, RecordValidationErrors>
 */
export const newRecord = (
  input: NewRecordInput,
  now: Date = new Date()
): Result<Record, RecordValidationErrors> => {
  const errors: RecordValidationErrors = {
    eatenAt: null,
    items: [],
  };
  let hasError = false;

  // eatenAtのバリデーション
  const eatenAtResult = newEatenAt(input.eatenAt, now);
  if (!eatenAtResult.ok) {
    errors.eatenAt = eatenAtResult.error.message;
    hasError = true;
  }

  // itemsのバリデーション
  if (input.items.length === 0) {
    errors.items = [{
      name: ERROR_MESSAGE_ITEMS_REQUIRED,
      calories: null,
    }];
    hasError = true;
  } else {
    // 各アイテムをバリデーション
    const itemResults = input.items.map((item) => newRecordItem(item));
    errors.items = itemResults.map((result) => {
      if (result.ok) {
        return { name: null, calories: null };
      }
      hasError = true;
      return {
        name: result.error.name?.message ?? null,
        calories: result.error.calories?.message ?? null,
      };
    });
  }

  if (hasError) {
    return err(errors);
  }

  // バリデーション成功時、RecordItemを生成
  const items = input.items.map((item) => {
    const result = newRecordItem(item);
    return (result as { ok: true; value: RecordItem }).value;
  });

  const record: Record = Object.freeze({
    eatenAt: (eatenAtResult as { ok: true; value: EatenAt }).value,
    items,
    totalCalories: () => items.reduce((sum, item) => sum + item.calories.value, 0),
  });

  return ok(record);
};
