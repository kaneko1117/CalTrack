import { Result, ok, err } from "../shared/result";
import {
  Email,
  EmailError,
  newEmail,
  Password,
  PasswordError,
  newPassword,
  Nickname,
  NicknameError,
  newNickname,
  Weight,
  WeightError,
  newWeight,
  Height,
  HeightError,
  newHeight,
  BirthDate,
  BirthDateError,
  newBirthDate,
  Gender,
  GenderError,
  newGender,
  ActivityLevel,
  ActivityLevelError,
  newActivityLevel,
} from "../valueObjects";

export type User = Readonly<{
  email: Email;
  password: Password;
  nickname: Nickname;
  weight: Weight;
  height: Height;
  birthDate: BirthDate;
  gender: Gender;
  activityLevel: ActivityLevel;
}>;

export type UserValidationErrors = {
  email?: EmailError;
  password?: PasswordError;
  nickname?: NicknameError;
  weight?: WeightError;
  height?: HeightError;
  birthDate?: BirthDateError;
  gender?: GenderError;
  activityLevel?: ActivityLevelError;
};

export type NewUserInput = {
  email: string;
  password: string;
  nickname: string;
  weight: number;
  height: number;
  birthDate: Date;
  gender: string;
  activityLevel: string;
};

export const newUser = (input: NewUserInput): Result<User, UserValidationErrors> => {
  const errors: UserValidationErrors = {};

  const emailResult = newEmail(input.email);
  if (!emailResult.ok) errors.email = emailResult.error;

  const passwordResult = newPassword(input.password);
  if (!passwordResult.ok) errors.password = passwordResult.error;

  const nicknameResult = newNickname(input.nickname);
  if (!nicknameResult.ok) errors.nickname = nicknameResult.error;

  const weightResult = newWeight(input.weight);
  if (!weightResult.ok) errors.weight = weightResult.error;

  const heightResult = newHeight(input.height);
  if (!heightResult.ok) errors.height = heightResult.error;

  const birthDateResult = newBirthDate(input.birthDate);
  if (!birthDateResult.ok) errors.birthDate = birthDateResult.error;

  const genderResult = newGender(input.gender);
  if (!genderResult.ok) errors.gender = genderResult.error;

  const activityLevelResult = newActivityLevel(input.activityLevel);
  if (!activityLevelResult.ok) errors.activityLevel = activityLevelResult.error;

  if (Object.keys(errors).length > 0) {
    return err(errors);
  }

  const user: User = Object.freeze({
    email: (emailResult as { ok: true; value: Email }).value,
    password: (passwordResult as { ok: true; value: Password }).value,
    nickname: (nicknameResult as { ok: true; value: Nickname }).value,
    weight: (weightResult as { ok: true; value: Weight }).value,
    height: (heightResult as { ok: true; value: Height }).value,
    birthDate: (birthDateResult as { ok: true; value: BirthDate }).value,
    gender: (genderResult as { ok: true; value: Gender }).value,
    activityLevel: (activityLevelResult as { ok: true; value: ActivityLevel }).value,
  });

  return ok(user);
};
