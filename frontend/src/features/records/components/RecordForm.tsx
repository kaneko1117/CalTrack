/**
 * RecordForm - カロリー記録フォームコンポーネント
 * Formik の Field, Form, FieldArray を使用
 * 画像解析機能付き
 */
import { useCallback, useState, useRef } from "react";
import { Formik, Form, Field, FieldArray, type FormikProps } from "formik";
import * as yup from "yup";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { FormField } from "@/components/form";
import { useRequestMutation } from "@/features/common/hooks";
import { newEatenAt } from "@/domain/valueObjects/eatenAt";
import { newItemName } from "@/domain/valueObjects/itemName";
import { newCalories, sumCalories } from "@/domain/valueObjects/calories";
import { newRecord } from "@/domain/entities/record";
import { ImageInput } from "./ImageInput";
import { useImageAnalysis, type ImageAnalysisResponse } from "../hooks";
import type { ImageFile } from "@/domain/valueObjects/imageFile";

/** フォームの食品アイテム型 */
type RecordFormItem = {
  name: string;
  calories: string;
};

/** フォームの値型 */
type RecordFormValues = {
  eatenAt: string;
  items: RecordFormItem[];
};

/** 記録作成リクエスト */
type CreateRecordRequest = {
  eatenAt: string;
  items: Array<{ name: string; calories: number }>;
};

/** 記録作成レスポンス */
type CreateRecordResponse = {
  recordId: string;
  eatenAt: string;
  totalCalories: number;
};

/** RecordFormのProps */
export type RecordFormProps = {
  onSuccess?: (response: CreateRecordResponse) => void;
};

/** 空のアイテム */
const createEmptyItem = (): RecordFormItem => ({
  name: "",
  calories: "",
});

/** 現在日時をdatetime-local形式で取得 */
const getInitialEatenAt = (): string => {
  const now = new Date();
  const result = newEatenAt(now.toISOString(), now);
  return result.ok ? result.value.toDateTimeLocal() : "";
};

/** yupバリデーションスキーマ（VOでバリデーション） */
const recordFormSchema = yup.object().shape({
  eatenAt: yup.string().test("vo-validation", "", function (value) {
    const result = newEatenAt(value ?? "", new Date());
    if (!result.ok) {
      return this.createError({ message: result.error.message });
    }
    return true;
  }),
  items: yup.array().of(
    yup.object().shape({
      name: yup.string().test("vo-validation", "", function (value) {
        const result = newItemName(value ?? "");
        if (!result.ok) {
          return this.createError({ message: result.error.message });
        }
        return true;
      }),
      calories: yup.string().test("vo-validation", "", function (value) {
        const result = newCalories(value ?? "");
        if (!result.ok) {
          return this.createError({ message: result.error.message });
        }
        return true;
      }),
    }),
  ),
});

/** Plusアイコン */
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

/** Trashアイコン */
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
      <polyline points="3 6 5 6 21 6" />
      <path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2" />
    </svg>
  );
}

/** AlertCircleアイコン */
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

/** 解析結果を食品リストに適用するヘルパー */
const applyAnalysisToItems = (
  analysisData: ImageAnalysisResponse,
  currentItems: RecordFormItem[]
): RecordFormItem[] => {
  // 解析結果を RecordFormItem 形式に変換
  const analyzedItems: RecordFormItem[] = analysisData.items.map((item) => ({
    name: item.name,
    calories: String(item.calories),
  }));

  // 既存の空でないアイテムを取得
  const nonEmptyItems = currentItems.filter(
    (item) => item.name.trim() !== "" || item.calories.trim() !== ""
  );

  // 既存の非空アイテム + 解析結果を結合
  const mergedItems = [...nonEmptyItems, ...analyzedItems];

  // 結果が空の場合は空のアイテムを1つ追加
  return mergedItems.length > 0 ? mergedItems : [createEmptyItem()];
};

export function RecordForm({ onSuccess }: RecordFormProps) {
  const initialValues: RecordFormValues = {
    eatenAt: getInitialEatenAt(),
    items: [createEmptyItem()],
  };

  // 画像プレビュー用のstate
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);

  // Formikのインスタンス参照を保持（useRefで無限ループを防ぐ）
  const formikRef = useRef<FormikProps<RecordFormValues> | null>(null);

  // 画像解析成功時のコールバック
  const handleAnalysisSuccess = useCallback(
    (data: ImageAnalysisResponse) => {
      if (formikRef.current && data.items.length > 0) {
        const newItems = applyAnalysisToItems(data, formikRef.current.values.items);
        formikRef.current.setFieldValue("items", newItems);
      }
    },
    []
  );

  // 画像解析フック
  const { analyze, isAnalyzing, error: analysisError, reset: resetAnalysis } = useImageAnalysis({
    onSuccess: handleAnalysisSuccess,
  });

  // 画像選択時のハンドラ
  const handleImageSelect = useCallback(
    async (imageFile: ImageFile) => {
      setPreviewUrl(imageFile.dataUrl);
      await analyze(imageFile);
    },
    [analyze]
  );

  // 画像クリア時のハンドラ
  const handleImageClear = useCallback(() => {
    setPreviewUrl(null);
    resetAnalysis();
  }, [resetAnalysis]);

  // SWRベースのmutationフック
  const { trigger, isMutating } = useRequestMutation<
    CreateRecordResponse,
    CreateRecordRequest
  >("/api/v1/records", "POST");

  return (
    <Formik
      innerRef={formikRef}
      initialValues={initialValues}
      validationSchema={recordFormSchema}
      validateOnBlur={true}
      validateOnChange={true}
      validateOnMount={true}
      onSubmit={async (values, { setSubmitting, resetForm }) => {
        // Record Entityで最終バリデーション
        // number型input要素の場合、Formikがcaloriesをnumberに変換する可能性があるため、
        // 明示的にstring型に変換してからnewRecordに渡す
        const recordResult = newRecord({
          eatenAt: values.eatenAt,
          items: values.items.map((item) => ({
            name: item.name,
            calories: String(item.calories),
          })),
        });

        if (!recordResult.ok) {
          setSubmitting(false);
          return;
        }

        const record = recordResult.value;

        try {
          const response = await trigger({
            eatenAt: record.eatenAt.toISOString(),
            items: record.items.map((item) => ({
              name: item.name.value,
              calories: item.calories.value,
            })),
          });

          // 画像プレビューもクリア
          setPreviewUrl(null);
          resetAnalysis();

          resetForm({
            values: {
              eatenAt: getInitialEatenAt(),
              items: [createEmptyItem()],
            },
          });

          // trigger成功後にonSuccessを呼ぶ
          onSuccess?.(response);
        } catch {
          setSubmitting(false);
        }
      }}
    >
      {({ values, errors, isValid }) => {
        const totalCalories = sumCalories(
          values.items.map((item) => Number(item.calories) || 0),
        );
        const itemsError =
          typeof errors.items === "string" ? errors.items : undefined;

        return (
          <Form className="space-y-6">
            {/* 食事日時 */}
            <Field
              name="eatenAt"
              component={FormField}
              label="食事日時"
              type="datetime-local"
            />

            {/* 画像解析セクション */}
            <div className="space-y-2">
              <Label className="text-foreground font-medium">
                画像から入力 (オプション)
              </Label>
              <ImageInput
                onImageSelect={handleImageSelect}
                isAnalyzing={isAnalyzing}
                error={analysisError}
                previewUrl={previewUrl}
                disabled={isMutating}
                onClear={handleImageClear}
              />
              <p className="text-xs text-muted-foreground">
                食事の写真をアップロードすると、AIが食品とカロリーを自動入力します
              </p>
            </div>

            {/* 食品リスト */}
            <div className="space-y-4">
              <div className="flex items-center justify-between">
                <Label className="text-foreground font-medium">
                  食品リスト
                </Label>
                <span className="text-sm text-muted-foreground">
                  合計: {totalCalories.toLocaleString()} kcal
                </span>
              </div>

              {itemsError && (
                <p className="flex items-center gap-1.5 text-sm text-destructive">
                  <AlertCircleIcon className="w-4 h-4 flex-shrink-0" />
                  <span>{itemsError}</span>
                </p>
              )}

              <FieldArray name="items">
                {({ push, remove }) => (
                  <div className="space-y-4">
                    {values.items.map((_, index) => (
                      <div
                        key={index}
                        className="p-4 border rounded-lg bg-muted/30 space-y-3"
                      >
                        <div className="flex items-center justify-between">
                          <span className="text-sm font-medium text-muted-foreground">
                            食品 {index + 1}
                          </span>
                          {values.items.length > 1 && (
                            <Button
                              type="button"
                              variant="ghost"
                              size="sm"
                              onClick={() => remove(index)}
                              disabled={isMutating}
                              className="h-8 w-8 p-0 text-muted-foreground hover:text-destructive"
                              aria-label={`食品 ${index + 1} を削除`}
                            >
                              <TrashIcon className="w-4 h-4" />
                            </Button>
                          )}
                        </div>

                        <Field
                          name={`items.${index}.name`}
                          component={FormField}
                          label="食品名"
                          placeholder="例: 白米"
                        />

                        <Field
                          name={`items.${index}.calories`}
                          component={FormField}
                          label="カロリー (kcal)"
                          type="number"
                          placeholder="例: 250"
                          min={1}
                        />
                      </div>
                    ))}

                    <Button
                      type="button"
                      variant="outline"
                      onClick={() => push(createEmptyItem())}
                      disabled={isMutating}
                      className="w-full h-10 border-dashed"
                    >
                      <PlusIcon className="w-4 h-4 mr-2" />
                      食品を追加
                    </Button>
                  </div>
                )}
              </FieldArray>
            </div>

            {/* 送信ボタン */}
            <div className="!mt-6">
              <Button
                type="submit"
                className="w-full h-12 text-base font-medium"
                disabled={isMutating || !isValid || isAnalyzing}
              >
                {isMutating ? "記録中..." : "記録する"}
              </Button>
            </div>
          </Form>
        );
      }}
    </Formik>
  );
}
