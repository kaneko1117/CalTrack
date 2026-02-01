/**
 * Select コンポーネントのStorybook
 * ドロップダウン選択の表示確認
 */
import type { Meta, StoryObj } from '@storybook/react-vite';
import { Select, SelectOption } from './select';
import { Label } from './label';

const meta: Meta<typeof Select> = {
  title: 'UI/Select',
  component: Select,
  tags: ['autodocs'],
  argTypes: {
    disabled: {
      control: 'boolean',
      description: '無効化状態',
    },
  },
};

export default meta;
type Story = StoryObj<typeof Select>;

/** 基本的なセレクト */
export const Default: Story = {
  render: () => (
    <Select>
      <SelectOption value="">選択してください</SelectOption>
      <SelectOption value="1">オプション1</SelectOption>
      <SelectOption value="2">オプション2</SelectOption>
      <SelectOption value="3">オプション3</SelectOption>
    </Select>
  ),
};

/** ラベル付き */
export const WithLabel: Story = {
  render: () => (
    <div className="space-y-2 w-[250px]">
      <Label htmlFor="select-with-label">カテゴリー</Label>
      <Select id="select-with-label">
        <SelectOption value="">選択してください</SelectOption>
        <SelectOption value="food">食品</SelectOption>
        <SelectOption value="drink">飲料</SelectOption>
        <SelectOption value="snack">お菓子</SelectOption>
      </Select>
    </div>
  ),
};

/** 性別選択 */
export const GenderSelect: Story = {
  render: () => (
    <div className="space-y-2 w-[250px]">
      <Label htmlFor="gender">性別</Label>
      <Select id="gender">
        <SelectOption value="">選択してください</SelectOption>
        <SelectOption value="male">男性</SelectOption>
        <SelectOption value="female">女性</SelectOption>
        <SelectOption value="other">その他</SelectOption>
      </Select>
    </div>
  ),
};

/** 活動レベル選択 */
export const ActivityLevelSelect: Story = {
  render: () => (
    <div className="space-y-2 w-[300px]">
      <Label htmlFor="activity">活動レベル</Label>
      <Select id="activity">
        <SelectOption value="">選択してください</SelectOption>
        <SelectOption value="sedentary">座りがち（運動なし）</SelectOption>
        <SelectOption value="light">軽い（週1-3回運動）</SelectOption>
        <SelectOption value="moderate">適度（週3-5回運動）</SelectOption>
        <SelectOption value="active">活動的（週6-7回運動）</SelectOption>
        <SelectOption value="veryActive">非常に活動的（毎日激しい運動）</SelectOption>
      </Select>
    </div>
  ),
};

/** 無効化状態 */
export const Disabled: Story = {
  render: () => (
    <div className="space-y-2 w-[250px]">
      <Label htmlFor="disabled-select">無効化</Label>
      <Select id="disabled-select" disabled>
        <SelectOption value="">選択できません</SelectOption>
        <SelectOption value="1">オプション1</SelectOption>
      </Select>
    </div>
  ),
};

/** エラー状態 */
export const WithError: Story = {
  render: () => (
    <div className="space-y-2 w-[250px]">
      <Label htmlFor="error-select">性別</Label>
      <Select
        id="error-select"
        aria-invalid="true"
        className="border-destructive focus-visible:ring-destructive"
      >
        <SelectOption value="">選択してください</SelectOption>
        <SelectOption value="male">男性</SelectOption>
        <SelectOption value="female">女性</SelectOption>
        <SelectOption value="other">その他</SelectOption>
      </Select>
      <p className="text-sm text-destructive">性別を選択してください</p>
    </div>
  ),
};

/** 選択済み */
export const Preselected: Story = {
  render: () => (
    <div className="space-y-2 w-[250px]">
      <Label htmlFor="preselected">性別</Label>
      <Select id="preselected" defaultValue="female">
        <SelectOption value="">選択してください</SelectOption>
        <SelectOption value="male">男性</SelectOption>
        <SelectOption value="female">女性</SelectOption>
        <SelectOption value="other">その他</SelectOption>
      </Select>
    </div>
  ),
};
