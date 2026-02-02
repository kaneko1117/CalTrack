/**
 * LoginForm - ログインフォームコンポーネント
 */
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
import { useForm, getApiErrorMessage } from "@/features/common";
import { newEmail } from "@/domain/valueObjects/email";
import { newPassword } from "@/domain/valueObjects/password";
import { post } from "@/lib/api";

/** ログインレスポンス */
export type LoginResponse = {
  userId: string;
  email: string;
  nickname: string;
};

/** ログインAPI */
const login = (data: { email: string; password: string }) =>
  post<LoginResponse>("/api/v1/auth/login", data);

/** LoginFormコンポーネントのProps */
export type LoginFormProps = {
  onSuccess?: (response: LoginResponse) => void;
};

/** フォームフィールド型 */
type LoginField = "email" | "password";

/** フォームの初期状態 */
const initialFormState: Record<LoginField, string> = {
  email: "",
  password: "",
};

/** エラーの初期状態 */
const initialErrors: Record<LoginField, string | null> = {
  email: null,
  password: null,
};

/** VOファクトリ設定 */
const formConfig = {
  email: newEmail,
  password: newPassword,
};

/**
 * AlertCircleアイコン - エラー表示用
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
  const {
    formState,
    errors,
    apiError,
    handleChange,
    handleSubmit,
    isValid,
    isPending,
  } = useForm(
    formConfig,
    initialFormState,
    initialErrors,
    (data) => login({ email: data.email, password: data.password }),
    onSuccess
  );

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
          {apiError && (
            <div
              className="flex items-start gap-3 p-4 text-sm rounded-lg bg-destructive/10 border border-destructive/20"
              role="alert"
            >
              <AlertCircleIcon className="w-5 h-5 text-destructive flex-shrink-0 mt-0.5" />
              <div className="flex-1">
                <p className="font-medium text-destructive">
                  {getApiErrorMessage(apiError.code)}
                </p>
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
              onChange={(e) => handleChange("email")(e.target.value)}
              placeholder="example@example.com"
              disabled={isPending}
              aria-invalid={!!errors.email}
              aria-describedby={errors.email ? "email-error" : undefined}
              className="h-11 bg-background border-input focus:border-primary focus:ring-primary/20"
            />
            {errors.email && (
              <FieldError id="email-error" message={errors.email} />
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
              onChange={(e) => handleChange("password")(e.target.value)}
              placeholder="パスワードを入力"
              disabled={isPending}
              aria-invalid={!!errors.password}
              aria-describedby={errors.password ? "password-error" : undefined}
              className="h-11 bg-background border-input focus:border-primary focus:ring-primary/20"
            />
            {errors.password && (
              <FieldError id="password-error" message={errors.password} />
            )}
          </div>

          {/* 送信ボタン */}
          <Button
            type="submit"
            className="w-full h-12 text-base font-medium mt-6 bg-primary hover:bg-primary/90 transition-colors"
            disabled={!isValid || isPending}
          >
            {isPending ? "ログイン中..." : "ログイン"}
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
