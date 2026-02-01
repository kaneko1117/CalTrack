import { Result, ok, err } from "../shared/result";
import {
  EatenAt,
  EatenAtError,
  newEatenAt,
} from "../valueObjects";
import { RecordItem } from "./recordItem";

export type Record = Readonly<{
  eatenAt: EatenAt;
  items: readonly RecordItem[];
  totalCalories: () => number;
}>;

export type RecordValidationErrors = {
  eatenAt?: EatenAtError;
};

export type NewRecordInput = {
  eatenAt: Date;
  items: readonly RecordItem[];
};

export const newRecord = (
  input: NewRecordInput,
  now: Date = new Date()
): Result<Record, RecordValidationErrors> => {
  const errors: RecordValidationErrors = {};

  const eatenAtResult = newEatenAt(input.eatenAt, now);
  if (!eatenAtResult.ok) errors.eatenAt = eatenAtResult.error;

  if (Object.keys(errors).length > 0) {
    return err(errors);
  }

  const items = input.items;

  const record: Record = Object.freeze({
    eatenAt: (eatenAtResult as { ok: true; value: EatenAt }).value,
    items,
    totalCalories: () => items.reduce((sum, item) => sum + item.calories.value, 0),
  });

  return ok(record);
};
