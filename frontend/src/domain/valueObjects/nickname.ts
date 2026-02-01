import { Result, ok, err } from "../shared/result";
import { DomainError, domainError } from "../shared/errors";

export type NicknameErrorCode = "NICKNAME_REQUIRED" | "NICKNAME_TOO_LONG";

export type NicknameError = DomainError<NicknameErrorCode>;

export type Nickname = Readonly<{
  value: string;
  equals: (other: Nickname) => boolean;
}>;

const MAX_LENGTH = 50;

const ERROR_MESSAGE_NICKNAME_REQUIRED = "ニックネームを入力してください";
const ERROR_MESSAGE_NICKNAME_TOO_LONG = "ニックネームは50文字以内で入力してください";

export const newNickname = (value: string): Result<Nickname, NicknameError> => {
  if (!value || value.trim() === "") {
    return err(domainError("NICKNAME_REQUIRED", ERROR_MESSAGE_NICKNAME_REQUIRED));
  }
  if (value.length > MAX_LENGTH) {
    return err(domainError("NICKNAME_TOO_LONG", ERROR_MESSAGE_NICKNAME_TOO_LONG));
  }

  const nickname: Nickname = Object.freeze({
    value,
    equals: (other: Nickname) => value === other.value,
  });
  return ok(nickname);
};
