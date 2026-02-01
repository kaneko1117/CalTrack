import { Result, ok, err } from "../shared/result";
import { DomainError, domainError } from "../shared/errors";

export type GenderValue = "male" | "female" | "other";

export type GenderErrorCode = "GENDER_INVALID";

export type GenderError = DomainError<GenderErrorCode>;

export type Gender = Readonly<{
  value: GenderValue;
  equals: (other: Gender) => boolean;
}>;

const VALID_GENDERS: GenderValue[] = ["male", "female", "other"];

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
