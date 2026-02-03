/**
 * CalorieChart - Storybookストーリー
 */
import type { Meta, StoryObj } from "@storybook/react-vite";
import { CalorieChart } from "./CalorieChart";
import type { DailyStatistic } from "../hooks/useStatistics";

const meta: Meta<typeof CalorieChart> = {
  title: "Features/Statistics/CalorieChart",
  component: CalorieChart,
  tags: ["autodocs"],
  parameters: {
    layout: "centered",
  },
  decorators: [
    (Story) => (
      <div className="w-[800px] p-4">
        <Story />
      </div>
    ),
  ],
};

export default meta;
type Story = StoryObj<typeof CalorieChart>;

/** モックデータ: 週間データ */
const mockWeekData: DailyStatistic[] = [
  { date: "2024-01-08", totalCalories: 1800 },
  { date: "2024-01-09", totalCalories: 2100 },
  { date: "2024-01-10", totalCalories: 1950 },
  { date: "2024-01-11", totalCalories: 2200 },
  { date: "2024-01-12", totalCalories: 1700 },
  { date: "2024-01-13", totalCalories: 2400 },
  { date: "2024-01-14", totalCalories: 1900 },
];

/** モックデータ: 目標を大きく下回る週間データ */
const mockLowCalorieData: DailyStatistic[] = [
  { date: "2024-01-08", totalCalories: 1200 },
  { date: "2024-01-09", totalCalories: 1100 },
  { date: "2024-01-10", totalCalories: 1350 },
  { date: "2024-01-11", totalCalories: 1400 },
  { date: "2024-01-12", totalCalories: 1250 },
  { date: "2024-01-13", totalCalories: 1300 },
  { date: "2024-01-14", totalCalories: 1150 },
];

/** デフォルト表示（週間データあり） */
export const Default: Story = {
  args: {
    data: mockWeekData,
    targetCalories: 2000,
    isLoading: false,
  },
};

/** ローディング状態 */
export const Loading: Story = {
  args: {
    data: [],
    targetCalories: 2000,
    isLoading: true,
  },
};

/** データなし（Empty状態） */
export const Empty: Story = {
  args: {
    data: [],
    targetCalories: 2000,
    isLoading: false,
  },
};

/** 目標を下回るデータ */
export const BelowTarget: Story = {
  args: {
    data: mockLowCalorieData,
    targetCalories: 2000,
    isLoading: false,
  },
};
