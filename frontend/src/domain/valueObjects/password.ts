import { Result, ok, err } from "../shared/result";
import { DomainError, domainError } from "../shared/errors";

export type PasswordErrorCode = "PASSWORD_REQUIRED" | "PASSWORD_TOO_SHORT";

export type PasswordError = DomainError<PasswordErrorCode>;

export type Password = Readonly<{
  value: string;
  equals: (other: Password) => boolean;
}>;

const MIN_LENGTH = 8;

const ERROR_MESSAGE_PASSWORD_REQUIRED = "パスワードを入力してください";
const ERROR_MESSAGE_PASSWORD_TOO_SHORT = "パスワードは8文字以上で入力してください";

export const newPassword = (value: string): Result<Password, PasswordError> => {
  if (!value || value.trim() === "") {
    return err(domainError("PASSWORD_REQUIRED", ERROR_MESSAGE_PASSWORD_REQUIRED));
  }
  if (value.length < MIN_LENGTH) {
    return err(domainError("PASSWORD_TOO_SHORT", ERROR_MESSAGE_PASSWORD_TOO_SHORT));
  }

  const password: Password = Object.freeze({
    value,
    equals: (other: Password) => value === other.value,
  });
  return ok(password);
};
