import { Result, ok, err } from "../shared/result";
import { DomainError, domainError } from "../shared/errors";

export type HeightErrorCode = "HEIGHT_MUST_BE_POSITIVE" | "HEIGHT_TOO_TALL";

export type HeightError = DomainError<HeightErrorCode>;

export type Height = Readonly<{
  value: number;
  equals: (other: Height) => boolean;
}>;

const MAX_HEIGHT = 300;

const ERROR_MESSAGE_HEIGHT_MUST_BE_POSITIVE = "身長は0より大きい値を入力してください";
const ERROR_MESSAGE_HEIGHT_TOO_TALL = "身長は300cm以内で入力してください";

export const newHeight = (value: number): Result<Height, HeightError> => {
  if (value <= 0) {
    return err(domainError("HEIGHT_MUST_BE_POSITIVE", ERROR_MESSAGE_HEIGHT_MUST_BE_POSITIVE));
  }
  if (value > MAX_HEIGHT) {
    return err(domainError("HEIGHT_TOO_TALL", ERROR_MESSAGE_HEIGHT_TOO_TALL));
  }

  const height: Height = Object.freeze({
    value,
    equals: (other: Height) => value === other.value,
  });
  return ok(height);
};
