/**
 * ProfileEditForm - Storybookストーリー
 * プロフィール編集フォームの各状態の表示確認
 */
import type { Meta, StoryObj } from "@storybook/react-vite";
import { SWRConfig } from "swr";
import { ProfileEditForm } from "./ProfileEditForm";
import type { CurrentUserResponse } from "../api";

const meta: Meta<typeof ProfileEditForm> = {
  title: "Features/User/ProfileEditForm",
  component: ProfileEditForm,
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
type Story = StoryObj<typeof ProfileEditForm>;

/** モックユーザーデータ */
const mockUser: CurrentUserResponse = {
  email: "test@example.com",
  nickname: "テストユーザー",
  weight: 70,
  height: 175,
  birthDate: "1990-01-01",
  gender: "male",
  activityLevel: "moderate",
};

/** デフォルト表示（ユーザーデータ読み込み完了後） */
export const Default: Story = {
  decorators: [
    (Story) => (
      <SWRConfig
        value={{
          dedupingInterval: 0,
          provider: () => new Map(),
          fetcher: async () => mockUser,
        }}
      >
        <Story />
      </SWRConfig>
    ),
  ],
  args: {
    onSuccess: (response) => {
      console.log("更新成功:", response);
      alert("プロフィールを更新しました！");
    },
  },
};

/** Loading状態（データ読み込み中） */
export const Loading: Story = {
  decorators: [
    (Story) => (
      <SWRConfig
        value={{
          dedupingInterval: 0,
          provider: () => new Map(),
          fetcher: async () => {
            // 永遠に解決しないPromiseでローディング状態を再現
            return new Promise(() => {});
          },
        }}
      >
        <Story />
      </SWRConfig>
    ),
  ],
  args: {},
};

/** FetchError状態（データ取得エラー） */
export const FetchError: Story = {
  decorators: [
    (Story) => (
      <SWRConfig
        value={{
          dedupingInterval: 0,
          provider: () => new Map(),
          fetcher: async () => {
            throw { code: "INTERNAL_ERROR", message: "サーバーエラーが発生しました" };
          },
        }}
      >
        <Story />
      </SWRConfig>
    ),
  ],
  args: {},
};
