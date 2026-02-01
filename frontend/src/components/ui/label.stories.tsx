/**
 * Label コンポーネントのStorybook
 * フォームラベルの表示確認
 */
import type { Meta, StoryObj } from '@storybook/react-vite';
import { Label } from './label';
import { Input } from './input';

const meta: Meta<typeof Label> = {
  title: 'UI/Label',
  component: Label,
  tags: ['autodocs'],
};

export default meta;
type Story = StoryObj<typeof Label>;

/** 基本的なラベル */
export const Default: Story = {
  args: {
    children: 'ラベル',
  },
};

/** 入力フィールドと組み合わせ */
export const WithInput: Story = {
  render: () => (
    <div className="space-y-2">
      <Label htmlFor="email">メールアドレス</Label>
      <Input id="email" type="email" placeholder="example@example.com" />
    </div>
  ),
};

/** 必須項目 */
export const Required: Story = {
  render: () => (
    <div className="space-y-2">
      <Label htmlFor="name">
        名前 <span className="text-destructive">*</span>
      </Label>
      <Input id="name" type="text" placeholder="名前を入力" />
    </div>
  ),
};

/** 無効化されたフィールドのラベル */
export const DisabledField: Story = {
  render: () => (
    <div className="space-y-2">
      <Label htmlFor="disabled-input" className="peer-disabled:opacity-70">
        無効化されたフィールド
      </Label>
      <Input id="disabled-input" type="text" disabled placeholder="入力不可" />
    </div>
  ),
};

/** 複数フィールド */
export const MultipleFields: Story = {
  render: () => (
    <div className="space-y-4">
      <div className="space-y-2">
        <Label htmlFor="username">ユーザー名</Label>
        <Input id="username" type="text" placeholder="ユーザー名" />
      </div>
      <div className="space-y-2">
        <Label htmlFor="password">パスワード</Label>
        <Input id="password" type="password" placeholder="パスワード" />
      </div>
    </div>
  ),
};
