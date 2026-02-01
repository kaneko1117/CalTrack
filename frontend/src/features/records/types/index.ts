/**
 * カロリー記録関連の型定義
 */

/** 食品アイテム（リクエスト用） */
export type RecordItemRequest = {
  name: string;
  calories: number;
};

/** 食品アイテム（レスポンス用） */
export type RecordItemResponse = {
  itemId: string;
  name: string;
  calories: number;
};

/** カロリー記録作成リクエスト */
export type CreateRecordRequest = {
  eatenAt: string;
  items: RecordItemRequest[];
};

/** カロリー記録作成レスポンス */
export type CreateRecordResponse = {
  recordId: string;
  eatenAt: string;
  totalCalories: number;
  items: RecordItemResponse[];
};
