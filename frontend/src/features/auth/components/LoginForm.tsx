/**
 * LoginForm - ログインフォームコンポーネント
 * メールアドレスとパスワードによるログインUI
 * Warm & Organicトーンのデザイン
 */
import * as React from "react";
import { useState } from "react";
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
} from "@/components/ui/card";
import { useLogin } from "../hooks";
import type { LoginResponse } from "../types";

/** LoginFormコンポーネントのProps */
export type LoginFormProps = {
  /** ログイン成功時のコールバック */
  onSuccess?: (response: LoginResponse) => void;
};

/** フォームの内部状態 */
type FormState = {
  email: string;
  password: string;
};

/** バリデーションエラー */
type FormErrors = {
  email?: string;
  password?: string;
};

/** フォームの初期状態 */
const initialFormState: FormState = {
  email: "",
  password: "",
};

/** メールアドレスバリデーション用パターン（モジュールレベルでホイスト） */
const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

/**
 * フォームバリデーション関数
 * @param form - フォームの状態
 * @returns バリデーションエラー
 */
function validateForm(form: FormState): FormErrors {
  const errors: FormErrors = {};

  // email: 必須、形式チェック
  if (!form.email.trim()) {
    errors.email = "メールアドレスを入力してください";
  } else if (!EMAIL_PATTERN.test(form.email)) {
    errors.email = "正しいメールアドレス形式で入力してください";
  }

  // password: 必須
  if (!form.password) {
    errors.password = "パスワードを入力してください";
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
    case "INVALID_CREDENTIALS":
      return "メールアドレスまたはパスワードが間違っています";
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
 * LoginForm - ログインフォーム
 */
export function LoginForm({ onSuccess }: LoginFormProps) {
  const [formState, setFormState] = useState<FormState>(initialFormState);
  const [formErrors, setFormErrors] = useState<FormErrors>({});
  const { login, isLoading, error, reset } = useLogin();

  /**
   * フィールド値の変更ハンドラ
   */
  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
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
    await login(
      {
        email: formState.email,
        password: formState.password,
      },
      onSuccess
    );
  };

  return (
    <Card className="w-full shadow-warm-lg border-0">
      <CardHeader className="space-y-1 pb-6">
        <CardTitle className="text-2xl font-semibold text-center">
          ログイン
        </CardTitle>
        <CardDescription className="text-center text-muted-foreground">
          アカウントにログインしてください
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
              placeholder="パスワードを入力"
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

          {/* 送信ボタン */}
          <Button
            type="submit"
            className="w-full h-12 text-base font-medium mt-6 bg-primary hover:bg-primary/90 transition-colors"
            disabled={isLoading}
          >
            {isLoading ? "ログイン中..." : "ログイン"}
          </Button>
        </form>
      </CardContent>
      <CardFooter className="flex justify-center pb-6">
        <p className="text-sm text-muted-foreground">
          アカウントをお持ちでない方は{" "}
          <Link
            to="/register"
            className="font-medium text-primary hover:underline"
          >
            新規登録
          </Link>
        </p>
      </CardFooter>
    </Card>
  );
}
