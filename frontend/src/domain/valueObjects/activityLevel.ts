import { Result, ok, err } from "../shared/result";
import { DomainError, domainError } from "../shared/errors";

export type ActivityLevelValue =
  | "sedentary"
  | "light"
  | "moderate"
  | "active"
  | "veryActive";

export type ActivityLevelErrorCode = "ACTIVITY_LEVEL_INVALID";

export type ActivityLevelError = DomainError<ActivityLevelErrorCode>;

export type ActivityLevel = Readonly<{
  value: ActivityLevelValue;
  equals: (other: ActivityLevel) => boolean;
}>;

const VALID_LEVELS: ActivityLevelValue[] = [
  "sedentary",
  "light",
  "moderate",
  "active",
  "veryActive",
];

const ERROR_MESSAGE_ACTIVITY_LEVEL_INVALID = "活動レベルを選択してください";

export const newActivityLevel = (value: string): Result<ActivityLevel, ActivityLevelError> => {
  if (!VALID_LEVELS.includes(value as ActivityLevelValue)) {
    return err(domainError("ACTIVITY_LEVEL_INVALID", ERROR_MESSAGE_ACTIVITY_LEVEL_INVALID));
  }

  const activityLevel: ActivityLevel = Object.freeze({
    value: value as ActivityLevelValue,
    equals: (other: ActivityLevel) => value === other.value,
  });
  return ok(activityLevel);
};
