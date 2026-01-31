/**
 * RegisterForm - ユーザー登録フォームコンポーネント
 * 新規ユーザー登録のためのフォームUI
 */
import * as React from "react";
import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectOption } from "@/components/ui/select";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "@/components/ui/card";
import { useRegisterUser } from "../hooks";
import type { Gender, ActivityLevel } from "../types";

/** RegisterFormコンポーネントのProps */
export interface RegisterFormProps {
  /** 登録成功時のコールバック */
  onSuccess?: () => void;
}

/** フォームの内部状態 */
interface FormState {
  email: string;
  password: string;
  nickname: string;
  weight: string;
  height: string;
  birthDate: string;
  gender: Gender | "";
  activityLevel: ActivityLevel | "";
}

/** バリデーションエラー */
interface FormErrors {
  email?: string;
  password?: string;
  nickname?: string;
  weight?: string;
  height?: string;
  birthDate?: string;
  gender?: string;
  activityLevel?: string;
}

/** フォームの初期状態 */
const initialFormState: FormState = {
  email: "",
  password: "",
  nickname: "",
  weight: "",
  height: "",
  birthDate: "",
  gender: "",
  activityLevel: "",
};

/** 性別の選択肢 */
const GENDER_OPTIONS = [
  { value: "male", label: "男性" },
  { value: "female", label: "女性" },
  { value: "other", label: "その他" },
] as const;

/** 活動レベルの選択肢 */
const ACTIVITY_LEVEL_OPTIONS = [
  { value: "sedentary", label: "座りがち（運動なし）" },
  { value: "light", label: "軽い（週1-3回運動）" },
  { value: "moderate", label: "適度（週3-5回運動）" },
  { value: "active", label: "活動的（週6-7回運動）" },
  { value: "veryActive", label: "非常に活動的（毎日激しい運動）" },
] as const;

/**
 * フォームバリデーション関数
 * @param form - フォームの状態
 * @returns バリデーションエラー
 */
function validateForm(form: FormState): FormErrors {
  const errors: FormErrors = {};

  // nickname: 必須
  if (!form.nickname.trim()) {
    errors.nickname = "ニックネームを入力してください";
  }

  // email: 必須、形式チェック
  if (!form.email.trim()) {
    errors.email = "メールアドレスを入力してください";
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
    errors.email = "正しいメールアドレス形式で入力してください";
  }

  // password: 必須、8文字以上
  if (!form.password) {
    errors.password = "パスワードを入力してください";
  } else if (form.password.length < 8) {
    errors.password = "パスワードは8文字以上で入力してください";
  }

  // weight: 必須、正の数
  const weight = parseFloat(form.weight);
  if (!form.weight) {
    errors.weight = "体重を入力してください";
  } else if (isNaN(weight) || weight <= 0) {
    errors.weight = "正しい体重を入力してください";
  }

  // height: 必須、正の数
  const height = parseFloat(form.height);
  if (!form.height) {
    errors.height = "身長を入力してください";
  } else if (isNaN(height) || height <= 0) {
    errors.height = "正しい身長を入力してください";
  }

  // birthDate: 必須、過去の日付
  if (!form.birthDate) {
    errors.birthDate = "生年月日を入力してください";
  } else if (new Date(form.birthDate) >= new Date()) {
    errors.birthDate = "過去の日付を入力してください";
  }

  // gender: 必須
  if (!form.gender) {
    errors.gender = "性別を選択してください";
  }

  // activityLevel: 必須
  if (!form.activityLevel) {
    errors.activityLevel = "活動レベルを選択してください";
  }

  return errors;
}

/**
 * APIエラーコードからユーザー向けメッセージを取得
 * @param code - エラーコード
 * @returns ユーザー向けメッセージ
 */
function getErrorMessage(code: string): string {
  switch (code) {
    case "EMAIL_ALREADY_EXISTS":
      return "このメールアドレスは既に登録されています";
    case "VALIDATION_ERROR":
      return "入力内容に誤りがあります";
    default:
      return "予期しないエラーが発生しました";
  }
}

/**
 * RegisterForm - ユーザー登録フォーム
 */
export function RegisterForm({ onSuccess }: RegisterFormProps) {
  const [formState, setFormState] = useState<FormState>(initialFormState);
  const [formErrors, setFormErrors] = useState<FormErrors>({});
  const { register, isLoading, error, isSuccess, reset } = useRegisterUser();

  // 成功時にコールバックを呼び出す
  useEffect(() => {
    if (isSuccess) {
      onSuccess?.();
    }
  }, [isSuccess, onSuccess]);

  /**
   * フィールド値の変更ハンドラ
   */
  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    const { name, value } = e.target;
    setFormState((prev) => ({ ...prev, [name]: value }));
    // 該当フィールドのエラーをクリア
    if (formErrors[name as keyof FormErrors]) {
      setFormErrors((prev) => ({ ...prev, [name]: undefined }));
    }
    // APIエラーをリセット
    if (error) {
      reset();
    }
  };

  /**
   * フォーム送信ハンドラ
   */
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // バリデーション
    const errors = validateForm(formState);
    if (Object.keys(errors).length > 0) {
      setFormErrors(errors);
      return;
    }

    // API呼び出し
    await register({
      email: formState.email,
      password: formState.password,
      nickname: formState.nickname,
      weight: parseFloat(formState.weight),
      height: parseFloat(formState.height),
      birthDate: formState.birthDate,
      gender: formState.gender as Gender,
      activityLevel: formState.activityLevel as ActivityLevel,
    });
  };

  return (
    <Card className="w-full max-w-md mx-auto">
      <CardHeader>
        <CardTitle>新規登録</CardTitle>
        <CardDescription>
          アカウントを作成して、カロリー管理を始めましょう
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-4">
          {/* APIエラー表示 */}
          {error && (
            <div
              className="p-3 text-sm text-red-500 bg-red-50 border border-red-200 rounded-md"
              role="alert"
            >
              <p>{getErrorMessage(error.code)}</p>
              {error.details && error.details.length > 0 && (
                <ul className="mt-1 list-disc list-inside">
                  {error.details.map((detail, index) => (
                    <li key={index}>{detail}</li>
                  ))}
                </ul>
              )}
            </div>
          )}

          {/* 成功メッセージ */}
          {isSuccess && (
            <div
              className="p-3 text-sm text-green-500 bg-green-50 border border-green-200 rounded-md"
              role="status"
            >
              登録が完了しました
            </div>
          )}

          {/* ニックネーム */}
          <div className="space-y-2">
            <Label htmlFor="nickname">ニックネーム</Label>
            <Input
              id="nickname"
              name="nickname"
              type="text"
              value={formState.nickname}
              onChange={handleChange}
              placeholder="ニックネームを入力"
              disabled={isLoading}
              aria-invalid={!!formErrors.nickname}
              aria-describedby={formErrors.nickname ? "nickname-error" : undefined}
            />
            {formErrors.nickname && (
              <p id="nickname-error" className="text-sm text-red-500">
                {formErrors.nickname}
              </p>
            )}
          </div>

          {/* メールアドレス */}
          <div className="space-y-2">
            <Label htmlFor="email">メールアドレス</Label>
            <Input
              id="email"
              name="email"
              type="email"
              value={formState.email}
              onChange={handleChange}
              placeholder="example@example.com"
              disabled={isLoading}
              aria-invalid={!!formErrors.email}
              aria-describedby={formErrors.email ? "email-error" : undefined}
            />
            {formErrors.email && (
              <p id="email-error" className="text-sm text-red-500">
                {formErrors.email}
              </p>
            )}
          </div>

          {/* パスワード */}
          <div className="space-y-2">
            <Label htmlFor="password">パスワード</Label>
            <Input
              id="password"
              name="password"
              type="password"
              value={formState.password}
              onChange={handleChange}
              placeholder="8文字以上で入力"
              disabled={isLoading}
              aria-invalid={!!formErrors.password}
              aria-describedby={formErrors.password ? "password-error" : undefined}
            />
            {formErrors.password && (
              <p id="password-error" className="text-sm text-red-500">
                {formErrors.password}
              </p>
            )}
          </div>

          {/* 体重 */}
          <div className="space-y-2">
            <Label htmlFor="weight">体重 (kg)</Label>
            <Input
              id="weight"
              name="weight"
              type="number"
              step="0.1"
              min="0"
              value={formState.weight}
              onChange={handleChange}
              placeholder="60"
              disabled={isLoading}
              aria-invalid={!!formErrors.weight}
              aria-describedby={formErrors.weight ? "weight-error" : undefined}
            />
            {formErrors.weight && (
              <p id="weight-error" className="text-sm text-red-500">
                {formErrors.weight}
              </p>
            )}
          </div>

          {/* 身長 */}
          <div className="space-y-2">
            <Label htmlFor="height">身長 (cm)</Label>
            <Input
              id="height"
              name="height"
              type="number"
              step="0.1"
              min="0"
              value={formState.height}
              onChange={handleChange}
              placeholder="170"
              disabled={isLoading}
              aria-invalid={!!formErrors.height}
              aria-describedby={formErrors.height ? "height-error" : undefined}
            />
            {formErrors.height && (
              <p id="height-error" className="text-sm text-red-500">
                {formErrors.height}
              </p>
            )}
          </div>

          {/* 生年月日 */}
          <div className="space-y-2">
            <Label htmlFor="birthDate">生年月日</Label>
            <Input
              id="birthDate"
              name="birthDate"
              type="date"
              value={formState.birthDate}
              onChange={handleChange}
              disabled={isLoading}
              aria-invalid={!!formErrors.birthDate}
              aria-describedby={formErrors.birthDate ? "birthDate-error" : undefined}
            />
            {formErrors.birthDate && (
              <p id="birthDate-error" className="text-sm text-red-500">
                {formErrors.birthDate}
              </p>
            )}
          </div>

          {/* 性別 */}
          <div className="space-y-2">
            <Label htmlFor="gender">性別</Label>
            <Select
              id="gender"
              name="gender"
              value={formState.gender}
              onChange={handleChange}
              disabled={isLoading}
              aria-invalid={!!formErrors.gender}
              aria-describedby={formErrors.gender ? "gender-error" : undefined}
            >
              <SelectOption value="">選択してください</SelectOption>
              {GENDER_OPTIONS.map((option) => (
                <SelectOption key={option.value} value={option.value}>
                  {option.label}
                </SelectOption>
              ))}
            </Select>
            {formErrors.gender && (
              <p id="gender-error" className="text-sm text-red-500">
                {formErrors.gender}
              </p>
            )}
          </div>

          {/* 活動レベル */}
          <div className="space-y-2">
            <Label htmlFor="activityLevel">活動レベル</Label>
            <Select
              id="activityLevel"
              name="activityLevel"
              value={formState.activityLevel}
              onChange={handleChange}
              disabled={isLoading}
              aria-invalid={!!formErrors.activityLevel}
              aria-describedby={
                formErrors.activityLevel ? "activityLevel-error" : undefined
              }
            >
              <SelectOption value="">選択してください</SelectOption>
              {ACTIVITY_LEVEL_OPTIONS.map((option) => (
                <SelectOption key={option.value} value={option.value}>
                  {option.label}
                </SelectOption>
              ))}
            </Select>
            {formErrors.activityLevel && (
              <p id="activityLevel-error" className="text-sm text-red-500">
                {formErrors.activityLevel}
              </p>
            )}
          </div>

          {/* 送信ボタン */}
          <Button type="submit" className="w-full" disabled={isLoading}>
            {isLoading ? "登録中..." : "登録する"}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
