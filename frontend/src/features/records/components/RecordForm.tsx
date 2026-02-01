/**
 * RecordForm - カロリー記録フォームコンポーネント
 * 食事日時と食品アイテムを入力してカロリー記録を作成
 * 動的なアイテム追加・削除、リアルタイム合計カロリー表示
 */
import * as React from "react";
import { useState, useMemo } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "@/components/ui/card";
import { useCreateRecord } from "../hooks";
import type { RecordItemRequest, CreateRecordResponse } from "../types";

/** RecordFormコンポーネントのProps */
export type RecordFormProps = {
  /** 記録作成成功時のコールバック */
  onSuccess?: (response: CreateRecordResponse) => void;
};

/** 食品アイテムの内部状態（IDを含む） */
type ItemState = RecordItemRequest & {
  id: string;
};

/** フォームの内部状態 */
type FormState = {
  eatenAt: string;
  items: ItemState[];
};

/** バリデーションエラー */
type FormErrors = {
  eatenAt?: string;
  items?: string;
  itemErrors?: { [key: string]: { name?: string; calories?: string } };
};

/** 一意なIDを生成する関数 */
function generateId(): string {
  return `item-${Date.now()}-${Math.random().toString(36).slice(2, 9)}`;
}

/** 現在日時をdatetime-local形式で取得 */
function getCurrentDateTimeLocal(): string {
  const now = new Date();
  const offset = now.getTimezoneOffset();
  const localDate = new Date(now.getTime() - offset * 60 * 1000);
  return localDate.toISOString().slice(0, 16);
}

/** フォームの初期状態 */
const createInitialFormState = (): FormState => ({
  eatenAt: getCurrentDateTimeLocal(),
  items: [{ id: generateId(), name: "", calories: 0 }],
});

/**
 * フォームバリデーション関数
 * @param form - フォームの状態
 * @returns バリデーションエラー
 */
function validateForm(form: FormState): FormErrors {
  const errors: FormErrors = {};
  const itemErrors: { [key: string]: { name?: string; calories?: string } } = {};

  // eatenAt: 必須
  if (!form.eatenAt) {
    errors.eatenAt = "食事日時を入力してください";
  }

  // items: 少なくとも1つ必要
  if (form.items.length === 0) {
    errors.items = "少なくとも1つの食品を追加してください";
  }

  // 各アイテムのバリデーション
  form.items.forEach((item) => {
    const itemError: { name?: string; calories?: string } = {};

    if (!item.name.trim()) {
      itemError.name = "食品名を入力してください";
    }

    if (item.calories < 0) {
      itemError.calories = "カロリーは0以上で入力してください";
    }

    if (Object.keys(itemError).length > 0) {
      itemErrors[item.id] = itemError;
    }
  });

  if (Object.keys(itemErrors).length > 0) {
    errors.itemErrors = itemErrors;
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
    case "VALIDATION_ERROR":
      return "入力内容に誤りがあります";
    case "UNAUTHORIZED":
      return "ログインが必要です";
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
 * PlusIcon - 追加ボタン用アイコン
 */
function PlusIcon({ className }: { className?: string }) {
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
      <line x1="12" y1="5" x2="12" y2="19" />
      <line x1="5" y1="12" x2="19" y2="12" />
    </svg>
  );
}

/**
 * TrashIcon - 削除ボタン用アイコン
 */
function TrashIcon({ className }: { className?: string }) {
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
      <polyline points="3,6 5,6 21,6" />
      <path d="M19,6v14a2,2 0 0,1-2,2H7a2,2 0 0,1-2-2V6m3,0V4a2,2 0 0,1,2-2h4a2,2 0 0,1,2,2v2" />
      <line x1="10" y1="11" x2="10" y2="17" />
      <line x1="14" y1="11" x2="14" y2="17" />
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
 * RecordForm - カロリー記録フォーム
 */
export function RecordForm({ onSuccess }: RecordFormProps) {
  const [formState, setFormState] = useState<FormState>(createInitialFormState);
  const [formErrors, setFormErrors] = useState<FormErrors>({});
  const { createRecord, isLoading, error, reset } = useCreateRecord();

  /** 合計カロリーの計算（メモ化） */
  const totalCalories = useMemo(() => {
    return formState.items.reduce((sum, item) => sum + (item.calories || 0), 0);
  }, [formState.items]);

  /**
   * 食事日時の変更ハンドラ
   */
  const handleEatenAtChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setFormState((prev) => ({ ...prev, eatenAt: value }));
    // エラーをクリア
    if (formErrors.eatenAt) {
      setFormErrors((prev) => ({ ...prev, eatenAt: undefined }));
    }
    if (error) {
      reset();
    }
  };

  /**
   * アイテムのフィールド変更ハンドラ
   */
  const handleItemChange = (
    itemId: string,
    field: "name" | "calories",
    value: string | number
  ) => {
    setFormState((prev) => ({
      ...prev,
      items: prev.items.map((item) =>
        item.id === itemId ? { ...item, [field]: value } : item
      ),
    }));
    // 該当アイテムのエラーをクリア
    if (formErrors.itemErrors?.[itemId]?.[field]) {
      setFormErrors((prev) => {
        const newItemErrors = { ...prev.itemErrors };
        if (newItemErrors[itemId]) {
          newItemErrors[itemId] = { ...newItemErrors[itemId], [field]: undefined };
          if (!newItemErrors[itemId].name && !newItemErrors[itemId].calories) {
            delete newItemErrors[itemId];
          }
        }
        return {
          ...prev,
          itemErrors: Object.keys(newItemErrors).length > 0 ? newItemErrors : undefined,
        };
      });
    }
    if (error) {
      reset();
    }
  };

  /**
   * アイテム追加ハンドラ
   */
  const handleAddItem = () => {
    setFormState((prev) => ({
      ...prev,
      items: [...prev.items, { id: generateId(), name: "", calories: 0 }],
    }));
    // アイテムエラーをクリア
    if (formErrors.items) {
      setFormErrors((prev) => ({ ...prev, items: undefined }));
    }
  };

  /**
   * アイテム削除ハンドラ
   */
  const handleRemoveItem = (itemId: string) => {
    setFormState((prev) => ({
      ...prev,
      items: prev.items.filter((item) => item.id !== itemId),
    }));
    // 該当アイテムのエラーを削除
    if (formErrors.itemErrors?.[itemId]) {
      setFormErrors((prev) => {
        const newItemErrors = { ...prev.itemErrors };
        delete newItemErrors[itemId];
        return {
          ...prev,
          itemErrors: Object.keys(newItemErrors).length > 0 ? newItemErrors : undefined,
        };
      });
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

    // datetime-localをISO 8601形式に変換
    const eatenAtISO = new Date(formState.eatenAt).toISOString();

    // API呼び出し
    await createRecord(
      {
        eatenAt: eatenAtISO,
        items: formState.items.map(({ name, calories }) => ({ name, calories })),
      },
      (response) => {
        // 成功時にフォームをリセット
        setFormState(createInitialFormState());
        setFormErrors({});
        onSuccess?.(response);
      }
    );
  };

  return (
    <Card className="w-full shadow-warm-lg border-0">
      <CardHeader className="space-y-1 pb-6">
        <CardTitle className="text-2xl font-semibold text-center">
          カロリー記録
        </CardTitle>
        <CardDescription className="text-center text-muted-foreground">
          食事の記録を追加してください
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-6">
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

          {/* 食事日時 */}
          <div className="space-y-2">
            <Label htmlFor="eatenAt" className="text-foreground font-medium">
              食事日時
            </Label>
            <Input
              id="eatenAt"
              name="eatenAt"
              type="datetime-local"
              value={formState.eatenAt}
              onChange={handleEatenAtChange}
              disabled={isLoading}
              aria-invalid={!!formErrors.eatenAt}
              aria-describedby={formErrors.eatenAt ? "eatenAt-error" : undefined}
              className="h-11 bg-background border-input focus:border-primary focus:ring-primary/20"
            />
            {formErrors.eatenAt && (
              <FieldError id="eatenAt-error" message={formErrors.eatenAt} />
            )}
          </div>

          {/* 食品アイテム一覧 */}
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <Label className="text-foreground font-medium">食品アイテム</Label>
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={handleAddItem}
                disabled={isLoading}
                className="gap-1"
              >
                <PlusIcon className="w-4 h-4" />
                追加
              </Button>
            </div>

            {formErrors.items && (
              <FieldError id="items-error" message={formErrors.items} />
            )}

            <div className="space-y-3">
              {formState.items.map((item, index) => (
                <div
                  key={item.id}
                  className="flex gap-3 items-start p-4 rounded-lg bg-muted/50"
                >
                  <div className="flex-1 space-y-3">
                    <div className="flex gap-3">
                      {/* 食品名 */}
                      <div className="flex-1 space-y-1">
                        <Label
                          htmlFor={`item-name-${item.id}`}
                          className="text-sm text-muted-foreground"
                        >
                          食品名
                        </Label>
                        <Input
                          id={`item-name-${item.id}`}
                          type="text"
                          value={item.name}
                          onChange={(e) =>
                            handleItemChange(item.id, "name", e.target.value)
                          }
                          placeholder={`食品${index + 1}`}
                          disabled={isLoading}
                          aria-invalid={!!formErrors.itemErrors?.[item.id]?.name}
                          aria-describedby={
                            formErrors.itemErrors?.[item.id]?.name
                              ? `item-name-error-${item.id}`
                              : undefined
                          }
                          className="h-10 bg-background"
                        />
                        {formErrors.itemErrors?.[item.id]?.name && (
                          <FieldError
                            id={`item-name-error-${item.id}`}
                            message={formErrors.itemErrors[item.id].name!}
                          />
                        )}
                      </div>

                      {/* カロリー */}
                      <div className="w-32 space-y-1">
                        <Label
                          htmlFor={`item-calories-${item.id}`}
                          className="text-sm text-muted-foreground"
                        >
                          カロリー (kcal)
                        </Label>
                        <Input
                          id={`item-calories-${item.id}`}
                          type="number"
                          min="0"
                          value={item.calories}
                          onChange={(e) =>
                            handleItemChange(
                              item.id,
                              "calories",
                              parseInt(e.target.value, 10) || 0
                            )
                          }
                          disabled={isLoading}
                          aria-invalid={!!formErrors.itemErrors?.[item.id]?.calories}
                          aria-describedby={
                            formErrors.itemErrors?.[item.id]?.calories
                              ? `item-calories-error-${item.id}`
                              : undefined
                          }
                          className="h-10 bg-background"
                        />
                        {formErrors.itemErrors?.[item.id]?.calories && (
                          <FieldError
                            id={`item-calories-error-${item.id}`}
                            message={formErrors.itemErrors[item.id].calories!}
                          />
                        )}
                      </div>
                    </div>
                  </div>

                  {/* 削除ボタン（アイテムが2つ以上の場合のみ表示） */}
                  {formState.items.length > 1 && (
                    <Button
                      type="button"
                      variant="ghost"
                      size="icon"
                      onClick={() => handleRemoveItem(item.id)}
                      disabled={isLoading}
                      className="mt-6 text-muted-foreground hover:text-destructive"
                      aria-label={`${item.name || `食品${index + 1}`}を削除`}
                    >
                      <TrashIcon className="w-4 h-4" />
                    </Button>
                  )}
                </div>
              ))}
            </div>
          </div>

          {/* 合計カロリー表示 */}
          <div className="flex items-center justify-between p-4 rounded-lg bg-primary/10 border border-primary/20">
            <span className="font-medium text-foreground">合計カロリー</span>
            <span className="text-2xl font-bold text-primary">
              {totalCalories.toLocaleString()} kcal
            </span>
          </div>

          {/* 送信ボタン */}
          <Button
            type="submit"
            className="w-full h-12 text-base font-medium bg-primary hover:bg-primary/90 transition-colors"
            disabled={isLoading}
          >
            {isLoading ? "記録中..." : "記録する"}
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
