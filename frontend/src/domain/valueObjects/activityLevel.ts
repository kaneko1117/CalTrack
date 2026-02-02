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

/** 選択肢の型 */
export type ActivityLevelOption = {
  value: ActivityLevelValue;
  label: string;
};

/** 活動レベルの選択肢（UI用） */
export const ACTIVITY_LEVEL_OPTIONS: ActivityLevelOption[] = [
  { value: "sedentary", label: "座りがち（運動なし）" },
  { value: "light", label: "軽い（週1-3回運動）" },
  { value: "moderate", label: "適度（週3-5回運動）" },
  { value: "active", label: "活動的（週6-7回運動）" },
  { value: "veryActive", label: "非常に活動的（毎日激しい運動）" },
];

const VALID_LEVELS: ActivityLevelValue[] = ACTIVITY_LEVEL_OPTIONS.map(
  (o) => o.value,
);

const ERROR_MESSAGE_ACTIVITY_LEVEL_INVALID = "活動レベルを選択してください";

export const newActivityLevel = (
  value: string,
): Result<ActivityLevel, ActivityLevelError> => {
  if (!VALID_LEVELS.includes(value as ActivityLevelValue)) {
    return err(
      domainError(
        "ACTIVITY_LEVEL_INVALID",
        ERROR_MESSAGE_ACTIVITY_LEVEL_INVALID,
      ),
    );
  }

  const activityLevel: ActivityLevel = Object.freeze({
    value: value as ActivityLevelValue,
    equals: (other: ActivityLevel) => value === other.value,
  });
  return ok(activityLevel);
};
