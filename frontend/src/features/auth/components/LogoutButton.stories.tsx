/**
 * LogoutButton コンポーネントのStorybook
 */
import type { Meta, StoryObj } from "@storybook/react-vite";
import { LogoutButton } from "./LogoutButton";

const meta: Meta<typeof LogoutButton> = {
  title: "Features/Auth/LogoutButton",
  component: LogoutButton,
  tags: ["autodocs"],
  parameters: {
    layout: "centered",
  },
};

export default meta;
type Story = StoryObj<typeof LogoutButton>;

/** デフォルト状態 */
export const Default: Story = {};

/** 成功コールバック付き */
export const WithOnSuccess: Story = {
  args: {
    onSuccess: () => {
      console.log("ログアウト成功");
      alert("ログアウトしました！");
    },
  },
};
