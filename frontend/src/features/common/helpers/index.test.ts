import { describe, it, expect, vi } from "vitest";
import {
  createFieldHandler,
  createResetHandler,
  getApiErrorMessage,
} from "./index";
import { ok, err } from "@/domain/shared/result";
import { domainError } from "@/domain/shared/errors";

describe("createFieldHandler", () => {
  describe("正常系", () => {
    it("バリデーション成功時、値をセットしエラーをnullにする", () => {
      const setFormState = vi.fn();
      const setErrors = vi.fn();
      const factory = vi.fn(() => ok({ value: "test@example.com" }));

      const handler = createFieldHandler("email", factory, setFormState, setErrors);
      handler("test@example.com");

      expect(setFormState).toHaveBeenCalledWith(expect.any(Function));
      expect(setErrors).toHaveBeenCalledWith(expect.any(Function));

      // setFormStateの関数を実行して確認
      const formStateFn = setFormState.mock.calls[0][0];
      expect(formStateFn({ email: "", password: "" })).toEqual({
        email: "test@example.com",
        password: "",
      });

      // setErrorsの関数を実行して確認
      const errorsFn = setErrors.mock.calls[0][0];
      expect(errorsFn({ email: "error", password: null })).toEqual({
        email: null,
        password: null,
      });
    });
  });

  describe("異常系", () => {
    it("バリデーション失敗時、値をセットしエラーメッセージをセットする", () => {
      const setFormState = vi.fn();
      const setErrors = vi.fn();
      const factory = vi.fn(() =>
        err(domainError("EMAIL_REQUIRED", "メールアドレスを入力してください"))
      );

      const handler = createFieldHandler("email", factory, setFormState, setErrors);
      handler("");

      expect(setFormState).toHaveBeenCalled();
      expect(setErrors).toHaveBeenCalled();

      // setErrorsの関数を実行して確認
      const errorsFn = setErrors.mock.calls[0][0];
      expect(errorsFn({ email: null, password: null })).toEqual({
        email: "メールアドレスを入力してください",
        password: null,
      });
    });
  });

  describe("フィールド更新", () => {
    it("指定したフィールドのみ更新される", () => {
      const setFormState = vi.fn();
      const setErrors = vi.fn();
      const factory = vi.fn(() => ok({ value: "newvalue" }));

      const handler = createFieldHandler("password", factory, setFormState, setErrors);
      handler("newvalue");

      const formStateFn = setFormState.mock.calls[0][0];
      expect(formStateFn({ email: "test@example.com", password: "" })).toEqual({
        email: "test@example.com",
        password: "newvalue",
      });
    });
  });
});

describe("createResetHandler", () => {
  it("フォームとエラーをリセットする", () => {
    const resetFormState = vi.fn();
    const resetErrors = vi.fn();

    const handler = createResetHandler(resetFormState, resetErrors);
    handler();

    expect(resetFormState).toHaveBeenCalledTimes(1);
    expect(resetErrors).toHaveBeenCalledTimes(1);
  });
});

describe("getApiErrorMessage", () => {
  const testCases = [
    { code: "INVALID_REQUEST", expected: "リクエストが不正です" },
    { code: "VALIDATION_ERROR", expected: "入力内容に誤りがあります" },
    { code: "UNAUTHORIZED", expected: "認証が必要です" },
    {
      code: "EMAIL_ALREADY_EXISTS",
      expected: "このメールアドレスは既に登録されています",
    },
    {
      code: "INVALID_CREDENTIALS",
      expected: "メールアドレスまたはパスワードが間違っています",
    },
    { code: "INTERNAL_ERROR", expected: "予期しないエラーが発生しました" },
  ];

  it.each(testCases)(
    "エラーコード $code に対して正しいメッセージを返す",
    ({ code, expected }) => {
      expect(getApiErrorMessage(code)).toBe(expected);
    }
  );

  it("未知のエラーコードにはデフォルトメッセージを返す", () => {
    expect(getApiErrorMessage("UNKNOWN_ERROR")).toBe(
      "予期しないエラーが発生しました"
    );
  });
});
