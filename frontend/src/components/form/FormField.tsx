/**
 * FormField - Formik Field と shadcn/ui Input を統合したコンポーネント
 */
import { FieldProps } from "formik";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

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

export type FormFieldProps = FieldProps & {
  label: string;
  type?: string;
  placeholder?: string;
  min?: number;
};

export function FormField({
  field,
  form,
  label,
  type = "text",
  placeholder,
  min,
}: FormFieldProps) {
  const meta = form.getFieldMeta(field.name);
  const hasError = meta.touched && meta.error;

  return (
    <div className="space-y-1">
      <Label htmlFor={field.name} className="text-sm">{label}</Label>
      <Input
        {...field}
        id={field.name}
        type={type}
        placeholder={placeholder}
        min={min}
        disabled={form.isSubmitting}
        aria-invalid={!!hasError}
        aria-describedby={hasError ? `${field.name}-error` : undefined}
        className="h-10 bg-background"
        value={type === "number" && field.value === 0 ? "" : field.value}
      />
      {hasError && (
        <p id={`${field.name}-error`} className="flex items-center gap-1.5 text-sm text-destructive mt-1">
          <AlertCircleIcon className="w-4 h-4 flex-shrink-0" />
          <span>{meta.error}</span>
        </p>
      )}
    </div>
  );
}
