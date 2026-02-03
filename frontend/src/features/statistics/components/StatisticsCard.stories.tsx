/**
 * StatisticsCard - Storybookストーリー
 */
import type { Meta, StoryObj } from "@storybook/react-vite";
import { StatisticsCard } from "./StatisticsCard";
import type { StatisticsResponse } from "../hooks/useStatistics";

const meta: Meta<typeof StatisticsCard> = {
  title: "Features/Statistics/StatisticsCard",
  component: StatisticsCard,
  tags: ["autodocs"],
  parameters: {
    layout: "centered",
  },
  decorators: [
    (Story) => (
      <div className="w-[900px] p-4">
        <Story />
      </div>
    ),
  ],
};

export default meta;
type Story = StoryObj<typeof StatisticsCard>;

/** モックデータ: 良好な達成状況 */
const mockGoodData: StatisticsResponse = {
  period: "week",
  targetCalories: 2000,
  averageCalories: 1850,
  totalDays: 7,
  achievedDays: 5,
  overDays: 2,
  dailyStatistics: [],
};

/** モックデータ: 超過多め */
const mockOverData: StatisticsResponse = {
  period: "week",
  targetCalories: 2000,
  averageCalories: 2300,
  totalDays: 7,
  achievedDays: 2,
  overDays: 5,
  dailyStatistics: [],
};

/** モックデータ: 完璧な達成 */
const mockPerfectData: StatisticsResponse = {
  period: "week",
  targetCalories: 2000,
  averageCalories: 1900,
  totalDays: 7,
  achievedDays: 7,
  overDays: 0,
  dailyStatistics: [],
};

/** モックデータ: データなし（初期状態） */
const mockEmptyData: StatisticsResponse = {
  period: "week",
  targetCalories: 2000,
  averageCalories: 0,
  totalDays: 0,
  achievedDays: 0,
  overDays: 0,
  dailyStatistics: [],
};

/** デフォルト表示（良好な達成状況） */
export const Default: Story = {
  args: {
    data: mockGoodData,
  },
};

/** 超過多めの状態 */
export const OverTarget: Story = {
  args: {
    data: mockOverData,
  },
};

/** 完璧な達成状態 */
export const Perfect: Story = {
  args: {
    data: mockPerfectData,
  },
};

/** データなし（Empty状態） */
export const Empty: Story = {
  args: {
    data: mockEmptyData,
  },
};
