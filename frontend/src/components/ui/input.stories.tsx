/**
 * Input コンポーネントのStorybook
 * 各種入力タイプと状態の表示確認
 */
import type { Meta, StoryObj } from '@storybook/react-vite';
import { Input } from './input';
import { Label } from './label';

const meta: Meta<typeof Input> = {
  title: 'UI/Input',
  component: Input,
  tags: ['autodocs'],
  argTypes: {
    type: {
      control: 'select',
      options: ['text', 'email', 'password', 'number', 'date', 'search'],
      description: '入力タイプ',
    },
    disabled: {
      control: 'boolean',
      description: '無効化状態',
    },
    placeholder: {
      control: 'text',
      description: 'プレースホルダー',
    },
  },
};

export default meta;
type Story = StoryObj<typeof Input>;

/** テキスト入力 */
export const Text: Story = {
  args: {
    type: 'text',
    placeholder: 'テキストを入力',
  },
};

/** メールアドレス入力 */
export const Email: Story = {
  args: {
    type: 'email',
    placeholder: 'example@example.com',
  },
};

/** パスワード入力 */
export const Password: Story = {
  args: {
    type: 'password',
    placeholder: 'パスワードを入力',
  },
};

/** 数値入力 */
export const Number: Story = {
  args: {
    type: 'number',
    placeholder: '0',
  },
};

/** 日付入力 */
export const Date: Story = {
  args: {
    type: 'date',
  },
};

/** 無効化状態 */
export const Disabled: Story = {
  args: {
    type: 'text',
    placeholder: '入力不可',
    disabled: true,
  },
};

/** エラー状態 */
export const WithError: Story = {
  render: () => (
    <div className="space-y-2">
      <Label htmlFor="error-input">メールアドレス</Label>
      <Input
        id="error-input"
        type="email"
        placeholder="example@example.com"
        aria-invalid="true"
        className="border-destructive focus-visible:ring-destructive"
      />
      <p className="text-sm text-destructive">正しいメールアドレスを入力してください</p>
    </div>
  ),
};

/** ラベル付き */
export const WithLabel: Story = {
  render: () => (
    <div className="space-y-2">
      <Label htmlFor="labeled-input">ニックネーム</Label>
      <Input id="labeled-input" type="text" placeholder="ニックネームを入力" />
    </div>
  ),
};

/** 全タイプ一覧 */
export const AllTypes: Story = {
  render: () => (
    <div className="space-y-4 w-[300px]">
      <div className="space-y-2">
        <Label>テキスト</Label>
        <Input type="text" placeholder="テキスト" />
      </div>
      <div className="space-y-2">
        <Label>メール</Label>
        <Input type="email" placeholder="example@example.com" />
      </div>
      <div className="space-y-2">
        <Label>パスワード</Label>
        <Input type="password" placeholder="パスワード" />
      </div>
      <div className="space-y-2">
        <Label>数値</Label>
        <Input type="number" placeholder="0" />
      </div>
      <div className="space-y-2">
        <Label>日付</Label>
        <Input type="date" />
      </div>
    </div>
  ),
};
