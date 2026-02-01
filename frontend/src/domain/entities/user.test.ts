import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";
import { newUser, NewUserInput } from "./user";

describe("newUser", () => {
  // テストの安定性のため、現在時刻をモック
  beforeEach(() => {
    vi.useFakeTimers();
    vi.setSystemTime(new Date("2024-01-15T12:00:00.000Z"));
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  // 有効な入力データを生成するヘルパー
  const createValidInput = (): NewUserInput => ({
    email: "test@example.com",
    password: "password123",
    nickname: "たろう",
    weight: 65,
    height: 170,
    birthDate: new Date("1990-05-15"),
    gender: "male",
    activityLevel: "moderate",
  });

  describe("正常系", () => {
    const cases = [
      {
        name: "全フィールドが有効な場合Userを生成できる",
        input: () => createValidInput(),
        expectedValues: {
          email: "test@example.com",
          password: "password123",
          nickname: "たろう",
          weight: 65,
          height: 170,
          gender: "male",
          activityLevel: "moderate",
        },
      },
      {
        name: "女性ユーザーを生成できる",
        input: () => ({ ...createValidInput(), gender: "female" }),
        expectedValues: {
          gender: "female",
        },
      },
      {
        name: "活動レベルがsedentaryのユーザーを生成できる",
        input: () => ({ ...createValidInput(), activityLevel: "sedentary" }),
        expectedValues: {
          activityLevel: "sedentary",
        },
      },
      {
        name: "活動レベルがveryActiveのユーザーを生成できる",
        input: () => ({ ...createValidInput(), activityLevel: "veryActive" }),
        expectedValues: {
          activityLevel: "veryActive",
        },
      },
    ];

    cases.forEach(({ name, input, expectedValues }) => {
      it(name, () => {
        const result = newUser(input());
        expect(result.ok).toBe(true);
        if (result.ok) {
          if (expectedValues.email !== undefined) {
            expect(result.value.email.value).toBe(expectedValues.email);
          }
          if (expectedValues.password !== undefined) {
            expect(result.value.password.value).toBe(expectedValues.password);
          }
          if (expectedValues.nickname !== undefined) {
            expect(result.value.nickname.value).toBe(expectedValues.nickname);
          }
          if (expectedValues.weight !== undefined) {
            expect(result.value.weight.value).toBe(expectedValues.weight);
          }
          if (expectedValues.height !== undefined) {
            expect(result.value.height.value).toBe(expectedValues.height);
          }
          if (expectedValues.gender !== undefined) {
            expect(result.value.gender.value).toBe(expectedValues.gender);
          }
          if (expectedValues.activityLevel !== undefined) {
            expect(result.value.activityLevel.value).toBe(expectedValues.activityLevel);
          }
        }
      });
    });
  });

  describe("異常系 - 単一フィールドエラー", () => {
    const cases = [
      {
        name: "メールアドレスが無効",
        input: () => ({ ...createValidInput(), email: "" }),
        expectedErrorField: "email" as const,
        expectedErrorCode: "EMAIL_REQUIRED",
      },
      {
        name: "パスワードが短すぎる",
        input: () => ({ ...createValidInput(), password: "short" }),
        expectedErrorField: "password" as const,
        expectedErrorCode: "PASSWORD_TOO_SHORT",
      },
      {
        name: "ニックネームが空",
        input: () => ({ ...createValidInput(), nickname: "" }),
        expectedErrorField: "nickname" as const,
        expectedErrorCode: "NICKNAME_REQUIRED",
      },
      {
        name: "体重が0",
        input: () => ({ ...createValidInput(), weight: 0 }),
        expectedErrorField: "weight" as const,
        expectedErrorCode: "WEIGHT_MUST_BE_POSITIVE",
      },
      {
        name: "身長が0",
        input: () => ({ ...createValidInput(), height: 0 }),
        expectedErrorField: "height" as const,
        expectedErrorCode: "HEIGHT_MUST_BE_POSITIVE",
      },
      {
        name: "生年月日が未来",
        input: () => ({ ...createValidInput(), birthDate: new Date("2025-01-01") }),
        expectedErrorField: "birthDate" as const,
        expectedErrorCode: "BIRTH_DATE_MUST_BE_PAST",
      },
      {
        name: "性別が不正",
        input: () => ({ ...createValidInput(), gender: "invalid" }),
        expectedErrorField: "gender" as const,
        expectedErrorCode: "GENDER_INVALID",
      },
      {
        name: "活動レベルが不正",
        input: () => ({ ...createValidInput(), activityLevel: "invalid" }),
        expectedErrorField: "activityLevel" as const,
        expectedErrorCode: "ACTIVITY_LEVEL_INVALID",
      },
    ];

    cases.forEach(({ name, input, expectedErrorField, expectedErrorCode }) => {
      it(name, () => {
        const result = newUser(input());
        expect(result.ok).toBe(false);
        if (!result.ok) {
          expect(result.error[expectedErrorField]).toBeDefined();
          expect(result.error[expectedErrorField]?.code).toBe(expectedErrorCode);
        }
      });
    });
  });

  describe("異常系 - 複数フィールドエラー", () => {
    const cases = [
      {
        name: "emailとpassword両方が無効",
        input: (): NewUserInput => ({
          ...createValidInput(),
          email: "",
          password: "short",
        }),
        expectedErrors: {
          email: "EMAIL_REQUIRED",
          password: "PASSWORD_TOO_SHORT",
        },
        unexpectedErrors: ["nickname", "weight", "height", "birthDate", "gender", "activityLevel"],
      },
      {
        name: "全フィールドが無効",
        input: (): NewUserInput => ({
          email: "",
          password: "short",
          nickname: "",
          weight: 0,
          height: 0,
          birthDate: new Date("2025-01-01"),
          gender: "invalid",
          activityLevel: "invalid",
        }),
        expectedErrors: {
          email: "EMAIL_REQUIRED",
          password: "PASSWORD_TOO_SHORT",
          nickname: "NICKNAME_REQUIRED",
          weight: "WEIGHT_MUST_BE_POSITIVE",
          height: "HEIGHT_MUST_BE_POSITIVE",
          birthDate: "BIRTH_DATE_MUST_BE_PAST",
          gender: "GENDER_INVALID",
          activityLevel: "ACTIVITY_LEVEL_INVALID",
        },
        unexpectedErrors: [],
      },
    ];

    cases.forEach(({ name, input, expectedErrors, unexpectedErrors }) => {
      it(name, () => {
        const result = newUser(input());
        expect(result.ok).toBe(false);
        if (!result.ok) {
          // 期待するエラーが存在することを確認
          Object.entries(expectedErrors).forEach(([field, code]) => {
            const errorField = field as keyof typeof result.error;
            expect(result.error[errorField]).toBeDefined();
            expect(result.error[errorField]?.code).toBe(code);
          });

          // 期待しないエラーが存在しないことを確認
          unexpectedErrors.forEach((field) => {
            const errorField = field as keyof typeof result.error;
            expect(result.error[errorField]).toBeUndefined();
          });
        }
      });
    });
  });
});
