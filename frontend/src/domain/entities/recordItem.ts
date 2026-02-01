import { Result, ok, err } from "../shared/result";
import {
  ItemName,
  ItemNameError,
  newItemName,
  Calories,
  CaloriesError,
  newCalories,
} from "../valueObjects";

export type RecordItem = Readonly<{
  name: ItemName;
  calories: Calories;
}>;

export type RecordItemValidationErrors = {
  name?: ItemNameError;
  calories?: CaloriesError;
};

export type NewRecordItemInput = {
  name: string;
  calories: number;
};

export const newRecordItem = (
  input: NewRecordItemInput
): Result<RecordItem, RecordItemValidationErrors> => {
  const errors: RecordItemValidationErrors = {};

  const nameResult = newItemName(input.name);
  if (!nameResult.ok) errors.name = nameResult.error;

  const caloriesResult = newCalories(input.calories);
  if (!caloriesResult.ok) errors.calories = caloriesResult.error;

  if (Object.keys(errors).length > 0) {
    return err(errors);
  }

  const recordItem: RecordItem = Object.freeze({
    name: (nameResult as { ok: true; value: ItemName }).value,
    calories: (caloriesResult as { ok: true; value: Calories }).value,
  });

  return ok(recordItem);
};
