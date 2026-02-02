/**
 * RecordForm - Storybookストーリー
 */
import type { Meta, StoryObj } from "@storybook/react-vite";
import { RecordForm } from "./RecordForm";

const meta: Meta<typeof RecordForm> = {
  title: "Features/Records/RecordForm",
  component: RecordForm,
  tags: ["autodocs"],
  parameters: {
    layout: "centered",
  },
  decorators: [
    (Story) => (
      <div className="w-[480px] p-4">
        <Story />
      </div>
    ),
  ],
};

export default meta;
type Story = StoryObj<typeof RecordForm>;

/** デフォルト表示（空の状態） */
export const Default: Story = {
  args: {},
};

/** 成功コールバック付き */
export const WithOnSuccess: Story = {
  args: {
    onSuccess: (response) => {
      console.log("記録成功:", response);
      alert("記録成功！");
    },
  },
};
