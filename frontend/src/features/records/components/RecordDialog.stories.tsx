/**
 * RecordDialog - Storybookストーリー
 */
import type { Meta, StoryObj } from "@storybook/react-vite";
import { RecordDialog } from "./RecordDialog";

const meta: Meta<typeof RecordDialog> = {
  title: "Features/Records/RecordDialog",
  component: RecordDialog,
  tags: ["autodocs"],
  parameters: {
    layout: "centered",
  },
};

export default meta;
type Story = StoryObj<typeof RecordDialog>;

/** デフォルト表示（ダイアログ閉じた状態） */
export const Default: Story = {
  args: {},
};

/** 成功コールバック付き */
export const WithOnSuccess: Story = {
  args: {
    onSuccess: () => {
      console.log("記録成功");
      alert("記録成功！");
    },
  },
};
