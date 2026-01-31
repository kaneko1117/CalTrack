/**
 * Selectコンポーネント
 * shadcn/uiスタイルのドロップダウン選択
 * シンプルなnative select実装（Radix UIなしで動作）
 */
import * as React from "react"

import { cn } from "@/lib/utils"

/** Selectコンポーネントのプロパティ */
export type SelectProps = React.SelectHTMLAttributes<HTMLSelectElement>

const Select = React.forwardRef<HTMLSelectElement, SelectProps>(
  ({ className, children, ...props }, ref) => {
    return (
      <select
        className={cn(
          "flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
          className
        )}
        ref={ref}
        {...props}
      >
        {children}
      </select>
    )
  }
)
Select.displayName = "Select"

/** SelectOptionコンポーネントのプロパティ */
export type SelectOptionProps = React.OptionHTMLAttributes<HTMLOptionElement>

const SelectOption = React.forwardRef<HTMLOptionElement, SelectOptionProps>(
  ({ className, ...props }, ref) => {
    return <option ref={ref} className={className} {...props} />
  }
)
SelectOption.displayName = "SelectOption"

export { Select, SelectOption }
