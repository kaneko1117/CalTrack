import { Result, ok, err } from "../shared/result";
import { DomainError, domainError } from "../shared/errors";

export type WeightErrorCode = "WEIGHT_REQUIRED" | "WEIGHT_INVALID" | "WEIGHT_MUST_BE_POSITIVE" | "WEIGHT_TOO_HEAVY";

export type WeightError = DomainError<WeightErrorCode>;

export type Weight = Readonly<{
  value: number;
  equals: (other: Weight) => boolean;
}>;

const MAX_WEIGHT = 500;

const ERROR_MESSAGE_WEIGHT_REQUIRED = "体重を入力してください";
const ERROR_MESSAGE_WEIGHT_INVALID = "体重は有効な数値を入力してください";
const ERROR_MESSAGE_WEIGHT_MUST_BE_POSITIVE = "体重は0より大きい値を入力してください";
const ERROR_MESSAGE_WEIGHT_TOO_HEAVY = "体重は500kg以内で入力してください";

export const newWeight = (input: string): Result<Weight, WeightError> => {
  // 空文字チェック
  if (input.trim() === "") {
    return err(domainError("WEIGHT_REQUIRED", ERROR_MESSAGE_WEIGHT_REQUIRED));
  }

  // 数値変換とNaNチェック
  const value = Number(input);
  if (isNaN(value)) {
    return err(domainError("WEIGHT_INVALID", ERROR_MESSAGE_WEIGHT_INVALID));
  }

  if (value <= 0) {
    return err(domainError("WEIGHT_MUST_BE_POSITIVE", ERROR_MESSAGE_WEIGHT_MUST_BE_POSITIVE));
  }
  if (value > MAX_WEIGHT) {
    return err(domainError("WEIGHT_TOO_HEAVY", ERROR_MESSAGE_WEIGHT_TOO_HEAVY));
  }

  const weight: Weight = Object.freeze({
    value,
    equals: (other: Weight) => value === other.value,
  });
  return ok(weight);
};
