import { Result, ok, err } from "../shared/result";
import { DomainError, domainError } from "../shared/errors";

export type BirthDateErrorCode = "BIRTH_DATE_REQUIRED" | "BIRTH_DATE_INVALID" | "BIRTH_DATE_MUST_BE_PAST" | "BIRTH_DATE_TOO_OLD";

export type BirthDateError = DomainError<BirthDateErrorCode>;

export type BirthDate = Readonly<{
  value: Date;
  equals: (other: BirthDate) => boolean;
}>;

const MAX_AGE_YEARS = 150;

const ERROR_MESSAGE_BIRTH_DATE_REQUIRED = "生年月日を入力してください";
const ERROR_MESSAGE_BIRTH_DATE_INVALID = "生年月日は有効な日付を入力してください";
const ERROR_MESSAGE_BIRTH_DATE_MUST_BE_PAST = "生年月日は過去の日付を入力してください";
const ERROR_MESSAGE_BIRTH_DATE_TOO_OLD = "生年月日は150年以内の日付を入力してください";

export const newBirthDate = (input: string): Result<BirthDate, BirthDateError> => {
  // 空文字チェック
  if (input.trim() === "") {
    return err(domainError("BIRTH_DATE_REQUIRED", ERROR_MESSAGE_BIRTH_DATE_REQUIRED));
  }

  // Date変換とInvalid Dateチェック
  const value = new Date(input);
  if (isNaN(value.getTime())) {
    return err(domainError("BIRTH_DATE_INVALID", ERROR_MESSAGE_BIRTH_DATE_INVALID));
  }

  const now = new Date();
  if (value >= now) {
    return err(domainError("BIRTH_DATE_MUST_BE_PAST", ERROR_MESSAGE_BIRTH_DATE_MUST_BE_PAST));
  }

  const minDate = new Date();
  minDate.setFullYear(minDate.getFullYear() - MAX_AGE_YEARS);
  if (value < minDate) {
    return err(domainError("BIRTH_DATE_TOO_OLD", ERROR_MESSAGE_BIRTH_DATE_TOO_OLD));
  }

  const birthDate: BirthDate = Object.freeze({
    value,
    equals: (other: BirthDate) => value.getTime() === other.value.getTime(),
  });
  return ok(birthDate);
};
