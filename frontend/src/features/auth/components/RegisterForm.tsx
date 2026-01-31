/**
 * RegisterForm - ユーザー登録フォームコンポーネント
 * 新規ユーザー登録のためのフォームUI
 * Warm & Organicトーンのデザイン
 */
import * as React from "react";
import { useState } from "react";
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
export type RegisterFormProps = {
  /** 登録成功時のコールバック */
  onSuccess?: () => void;
};

/** フォームの内部状態 */
type FormState = {
  email: string;
  password: string;
  nickname: string;
  weight: string;
  height: string;
  birthDate: string;
  gender: Gender | "";
  activityLevel: ActivityLevel | "";
};

/** バリデーションエラー */
type FormErrors = {
  email?: string;
  password?: string;
  nickname?: string;
  weight?: string;
  height?: string;
  birthDate?: string;
  gender?: string;
  activityLevel?: string;
};

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

/** メールアドレスバリデーション用パターン（モジュールレベルでホイスト） */
const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

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
  } else if (!EMAIL_PATTERN.test(form.email)) {
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
 * AlertCircleアイコン - エラー表示用
 * SVGインラインアイコン
 */
function AlertCircleIcon({ className }: { className?: string }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      aria-hidden="true"
    >
      <circle cx="12" cy="12" r="10" />
      <line x1="12" y1="8" x2="12" y2="12" />
      <line x1="12" y1="16" x2="12.01" y2="16" />
    </svg>
  );
}

/**
 * CheckCircleアイコン - 成功表示用
 * SVGインラインアイコン
 */
function CheckCircleIcon({ className }: { className?: string }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      aria-hidden="true"
    >
      <path d="M22 11.08V12a10 10 0 1 1-5.93-9.14" />
      <polyline points="22 4 12 14.01 9 11.01" />
    </svg>
  );
}

/**
 * フィールドエラー表示コンポーネント
 * アイコン付きのエラーメッセージを表示
 */
function FieldError({ id, message }: { id: string; message: string }) {
  return (
    <p id={id} className="flex items-center gap-1.5 text-sm text-destructive">
      <AlertCircleIcon className="w-4 h-4 flex-shrink-0" />
      <span>{message}</span>
    </p>
  );
}

/**
 * RegisterForm - ユーザー登録フォーム
 */
export function RegisterForm({ onSuccess }: RegisterFormProps) {
  const [formState, setFormState] = useState<FormState>(initialFormState);
  const [formErrors, setFormErrors] = useState<FormErrors>({});
  const { register, isLoading, error, isSuccess, reset } = useRegisterUser();

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

    // API呼び出し（成功時のコールバックを引数として渡す）
    await register(
      {
        email: formState.email,
        password: formState.password,
        nickname: formState.nickname,
        weight: parseFloat(formState.weight),
        height: parseFloat(formState.height),
        birthDate: formState.birthDate,
        gender: formState.gender as Gender,
        activityLevel: formState.activityLevel as ActivityLevel,
      },
      onSuccess
    );
  };

  return (
    <Card className="w-full shadow-warm-lg border-0">
      <CardHeader className="space-y-1 pb-6">
        <CardTitle className="text-2xl font-semibold text-center">
          新規登録
        </CardTitle>
        <CardDescription className="text-center text-muted-foreground">
          アカウントを作成して、カロリー管理を始めましょう
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-5">
          {/* APIエラー表示 */}
          {error && (
            <div
              className="flex items-start gap-3 p-4 text-sm rounded-lg bg-destructive/10 border border-destructive/20"
              role="alert"
            >
              <AlertCircleIcon className="w-5 h-5 text-destructive flex-shrink-0 mt-0.5" />
              <div className="flex-1">
                <p className="font-medium text-destructive">
                  {getErrorMessage(error.code)}
                </p>
                {error.details && error.details.length > 0 && (
                  <ul className="mt-1.5 list-disc list-inside text-destructive/80">
                    {error.details.map((detail, index) => (
                      <li key={index}>{detail}</li>
                    ))}
                  </ul>
                )}
              </div>
            </div>
          )}

          {/* 成功メッセージ */}
          {isSuccess && (
            <div
              className="flex items-center gap-3 p-4 text-sm rounded-lg bg-success/10 border border-success/20"
              role="status"
            >
              <CheckCircleIcon className="w-5 h-5 text-success flex-shrink-0" />
              <p className="font-medium text-success">登録が完了しました</p>
            </div>
          )}

          {/* ニックネーム */}
          <div className="space-y-2">
            <Label htmlFor="nickname" className="text-foreground font-medium">
              ニックネーム
            </Label>
            <Input
              id="nickname"
              name="nickname"
              type="text"
              value={formState.nickname}
              onChange={handleChange}
              placeholder="ニックネームを入力"
              disabled={isLoading}
              aria-invalid={!!formErrors.nickname}
              aria-describedby={
                formErrors.nickname ? "nickname-error" : undefined
              }
              className="h-11 bg-background border-input focus:border-primary focus:ring-primary/20"
            />
            {formErrors.nickname && (
              <FieldError id="nickname-error" message={formErrors.nickname} />
            )}
          </div>

          {/* メールアドレス */}
          <div className="space-y-2">
            <Label htmlFor="email" className="text-foreground font-medium">
              メールアドレス
            </Label>
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
              className="h-11 bg-background border-input focus:border-primary focus:ring-primary/20"
            />
            {formErrors.email && (
              <FieldError id="email-error" message={formErrors.email} />
            )}
          </div>

          {/* パスワード */}
          <div className="space-y-2">
            <Label htmlFor="password" className="text-foreground font-medium">
              パスワード
            </Label>
            <Input
              id="password"
              name="password"
              type="password"
              value={formState.password}
              onChange={handleChange}
              placeholder="8文字以上で入力"
              disabled={isLoading}
              aria-invalid={!!formErrors.password}
              aria-describedby={
                formErrors.password ? "password-error" : undefined
              }
              className="h-11 bg-background border-input focus:border-primary focus:ring-primary/20"
            />
            {formErrors.password && (
              <FieldError id="password-error" message={formErrors.password} />
            )}
          </div>

          {/* 体重・身長（2カラム） */}
          <div className="grid grid-cols-2 gap-4">
            {/* 体重 */}
            <div className="space-y-2">
              <Label htmlFor="weight" className="text-foreground font-medium">
                体重 (kg)
              </Label>
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
                aria-describedby={
                  formErrors.weight ? "weight-error" : undefined
                }
                className="h-11 bg-background border-input focus:border-primary focus:ring-primary/20"
              />
              {formErrors.weight && (
                <FieldError id="weight-error" message={formErrors.weight} />
              )}
            </div>

            {/* 身長 */}
            <div className="space-y-2">
              <Label htmlFor="height" className="text-foreground font-medium">
                身長 (cm)
              </Label>
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
                aria-describedby={
                  formErrors.height ? "height-error" : undefined
                }
                className="h-11 bg-background border-input focus:border-primary focus:ring-primary/20"
              />
              {formErrors.height && (
                <FieldError id="height-error" message={formErrors.height} />
              )}
            </div>
          </div>

          {/* 生年月日 */}
          <div className="space-y-2">
            <Label htmlFor="birthDate" className="text-foreground font-medium">
              生年月日
            </Label>
            <Input
              id="birthDate"
              name="birthDate"
              type="date"
              value={formState.birthDate}
              onChange={handleChange}
              disabled={isLoading}
              aria-invalid={!!formErrors.birthDate}
              aria-describedby={
                formErrors.birthDate ? "birthDate-error" : undefined
              }
              className="h-11 bg-background border-input focus:border-primary focus:ring-primary/20"
            />
            {formErrors.birthDate && (
              <FieldError id="birthDate-error" message={formErrors.birthDate} />
            )}
          </div>

          {/* 性別 */}
          <div className="space-y-2">
            <Label htmlFor="gender" className="text-foreground font-medium">
              性別
            </Label>
            <Select
              id="gender"
              name="gender"
              value={formState.gender}
              onChange={handleChange}
              disabled={isLoading}
              aria-invalid={!!formErrors.gender}
              aria-describedby={formErrors.gender ? "gender-error" : undefined}
              className="h-11"
            >
              <SelectOption value="">選択してください</SelectOption>
              {GENDER_OPTIONS.map((option) => (
                <SelectOption key={option.value} value={option.value}>
                  {option.label}
                </SelectOption>
              ))}
            </Select>
            {formErrors.gender && (
              <FieldError id="gender-error" message={formErrors.gender} />
            )}
          </div>

          {/* 活動レベル */}
          <div className="space-y-2">
            <Label
              htmlFor="activityLevel"
              className="text-foreground font-medium"
            >
              活動レベル
            </Label>
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
              className="h-11"
            >
              <SelectOption value="">選択してください</SelectOption>
              {ACTIVITY_LEVEL_OPTIONS.map((option) => (
                <SelectOption key={option.value} value={option.value}>
                  {option.label}
                </SelectOption>
              ))}
            </Select>
            {formErrors.activityLevel && (
              <FieldError
                id="activityLevel-error"
                message={formErrors.activityLevel}
              />
            )}
          </div>

          {/* 送信ボタン */}
          <Button
            type="submit"
            className="w-full h-12 text-base font-medium mt-6 bg-primary hover:bg-primary/90 transition-colors"
            disabled={isLoading}
          >
            {isLoading ? "登録中..." : "登録する"}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
