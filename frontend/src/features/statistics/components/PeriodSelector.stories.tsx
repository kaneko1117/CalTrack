/**
 * PeriodSelector - Storybookストーリー
 */
import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { PeriodSelector } from "./PeriodSelector";
import type { Period } from "../hooks/useStatistics";

const meta: Meta<typeof PeriodSelector> = {
  title: "Features/Statistics/PeriodSelector",
  component: PeriodSelector,
  tags: ["autodocs"],
  parameters: {
    layout: "centered",
  },
};

export default meta;
type Story = StoryObj<typeof PeriodSelector>;

/** インタラクティブラッパー */
function InteractivePeriodSelector({ initialValue }: { initialValue: Period }) {
  const [value, setValue] = useState<Period>(initialValue);
  return <PeriodSelector value={value} onChange={setValue} />;
}

/** デフォルト表示（週間選択） */
export const Default: Story = {
  render: () => <InteractivePeriodSelector initialValue="week" />,
};

/** 月間選択状態 */
export const MonthSelected: Story = {
  render: () => <InteractivePeriodSelector initialValue="month" />,
};
