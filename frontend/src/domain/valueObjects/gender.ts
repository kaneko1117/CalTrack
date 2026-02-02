import { Result, ok, err } from "../shared/result";
import { DomainError, domainError } from "../shared/errors";

export type GenderValue = "male" | "female" | "other";

export type GenderErrorCode = "GENDER_INVALID";

export type GenderError = DomainError<GenderErrorCode>;

export type Gender = Readonly<{
  value: GenderValue;
  equals: (other: Gender) => boolean;
}>;

/** 選択肢の型 */
export type GenderOption = {
  value: GenderValue;
  label: string;
};

/** 性別の選択肢（UI用） */
export const GENDER_OPTIONS: GenderOption[] = [
  { value: "male", label: "男性" },
  { value: "female", label: "女性" },
  { value: "other", label: "その他" },
];

const VALID_GENDERS: GenderValue[] = GENDER_OPTIONS.map((o) => o.value);

const ERROR_MESSAGE_GENDER_INVALID = "性別を選択してください";

export const newGender = (value: string): Result<Gender, GenderError> => {
  if (!VALID_GENDERS.includes(value as GenderValue)) {
    return err(domainError("GENDER_INVALID", ERROR_MESSAGE_GENDER_INVALID));
  }

  const gender: Gender = Object.freeze({
    value: value as GenderValue,
    equals: (other: Gender) => value === other.value,
  });
  return ok(gender);
};
