import { Result, ok, err } from "../shared/result";
import { DomainError, domainError } from "../shared/errors";

export type CaloriesErrorCode = "CALORIES_MUST_BE_POSITIVE";

export type CaloriesError = DomainError<CaloriesErrorCode>;

export type Calories = Readonly<{
  value: number;
  equals: (other: Calories) => boolean;
}>;

const ERROR_MESSAGE_CALORIES_MUST_BE_POSITIVE = "カロリーは1以上の整数で入力してください";

export const newCalories = (value: number): Result<Calories, CaloriesError> => {
  if (value < 1) {
    return err(domainError("CALORIES_MUST_BE_POSITIVE", ERROR_MESSAGE_CALORIES_MUST_BE_POSITIVE));
  }

  const calories: Calories = Object.freeze({
    value,
    equals: (other: Calories) => value === other.value,
  });
  return ok(calories);
};

/**
 * カロリー値の配列を合計する
 * NaN値は0として扱う
 * @param values - カロリー値の配列
 * @returns 合計カロリー
 */
export const sumCalories = (values: number[]): number => {
  return values.reduce((sum, value) => sum + (isNaN(value) ? 0 : value), 0);
};
