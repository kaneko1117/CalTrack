/**
 * RecordForm - カロリー記録フォームコンポーネント
 * Formik の Field, Form, FieldArray を使用
 */
import { Formik, Form, Field, FieldArray } from "formik";
import * as yup from "yup";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { FormField } from "@/components/form";
import { post } from "@/lib/api";
import { newEatenAt } from "@/domain/valueObjects/eatenAt";
import { newItemName } from "@/domain/valueObjects/itemName";
import { newCalories, sumCalories } from "@/domain/valueObjects/calories";
import { newRecord } from "@/domain/entities/record";

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

/** 記録作成レスポンス */
type CreateRecordResponse = {
  recordId: string;
  eatenAt: string;
  totalCalories: number;
};

/** 記録作成API */
const createRecord = (data: { eatenAt: string; items: Array<{ name: string; calories: number }> }) =>
  post<CreateRecordResponse>("/api/v1/records", data);

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

export function RecordForm({ onSuccess }: RecordFormProps) {
  const initialValues: RecordFormValues = {
    eatenAt: getInitialEatenAt(),
    items: [createEmptyItem()],
  };

  return (
    <Formik
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
          const response = await createRecord({
            eatenAt: record.eatenAt.toISOString(),
            items: record.items.map((item) => ({
              name: item.name.value,
              calories: item.calories.value,
            })),
          });

          resetForm({
            values: {
              eatenAt: getInitialEatenAt(),
              items: [createEmptyItem()],
            },
          });
          onSuccess?.(response);
        } catch {
          setSubmitting(false);
        }
      }}
    >
      {({ values, errors, isSubmitting, isValid }) => {
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
                              disabled={isSubmitting}
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
                      disabled={isSubmitting}
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
                disabled={isSubmitting || !isValid}
              >
                {isSubmitting ? "記録中..." : "記録する"}
              </Button>
            </div>
          </Form>
        );
      }}
    </Formik>
  );
}
