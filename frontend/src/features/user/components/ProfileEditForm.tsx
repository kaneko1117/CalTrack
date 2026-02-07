/**
 * ProfileEditForm - プロフィール編集フォームコンポーネント
 * ユーザー情報（ニックネーム、体重、身長、活動レベル）の編集UI
 * Warm & Organicトーンのデザイン
 */
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Select, SelectOption } from "@/components/ui/select";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "@/components/ui/card";
import { getApiErrorMessage } from "@/features/common";
import { ACTIVITY_LEVEL_OPTIONS, GENDER_OPTIONS } from "@/domain/valueObjects";
import { useCurrentUser } from "../hooks/useCurrentUser";
import { useUpdateProfile } from "../hooks/useUpdateProfile";
import type { UpdateProfileResponse } from "../api";

/**
 * 読み取り専用フィールドコンポーネント
 * メールアドレス、生年月日、性別など編集不可なフィールドを表示
 */
function ReadOnlyField({ label, value }: { label: string; value: string }) {
  return (
    <div className="space-y-2">
      <Label className="text-muted-foreground font-medium text-sm">
        {label}
      </Label>
      <div className="h-11 flex items-center px-3 rounded-md bg-muted text-foreground text-sm">
        {value}
      </div>
    </div>
  );
}

/**
 * 性別値をラベルに変換
 */
function getGenderLabel(value: string): string {
  const option = GENDER_OPTIONS.find((o) => o.value === value);
  return option ? option.label : value;
}

/**
 * ISO日付を日本語形式にフォーマット (YYYY年M月D日)
 */
function formatBirthDate(isoDate: string): string {
  const date = new Date(isoDate);
  if (isNaN(date.getTime())) return isoDate;
  const year = date.getFullYear();
  const month = date.getMonth() + 1;
  const day = date.getDate();
  return `${year}年${month}月${day}日`;
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
 * ローディングスケルトン - データ読み込み中の表示（7フィールド対応）
 */
function LoadingSkeleton() {
  return (
    <Card className="w-full shadow-warm-lg border-0">
      <CardHeader className="space-y-1 pb-6">
        <Skeleton className="h-8 w-48 mx-auto" />
        <Skeleton className="h-4 w-64 mx-auto" />
      </CardHeader>
      <CardContent>
        <div className="space-y-5">
          {/* メールアドレス */}
          <div className="space-y-2">
            <Skeleton className="h-5 w-32" />
            <Skeleton className="h-11 w-full" />
          </div>
          {/* ニックネーム */}
          <div className="space-y-2">
            <Skeleton className="h-5 w-32" />
            <Skeleton className="h-11 w-full" />
          </div>
          {/* 生年月日・性別 */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Skeleton className="h-5 w-24" />
              <Skeleton className="h-11 w-full" />
            </div>
            <div className="space-y-2">
              <Skeleton className="h-5 w-24" />
              <Skeleton className="h-11 w-full" />
            </div>
          </div>
          {/* 体重・身長 */}
          <div className="grid grid-cols-2 gap-4">
            <div className="space-y-2">
              <Skeleton className="h-5 w-24" />
              <Skeleton className="h-11 w-full" />
            </div>
            <div className="space-y-2">
              <Skeleton className="h-5 w-24" />
              <Skeleton className="h-11 w-full" />
            </div>
          </div>
          {/* 活動レベル */}
          <div className="space-y-2">
            <Skeleton className="h-5 w-32" />
            <Skeleton className="h-11 w-full" />
          </div>
          {/* 送信ボタン */}
          <div className="!mt-10">
            <Skeleton className="h-12 w-full" />
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

/** ProfileEditFormコンポーネントのProps */
export type ProfileEditFormProps = {
  /** 更新成功時のコールバック */
  onSuccess?: (response: UpdateProfileResponse) => void;
};

/**
 * ProfileEditForm - プロフィール編集フォーム
 */
export function ProfileEditForm({ onSuccess }: ProfileEditFormProps) {
  const { data: currentUser, error: fetchError, isLoading } = useCurrentUser();
  const {
    formState,
    errors,
    apiError,
    handleChange,
    handleSubmit,
    isValid,
    isPending,
  } = useUpdateProfile(currentUser, onSuccess);

  // ローディング中
  if (isLoading) {
    return <LoadingSkeleton />;
  }

  // 取得エラー
  if (fetchError) {
    return (
      <Card className="w-full shadow-warm-lg border-0">
        <CardContent className="pt-6">
          <div
            className="flex items-start gap-3 p-4 text-sm rounded-lg bg-destructive/10 border border-destructive/20"
            role="alert"
          >
            <AlertCircleIcon className="w-5 h-5 text-destructive flex-shrink-0 mt-0.5" />
            <div className="flex-1">
              <p className="font-medium text-destructive">
                プロフィールの読み込みに失敗しました
              </p>
              <p className="mt-1 text-destructive/80">
                {getApiErrorMessage(fetchError.code)}
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="w-full shadow-warm-lg border-0">
      <CardHeader className="space-y-1 pb-6">
        <CardTitle className="text-2xl font-semibold text-center">
          プロフィール編集
        </CardTitle>
        <CardDescription className="text-center text-muted-foreground">
          ユーザー情報を確認・更新できます
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

          {/* メールアドレス（読み取り専用） */}
          {currentUser && (
            <ReadOnlyField label="メールアドレス" value={currentUser.email} />
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

          {/* 生年月日・性別（2カラム、読み取り専用） */}
          {currentUser && (
            <div className="grid grid-cols-2 gap-4">
              <ReadOnlyField
                label="生年月日"
                value={formatBirthDate(currentUser.birthDate)}
              />
              <ReadOnlyField
                label="性別"
                value={getGenderLabel(currentUser.gender)}
              />
            </div>
          )}

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
          <div className="!mt-10">
            <Button
              type="submit"
              className="w-full h-12 text-base font-medium bg-primary hover:bg-primary/90 transition-colors"
              disabled={!isValid || isPending}
            >
              {isPending ? "更新中..." : "更新する"}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}
