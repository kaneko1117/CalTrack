/**
 * TodaySummary - Storybookストーリー
 */
import type { Meta, StoryObj } from "@storybook/react-vite";
import { TodaySummary } from "./TodaySummary";
import type { TodayRecordsResponse } from "../hooks/useTodayRecords";
import type { ApiErrorResponse } from "@/lib/api";

const meta: Meta<typeof TodaySummary> = {
  title: "Features/Records/TodaySummary",
  component: TodaySummary,
  tags: ["autodocs"],
  parameters: {
    layout: "centered",
  },
  decorators: [
    (Story) => (
      <div className="w-[720px] p-4">
        <Story />
      </div>
    ),
  ],
};

export default meta;
type Story = StoryObj<typeof TodaySummary>;

/** モックデータ: 記録あり */
const mockDataWithRecords: TodayRecordsResponse = {
  date: "2024-01-15",
  totalCalories: 1200,
  targetCalories: 2000,
  difference: -800,
  records: [
    {
      id: "record-1",
      eatenAt: "2024-01-15T08:30:00Z",
      items: [
        { itemId: "item-1", name: "トースト", calories: 200 },
        { itemId: "item-2", name: "目玉焼き", calories: 150 },
      ],
    },
    {
      id: "record-2",
      eatenAt: "2024-01-15T12:00:00Z",
      items: [
        { itemId: "item-3", name: "チキンカレー", calories: 650 },
        { itemId: "item-4", name: "サラダ", calories: 50 },
      ],
    },
    {
      id: "record-3",
      eatenAt: "2024-01-15T15:30:00Z",
      items: [{ itemId: "item-5", name: "プロテインバー", calories: 150 }],
    },
  ],
};

/** モックデータ: 記録なし */
const mockDataEmpty: TodayRecordsResponse = {
  date: "2024-01-15",
  totalCalories: 0,
  targetCalories: 2000,
  difference: -2000,
  records: [],
};

/** モックデータ: 目標超過 */
const mockDataOverTarget: TodayRecordsResponse = {
  date: "2024-01-15",
  totalCalories: 2500,
  targetCalories: 2000,
  difference: 500,
  records: [
    {
      id: "record-1",
      eatenAt: "2024-01-15T08:00:00Z",
      items: [
        { itemId: "item-1", name: "モーニングセット", calories: 600 },
      ],
    },
    {
      id: "record-2",
      eatenAt: "2024-01-15T12:00:00Z",
      items: [
        { itemId: "item-2", name: "ラーメン大盛り", calories: 900 },
        { itemId: "item-3", name: "餃子", calories: 300 },
      ],
    },
    {
      id: "record-3",
      eatenAt: "2024-01-15T19:00:00Z",
      items: [
        { itemId: "item-4", name: "焼肉定食", calories: 700 },
      ],
    },
  ],
};

/** モックエラー */
const mockError: ApiErrorResponse = {
  code: "INTERNAL_ERROR",
  message: "サーバーエラーが発生しました",
};

/** デフォルト表示（記録あり） */
export const Default: Story = {
  args: {
    data: mockDataWithRecords,
    isPending: false,
    error: null,
  },
};

/** ローディング状態 */
export const Loading: Story = {
  args: {
    data: null,
    isPending: true,
    error: null,
  },
};

/** エラー状態 */
export const Error: Story = {
  args: {
    data: null,
    isPending: false,
    error: mockError,
  },
};

/** 記録が0件 */
export const Empty: Story = {
  args: {
    data: mockDataEmpty,
    isPending: false,
    error: null,
  },
};

/** 目標超過状態 */
export const OverTarget: Story = {
  args: {
    data: mockDataOverTarget,
    isPending: false,
    error: null,
  },
};
