import { Result, ok, err } from "../shared/result";
import { DomainError, domainError } from "../shared/errors";

export type EatenAtErrorCode =
  | "EATEN_AT_REQUIRED"
  | "EATEN_AT_INVALID"
  | "EATEN_AT_MUST_NOT_BE_FUTURE";

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
  mealTypeLabel: () => string;
  formattedTime: () => string;
  equals: (other: EatenAt) => boolean;
  toDateTimeLocal: () => string;
  toISOString: () => string;
}>;

/** エラーメッセージ定数 */
const ERROR_MESSAGES = {
  EATEN_AT_REQUIRED: "食事日時を入力してください",
  EATEN_AT_INVALID: "有効な日時を入力してください",
  EATEN_AT_MUST_NOT_BE_FUTURE: "食事日時は現在より過去を指定してください",
} as const;

/**
 * 時間帯から食事タイプを判定
 * @param date - 日時
 * @returns 食事タイプ
 */
const determineMealType = (date: Date): MealType => {
  const hour = date.getHours();
  if (hour >= 5 && hour < 11) return "breakfast";
  if (hour >= 11 && hour < 14) return "lunch";
  if (hour >= 14 && hour < 17) return "snack";
  if (hour >= 17 && hour < 21) return "dinner";
  return "lateNight";
};

/**
 * EatenAt Value Object を生成
 * @param value - datetime-local形式の文字列
 * @param now - 現在日時（テスト用にDI可能）
 * @returns Result<EatenAt, EatenAtError>
 */
export const newEatenAt = (
  value: string,
  now: Date = new Date()
): Result<EatenAt, EatenAtError> => {
  // 必須チェック
  if (!value) {
    return err(domainError("EATEN_AT_REQUIRED", ERROR_MESSAGES.EATEN_AT_REQUIRED));
  }

  // 形式チェック
  const date = new Date(value);
  if (isNaN(date.getTime())) {
    return err(domainError("EATEN_AT_INVALID", ERROR_MESSAGES.EATEN_AT_INVALID));
  }

  // 未来日時チェック
  if (date > now) {
    return err(domainError("EATEN_AT_MUST_NOT_BE_FUTURE", ERROR_MESSAGES.EATEN_AT_MUST_NOT_BE_FUTURE));
  }

  const mealType = determineMealType(date);

  const eatenAt: EatenAt = Object.freeze({
    value: date,
    mealType: () => mealType,
    mealTypeLabel: () => MEAL_TYPE_LABELS[mealType],
    formattedTime: () =>
      date.toLocaleTimeString("ja-JP", {
        hour: "2-digit",
        minute: "2-digit",
        hour12: false,
      }),
    equals: (other: EatenAt) => date.getTime() === other.value.getTime(),
    toDateTimeLocal: () => {
      const offset = date.getTimezoneOffset();
      const localDate = new Date(date.getTime() - offset * 60 * 1000);
      return localDate.toISOString().slice(0, 16);
    },
    toISOString: () => date.toISOString(),
  });

  return ok(eatenAt);
};
