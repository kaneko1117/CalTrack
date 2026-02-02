/**
 * RegisterForm - ユーザー登録フォームコンポーネント
 * 新規ユーザー登録のためのフォームUI
 * Warm & Organicトーンのデザイン
 */
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
import { useForm, getApiErrorMessage } from "@/features/common";
import {
  newEmail,
  newPassword,
  newNickname,
  newWeight,
  newHeight,
  newBirthDate,
  newGender,
  newActivityLevel,
  GENDER_OPTIONS,
  ACTIVITY_LEVEL_OPTIONS,
} from "@/domain/valueObjects";
import type { GenderValue, ActivityLevelValue } from "@/domain/valueObjects";
import { post } from "@/lib/api";
import { err } from "@/domain/shared/result";

/** ユーザー登録レスポンス */
export type RegisterUserResponse = {
  userId: string;
  email: string;
  nickname: string;
};

/** ユーザー登録リクエストデータ */
type RegisterUserRequest = {
  email: string;
  password: string;
  nickname: string;
  weight: number;
  height: number;
  birthDate: string;
  gender: GenderValue;
  activityLevel: ActivityLevelValue;
};

/** ユーザー登録API */
const registerUser = (data: RegisterUserRequest) =>
  post<RegisterUserResponse>("/api/v1/auth/register", data);

/** RegisterFormコンポーネントのProps */
export type RegisterFormProps = {
  /** 登録成功時のコールバック */
  onSuccess?: (response: RegisterUserResponse) => void;
};

/** フォームフィールド型 */
type RegisterField =
  | "email"
  | "password"
  | "nickname"
  | "weight"
  | "height"
  | "birthDate"
  | "gender"
  | "activityLevel";

/** フォームの初期状態 */
const initialFormState: Record<RegisterField, string> = {
  email: "",
  password: "",
  nickname: "",
  weight: "",
  height: "",
  birthDate: "",
  gender: "",
  activityLevel: "",
};

/** エラーの初期状態 */
const initialErrors: Record<RegisterField, string | null> = {
  email: null,
  password: null,
  nickname: null,
  weight: null,
  height: null,
  birthDate: null,
  gender: null,
  activityLevel: null,
};

/**
 * 文字列をnumberに変換してWeightを生成するラッパー
 */
const newWeightFromString = (value: string) => {
  const num = parseFloat(value);
  if (isNaN(num)) {
    return err({
      code: "WEIGHT_MUST_BE_POSITIVE" as const,
      message: "体重を入力してください",
    });
  }
  return newWeight(num);
};

/**
 * 文字列をnumberに変換してHeightを生成するラッパー
 */
const newHeightFromString = (value: string) => {
  const num = parseFloat(value);
  if (isNaN(num)) {
    return err({
      code: "HEIGHT_MUST_BE_POSITIVE" as const,
      message: "身長を入力してください",
    });
  }
  return newHeight(num);
};

/**
 * 文字列をDateに変換してBirthDateを生成するラッパー
 */
const newBirthDateFromString = (value: string) => {
  if (!value) {
    return err({
      code: "BIRTH_DATE_MUST_BE_PAST" as const,
      message: "生年月日を入力してください",
    });
  }
  const date = new Date(value);
  if (isNaN(date.getTime())) {
    return err({
      code: "BIRTH_DATE_MUST_BE_PAST" as const,
      message: "有効な日付を入力してください",
    });
  }
  return newBirthDate(date);
};

/** VOファクトリ設定 */
const formConfig = {
  email: newEmail,
  password: newPassword,
  nickname: newNickname,
  weight: newWeightFromString,
  height: newHeightFromString,
  birthDate: newBirthDateFromString,
  gender: newGender,
  activityLevel: newActivityLevel,
};

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
 * RegisterForm - ユーザー登録フォーム
 */
export function RegisterForm({ onSuccess }: RegisterFormProps) {
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
    (data) =>
      registerUser({
        email: data.email,
        password: data.password,
        nickname: data.nickname,
        weight: parseFloat(data.weight),
        height: parseFloat(data.height),
        birthDate: data.birthDate,
        gender: data.gender as GenderValue,
        activityLevel: data.activityLevel as ActivityLevelValue,
      }),
    onSuccess
  );

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
                {apiError.details && apiError.details.length > 0 && (
                  <ul className="mt-1.5 list-disc list-inside text-destructive/80">
                    {apiError.details.map((detail, index) => (
                      <li key={index}>{detail}</li>
                    ))}
                  </ul>
                )}
              </div>
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
              onChange={(e) => handleChange("nickname")(e.target.value)}
              placeholder="ニックネームを入力"
              disabled={isPending}
              aria-invalid={!!errors.nickname}
              aria-describedby={errors.nickname ? "nickname-error" : undefined}
              className="h-11 bg-background border-input focus:border-primary focus:ring-primary/20"
            />
            {errors.nickname && (
              <FieldError id="nickname-error" message={errors.nickname} />
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
              placeholder="8文字以上で入力"
              disabled={isPending}
              aria-invalid={!!errors.password}
              aria-describedby={errors.password ? "password-error" : undefined}
              className="h-11 bg-background border-input focus:border-primary focus:ring-primary/20"
            />
            {errors.password && (
              <FieldError id="password-error" message={errors.password} />
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
                onChange={(e) => handleChange("weight")(e.target.value)}
                placeholder="60"
                disabled={isPending}
                aria-invalid={!!errors.weight}
                aria-describedby={errors.weight ? "weight-error" : undefined}
                className="h-11 bg-background border-input focus:border-primary focus:ring-primary/20"
              />
              {errors.weight && (
                <FieldError id="weight-error" message={errors.weight} />
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
                onChange={(e) => handleChange("height")(e.target.value)}
                placeholder="170"
                disabled={isPending}
                aria-invalid={!!errors.height}
                aria-describedby={errors.height ? "height-error" : undefined}
                className="h-11 bg-background border-input focus:border-primary focus:ring-primary/20"
              />
              {errors.height && (
                <FieldError id="height-error" message={errors.height} />
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
              onChange={(e) => handleChange("birthDate")(e.target.value)}
              disabled={isPending}
              aria-invalid={!!errors.birthDate}
              aria-describedby={errors.birthDate ? "birthDate-error" : undefined}
              className="h-11 bg-background border-input focus:border-primary focus:ring-primary/20"
            />
            {errors.birthDate && (
              <FieldError id="birthDate-error" message={errors.birthDate} />
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
              onChange={(e) => handleChange("gender")(e.target.value)}
              disabled={isPending}
              aria-invalid={!!errors.gender}
              aria-describedby={errors.gender ? "gender-error" : undefined}
              className="h-11"
            >
              <SelectOption value="">選択してください</SelectOption>
              {GENDER_OPTIONS.map((option) => (
                <SelectOption key={option.value} value={option.value}>
                  {option.label}
                </SelectOption>
              ))}
            </Select>
            {errors.gender && (
              <FieldError id="gender-error" message={errors.gender} />
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
              onChange={(e) => handleChange("activityLevel")(e.target.value)}
              disabled={isPending}
              aria-invalid={!!errors.activityLevel}
              aria-describedby={
                errors.activityLevel ? "activityLevel-error" : undefined
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
            {errors.activityLevel && (
              <FieldError
                id="activityLevel-error"
                message={errors.activityLevel}
              />
            )}
          </div>

          {/* 送信ボタン */}
          <Button
            type="submit"
            className="w-full h-12 text-base font-medium mt-6 bg-primary hover:bg-primary/90 transition-colors"
            disabled={!isValid || isPending}
          >
            {isPending ? "登録中..." : "登録する"}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
