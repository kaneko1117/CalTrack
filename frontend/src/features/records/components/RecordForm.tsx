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
import { post } from "@/lib/api";
import type { ApiErrorResponse } from "@/lib/api";
import { newItemName } from "@/domain/valueObjects/itemName";
import { newCalories } from "@/domain/valueObjects/calories";
import { newEatenAt } from "@/domain/valueObjects/eatenAt";

/** 記録作成レスポンス */
export type CreateRecordResponse = {
  recordId: string;
  eatenAt: string;
  totalCalories: number;
};

/** 記録作成リクエスト */
type CreateRecordRequest = {
  eatenAt: string;
  items: Array<{ name: string; calories: number }>;
};

/** 記録作成API */
const createRecord = (data: CreateRecordRequest) =>
  post<CreateRecordResponse>("/api/v1/records", data);

/** RecordFormコンポーネントのProps */
export type RecordFormProps = {
  /** 記録作成成功時のコールバック */
  onSuccess?: (response: CreateRecordResponse) => void;
};

/** 食品アイテムの内部状態（IDを含む） */
type ItemState = {
  id: string;
  name: string;
  calories: number;
};

/** フォームフィールド型 */
type FormField = "eatenAt";

/** フォームの状態 */
type FormState = {
  eatenAt: string;
};

/** 新規アイテム入力の状態 */
type NewItemState = {
  name: string;
  calories: string;
};

/** 新規アイテムのエラー */
type NewItemErrors = {
  name: string | null;
  calories: string | null;
};

/** エラーの初期状態 */
const initialErrors: Record<FormField, string | null> = {
  eatenAt: null,
};

/** 新規アイテムの初期状態 */
const initialNewItemState: NewItemState = {
  name: "",
  calories: "",
};

/** 新規アイテムエラーの初期状態 */
const initialNewItemErrors: NewItemErrors = {
  name: null,
  calories: null,
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
  // フォーム状態（eatenAtはuseFormで管理）
  const [formState, setFormState] = useState<FormState>({
    eatenAt: getCurrentDateTimeLocal(),
  });
  const [errors, setErrors] = useState<Record<FormField, string | null>>(
    initialErrors
  );

  // アイテムリスト（動的配列のためuseState）
  const [items, setItems] = useState<ItemState[]>([]);
  const [itemsError, setItemsError] = useState<string | null>(null);

  // 新規アイテム入力
  const [newItem, setNewItem] = useState<NewItemState>(initialNewItemState);
  const [newItemErrors, setNewItemErrors] =
    useState<NewItemErrors>(initialNewItemErrors);

  // API状態
  const [isLoading, setIsLoading] = useState(false);
  const [apiError, setApiError] = useState<ApiErrorResponse | null>(null);

  /** 合計カロリーの計算（メモ化） */
  const totalCalories = useMemo(() => {
    return items.reduce((sum, item) => sum + item.calories, 0);
  }, [items]);

  /**
   * 食事日時の変更ハンドラ
   */
  const handleEatenAtChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setFormState((prev) => ({ ...prev, eatenAt: value }));

    // バリデーション
    if (!value) {
      setErrors((prev) => ({
        ...prev,
        eatenAt: "食事日時を入力してください",
      }));
      return;
    }

    const date = new Date(value);
    const result = newEatenAt(date);
    if (!result.ok) {
      setErrors((prev) => ({ ...prev, eatenAt: result.error.message }));
    } else {
      setErrors((prev) => ({ ...prev, eatenAt: null }));
    }

    if (apiError) {
      setApiError(null);
    }
  };

  /**
   * 新規アイテム名の変更ハンドラ
   */
  const handleNewItemNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setNewItem((prev) => ({ ...prev, name: value }));

    // リアルタイムバリデーション
    const result = newItemName(value);
    if (!result.ok) {
      setNewItemErrors((prev) => ({ ...prev, name: result.error.message }));
    } else {
      setNewItemErrors((prev) => ({ ...prev, name: null }));
    }
  };

  /**
   * 新規アイテムカロリーの変更ハンドラ
   */
  const handleNewItemCaloriesChange = (
    e: React.ChangeEvent<HTMLInputElement>
  ) => {
    const value = e.target.value;
    setNewItem((prev) => ({ ...prev, calories: value }));

    // リアルタイムバリデーション
    const num = parseInt(value, 10);
    if (isNaN(num) || value === "") {
      setNewItemErrors((prev) => ({
        ...prev,
        calories: "カロリーを入力してください",
      }));
      return;
    }

    const result = newCalories(num);
    if (!result.ok) {
      setNewItemErrors((prev) => ({ ...prev, calories: result.error.message }));
    } else {
      setNewItemErrors((prev) => ({ ...prev, calories: null }));
    }
  };

  /**
   * アイテム追加ハンドラ
   */
  const handleAddItem = () => {
    // バリデーション
    const nameResult = newItemName(newItem.name);
    const caloriesNum = parseInt(newItem.calories, 10);
    const caloriesResult = newCalories(isNaN(caloriesNum) ? 0 : caloriesNum);

    if (!nameResult.ok || !caloriesResult.ok) {
      setNewItemErrors({
        name: nameResult.ok ? null : nameResult.error.message,
        calories: caloriesResult.ok ? null : caloriesResult.error.message,
      });
      return;
    }

    // アイテム追加
    const newItemState: ItemState = {
      id: generateId(),
      name: nameResult.value.value,
      calories: caloriesResult.value.value,
    };
    setItems((prev) => [...prev, newItemState]);

    // 入力をリセット
    setNewItem(initialNewItemState);
    setNewItemErrors(initialNewItemErrors);

    // アイテムエラーをクリア
    if (itemsError) {
      setItemsError(null);
    }
  };

  /**
   * アイテム削除ハンドラ
   */
  const handleRemoveItem = (itemId: string) => {
    setItems((prev) => prev.filter((item) => item.id !== itemId));
  };

  /**
   * フォーム送信ハンドラ
   */
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // バリデーション
    let hasError = false;

    // eatenAtバリデーション
    if (!formState.eatenAt) {
      setErrors((prev) => ({
        ...prev,
        eatenAt: "食事日時を入力してください",
      }));
      hasError = true;
    } else {
      const date = new Date(formState.eatenAt);
      const eatenAtResult = newEatenAt(date);
      if (!eatenAtResult.ok) {
        setErrors((prev) => ({ ...prev, eatenAt: eatenAtResult.error.message }));
        hasError = true;
      }
    }

    // itemsバリデーション
    if (items.length === 0) {
      setItemsError("少なくとも1つの食品を追加してください");
      hasError = true;
    }

    if (hasError) {
      return;
    }

    // datetime-localをISO 8601形式に変換
    const eatenAtISO = new Date(formState.eatenAt).toISOString();

    // API呼び出し
    setIsLoading(true);
    setApiError(null);

    try {
      const response = await createRecord({
        eatenAt: eatenAtISO,
        items: items.map(({ name, calories }) => ({ name, calories })),
      });

      // 成功時にフォームをリセット
      setFormState({ eatenAt: getCurrentDateTimeLocal() });
      setErrors(initialErrors);
      setItems([]);
      setItemsError(null);
      setNewItem(initialNewItemState);
      setNewItemErrors(initialNewItemErrors);

      onSuccess?.(response);
    } catch (error) {
      setApiError(error as ApiErrorResponse);
    } finally {
      setIsLoading(false);
    }
  };

  /** フォームが有効かどうか */
  const isFormValid = useMemo(() => {
    return !errors.eatenAt && items.length > 0;
  }, [errors.eatenAt, items.length]);

  /** 追加ボタンが有効かどうか */
  const isAddButtonValid = useMemo(() => {
    return !newItemErrors.name && !newItemErrors.calories && newItem.name && newItem.calories;
  }, [newItemErrors.name, newItemErrors.calories, newItem.name, newItem.calories]);

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* APIエラー表示 */}
      {apiError && (
        <div
          className="flex items-start gap-3 p-4 text-sm rounded-lg bg-destructive/10 border border-destructive/20"
          role="alert"
        >
          <AlertCircleIcon className="w-5 h-5 text-destructive flex-shrink-0 mt-0.5" />
          <div className="flex-1">
            <p className="font-medium text-destructive">
              {getErrorMessage(apiError.code)}
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
          aria-invalid={!!errors.eatenAt}
          aria-describedby={errors.eatenAt ? "eatenAt-error" : undefined}
          className="h-11 bg-background border-input focus:border-primary focus:ring-primary/20"
        />
        {errors.eatenAt && (
          <FieldError id="eatenAt-error" message={errors.eatenAt} />
        )}
      </div>

      {/* 食品アイテム追加セクション */}
      <div className="space-y-4">
        <Label className="text-foreground font-medium">食品を追加</Label>

        <div className="flex gap-3 items-end">
          {/* 食品名入力 */}
          <div className="flex-1 space-y-1">
            <Label
              htmlFor="new-item-name"
              className="text-sm text-muted-foreground"
            >
              食品名
            </Label>
            <Input
              id="new-item-name"
              type="text"
              value={newItem.name}
              onChange={handleNewItemNameChange}
              placeholder="例: ご飯"
              disabled={isLoading}
              aria-invalid={!!newItemErrors.name}
              aria-describedby={
                newItemErrors.name ? "new-item-name-error" : undefined
              }
              className="h-10 bg-background"
            />
          </div>

          {/* カロリー入力 */}
          <div className="w-28 space-y-1">
            <Label
              htmlFor="new-item-calories"
              className="text-sm text-muted-foreground"
            >
              kcal
            </Label>
            <Input
              id="new-item-calories"
              type="number"
              min="1"
              value={newItem.calories}
              onChange={handleNewItemCaloriesChange}
              placeholder="300"
              disabled={isLoading}
              aria-invalid={!!newItemErrors.calories}
              aria-describedby={
                newItemErrors.calories ? "new-item-calories-error" : undefined
              }
              className="h-10 bg-background"
            />
          </div>

          {/* 追加ボタン */}
          <Button
            type="button"
            variant="outline"
            size="icon"
            onClick={handleAddItem}
            disabled={isLoading || !isAddButtonValid}
            className="h-10 w-10"
            aria-label="食品を追加"
          >
            <PlusIcon className="w-4 h-4" />
          </Button>
        </div>

        {/* 新規アイテムのエラー表示 */}
        {newItemErrors.name && (
          <FieldError id="new-item-name-error" message={newItemErrors.name} />
        )}
        {newItemErrors.calories && (
          <FieldError
            id="new-item-calories-error"
            message={newItemErrors.calories}
          />
        )}
      </div>

      {/* 追加済み食品一覧 */}
      {items.length > 0 && (
        <div className="space-y-3">
          <Label className="text-foreground font-medium">
            追加済み ({items.length}件)
          </Label>
          <div className="space-y-2">
            {items.map((item) => (
              <div
                key={item.id}
                className="flex items-center justify-between p-3 rounded-lg bg-muted/50"
              >
                <div className="flex-1 min-w-0">
                  <p className="font-medium truncate">{item.name}</p>
                  <p className="text-sm text-muted-foreground">
                    {item.calories.toLocaleString()} kcal
                  </p>
                </div>
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  onClick={() => handleRemoveItem(item.id)}
                  disabled={isLoading}
                  className="text-muted-foreground hover:text-destructive"
                  aria-label={`${item.name}を削除`}
                >
                  <TrashIcon className="w-4 h-4" />
                </Button>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* アイテムエラー */}
      {itemsError && <FieldError id="items-error" message={itemsError} />}

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
        disabled={!isFormValid || isLoading}
      >
        {isLoading ? "記録中..." : "記録する"}
      </Button>
    </form>
  );
}
