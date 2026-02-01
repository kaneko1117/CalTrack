import { Result, ok, err } from "../shared/result";
import { DomainError, domainError } from "../shared/errors";

export type EatenAtErrorCode = "EATEN_AT_MUST_NOT_BE_FUTURE";

export type EatenAtError = DomainError<EatenAtErrorCode>;

export type MealType = "breakfast" | "lunch" | "snack" | "dinner" | "lateNight";

export const MEAL_TYPE_LABELS: Record<MealType, string> = {
  breakfast: "朝食",
  lunch: "昼食",
  snack: "間食",
  dinner: "夕食",
  lateNight: "夜食",
};

export type EatenAt = Readonly<{
  value: Date;
  mealType: () => MealType;
  equals: (other: EatenAt) => boolean;
}>;

const ERROR_MESSAGE_EATEN_AT_MUST_NOT_BE_FUTURE = "食事日時は現在より過去を指定してください";

const determineMealType = (date: Date): MealType => {
  const hour = date.getHours();
  if (hour >= 5 && hour < 11) return "breakfast";
  if (hour >= 11 && hour < 14) return "lunch";
  if (hour >= 14 && hour < 17) return "snack";
  if (hour >= 17 && hour < 21) return "dinner";
  return "lateNight";
};

export const newEatenAt = (value: Date, now: Date = new Date()): Result<EatenAt, EatenAtError> => {
  if (value > now) {
    return err(domainError("EATEN_AT_MUST_NOT_BE_FUTURE", ERROR_MESSAGE_EATEN_AT_MUST_NOT_BE_FUTURE));
  }

  const eatenAt: EatenAt = Object.freeze({
    value,
    mealType: () => determineMealType(value),
    equals: (other: EatenAt) => value.getTime() === other.value.getTime(),
  });
  return ok(eatenAt);
};
