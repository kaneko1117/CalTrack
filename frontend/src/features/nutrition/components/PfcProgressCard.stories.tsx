import type { Meta, StoryObj } from "@storybook/react-vite";
import { PfcProgressCard } from "./PfcProgressCard";
import type { ApiErrorResponse } from "@/lib/api";

const meta: Meta<typeof PfcProgressCard> = {
  title: "Features/Nutrition/PfcProgressCard",
  component: PfcProgressCard,
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
type Story = StoryObj<typeof PfcProgressCard>;

/** デフォルト表示（通常状態 0-79%） */
export const Default: Story = {
  args: {
    data: {
      date: "2026-02-09T00:00:00Z",
      current: { protein: 35.0, fat: 25.0, carbs: 100.0 },
      target: { protein: 120.0, fat: 65.0, carbs: 300.0 },
    },
    isLoading: false,
    error: null,
  },
};

/** 適切状態（80-100%） */
export const Optimal: Story = {
  args: {
    data: {
      date: "2026-02-09T00:00:00Z",
      current: { protein: 100.0, fat: 55.0, carbs: 260.0 },
      target: { protein: 120.0, fat: 65.0, carbs: 300.0 },
    },
    isLoading: false,
    error: null,
  },
};

/** 超過状態（100%超） */
export const Over: Story = {
  args: {
    data: {
      date: "2026-02-09T00:00:00Z",
      current: { protein: 200.0, fat: 100.0, carbs: 500.0 },
      target: { protein: 120.0, fat: 65.0, carbs: 300.0 },
    },
    isLoading: false,
    error: null,
  },
};

/** 混合状態（各栄養素で異なるステータス） */
export const Mixed: Story = {
  args: {
    data: {
      date: "2026-02-09T00:00:00Z",
      current: { protein: 35.0, fat: 55.0, carbs: 500.0 },
      target: { protein: 120.0, fat: 65.0, carbs: 300.0 },
    },
    isLoading: false,
    error: null,
  },
};

/** ローディング状態 */
export const Loading: Story = {
  args: {
    data: null,
    isLoading: true,
    error: null,
  },
};

/** エラー状態 */
export const Error: Story = {
  args: {
    data: null,
    isLoading: false,
    error: {
      code: "INTERNAL_ERROR",
      message: "サーバーエラーが発生しました",
    } as ApiErrorResponse,
  },
};

/** データなし */
export const Empty: Story = {
  args: {
    data: null,
    isLoading: false,
    error: null,
  },
};
