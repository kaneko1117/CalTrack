/**
 * FormField コンポーネントのテスト
 */
import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { Formik, Form, Field, FieldProps } from "formik";
import { FormField } from "./FormField";

/**
 * FormFieldをFormikでラップしてレンダリングするヘルパー
 */
function renderFormField({
  initialValue = "",
  initialTouched = false,
  initialError = "",
  label = "テストラベル",
  type = "text",
  placeholder = "",
  min,
}: {
  initialValue?: string | number;
  initialTouched?: boolean;
  initialError?: string;
  label?: string;
  type?: string;
  placeholder?: string;
  min?: number;
} = {}) {
  const onSubmit = vi.fn();

  return render(
    <Formik
      initialValues={{ testField: initialValue }}
      initialTouched={{ testField: initialTouched }}
      initialErrors={initialError ? { testField: initialError } : {}}
      onSubmit={onSubmit}
    >
      <Form>
        <Field
          name="testField"
          component={FormField}
          label={label}
          type={type}
          placeholder={placeholder}
          min={min}
        />
      </Form>
    </Formik>
  );
}

/**
 * モックフォームオブジェクトでFormFieldを直接レンダリングするヘルパー
 * isSubmitting状態のテスト用
 */
function renderFormFieldWithMockedForm({
  isSubmitting = false,
  touched = false,
  error = "",
  value = "",
  label = "テストラベル",
}: {
  isSubmitting?: boolean;
  touched?: boolean;
  error?: string;
  value?: string;
  label?: string;
} = {}) {
  const mockField: FieldProps["field"] = {
    name: "testField",
    value,
    onChange: vi.fn(),
    onBlur: vi.fn(),
  };

  const mockMeta: FieldProps["meta"] = {
    touched,
    error,
    value,
    initialTouched: false,
    initialValue: "",
    initialError: undefined,
  };

  const mockForm = {
    isSubmitting,
    getFieldMeta: () => mockMeta,
  } as unknown as FieldProps["form"];

  return render(
    <FormField field={mockField} form={mockForm} meta={mockMeta} label={label} />
  );
}

describe("FormField", () => {
  describe("レンダリング", () => {
    it("ラベルが表示される", () => {
      renderFormField({ label: "テストラベル" });

      expect(screen.getByText("テストラベル")).toBeInTheDocument();
    });

    it("入力値が反映される", () => {
      renderFormField({ initialValue: "初期値" });

      const input = screen.getByRole("textbox");
      expect(input).toHaveValue("初期値");
    });

    it("プレースホルダーが表示される", () => {
      renderFormField({ placeholder: "テストプレースホルダー" });

      const input = screen.getByPlaceholderText("テストプレースホルダー");
      expect(input).toBeInTheDocument();
    });

    it("数値フィールドが正しく表示される", () => {
      renderFormField({ type: "number", initialValue: 100 });

      const input = screen.getByRole("spinbutton");
      expect(input).toHaveValue(100);
    });

    it("数値フィールドで0の場合は空欄として表示される", () => {
      renderFormField({ type: "number", initialValue: 0 });

      const input = screen.getByRole("spinbutton");
      expect(input).toHaveValue(null);
    });
  });

  describe("エラー表示", () => {
    it("エラーがある場合にエラーメッセージが表示される", () => {
      renderFormField({
        initialTouched: true,
        initialError: "エラーメッセージ",
      });

      expect(screen.getByText("エラーメッセージ")).toBeInTheDocument();
    });

    it("エラーがない場合はエラーメッセージが表示されない", () => {
      renderFormField({ initialTouched: true });

      expect(screen.queryByRole("alert")).not.toBeInTheDocument();
      expect(screen.queryByText(/エラー/)).not.toBeInTheDocument();
    });

    it("touchedでない場合はエラーメッセージが表示されない", () => {
      renderFormField({
        initialTouched: false,
        initialError: "エラーメッセージ",
      });

      expect(screen.queryByText("エラーメッセージ")).not.toBeInTheDocument();
    });

    it("エラーがある場合にaria-invalidがtrueになる", () => {
      renderFormField({
        initialTouched: true,
        initialError: "エラーメッセージ",
      });

      const input = screen.getByRole("textbox");
      expect(input).toHaveAttribute("aria-invalid", "true");
    });

    it("エラーがない場合にaria-invalidがfalseになる", () => {
      renderFormField({ initialTouched: true });

      const input = screen.getByRole("textbox");
      expect(input).toHaveAttribute("aria-invalid", "false");
    });
  });

  describe("disabled状態", () => {
    it("disabled状態が反映される", () => {
      renderFormFieldWithMockedForm({ isSubmitting: true });

      const input = screen.getByRole("textbox");
      expect(input).toBeDisabled();
    });

    it("送信中でない場合は有効化される", () => {
      renderFormFieldWithMockedForm({ isSubmitting: false });

      const input = screen.getByRole("textbox");
      expect(input).not.toBeDisabled();
    });
  });

  describe("ユーザー操作", () => {
    it("入力値が変更できる", async () => {
      const user = userEvent.setup();
      renderFormField();

      const input = screen.getByRole("textbox");
      await user.type(input, "新しい値");

      expect(input).toHaveValue("新しい値");
    });

    it("数値フィールドで値を入力できる", async () => {
      const user = userEvent.setup();
      renderFormField({ type: "number" });

      const input = screen.getByRole("spinbutton");
      await user.type(input, "250");

      expect(input).toHaveValue(250);
    });
  });

  describe("アクセシビリティ", () => {
    it("ラベルとinputが正しく関連付けられる", () => {
      renderFormField({ label: "テストラベル" });

      const input = screen.getByLabelText("テストラベル");
      expect(input).toBeInTheDocument();
    });

    it("エラーメッセージがaria-describedbyで関連付けられる", () => {
      renderFormField({
        initialTouched: true,
        initialError: "エラーメッセージ",
      });

      const input = screen.getByRole("textbox");
      expect(input).toHaveAttribute("aria-describedby", "testField-error");

      const errorMessage = screen.getByText("エラーメッセージ");
      expect(errorMessage.closest("[id]")).toHaveAttribute("id", "testField-error");
    });
  });
});
