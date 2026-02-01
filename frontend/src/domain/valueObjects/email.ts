import { Result, ok, err } from "../shared/result";
import { DomainError, domainError } from "../shared/errors";

export type EmailErrorCode =
  | "EMAIL_REQUIRED"
  | "EMAIL_TOO_LONG"
  | "EMAIL_INVALID_FORMAT";

export type EmailError = DomainError<EmailErrorCode>;

export type Email = Readonly<{
  value: string;
  equals: (other: Email) => boolean;
}>;

const MAX_LENGTH = 254;
const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

const ERROR_MESSAGE_EMAIL_REQUIRED = "メールアドレスを入力してください";
const ERROR_MESSAGE_EMAIL_TOO_LONG = "メールアドレスは254文字以内で入力してください";
const ERROR_MESSAGE_EMAIL_INVALID_FORMAT = "有効なメールアドレスを入力してください";

export const newEmail = (value: string): Result<Email, EmailError> => {
  if (!value || value.trim() === "") {
    return err(domainError("EMAIL_REQUIRED", ERROR_MESSAGE_EMAIL_REQUIRED));
  }
  if (value.length > MAX_LENGTH) {
    return err(domainError("EMAIL_TOO_LONG", ERROR_MESSAGE_EMAIL_TOO_LONG));
  }
  if (!EMAIL_REGEX.test(value)) {
    return err(domainError("EMAIL_INVALID_FORMAT", ERROR_MESSAGE_EMAIL_INVALID_FORMAT));
  }

  const email: Email = Object.freeze({
    value,
    equals: (other: Email) => value === other.value,
  });
  return ok(email);
};
