/**
 * Button コンポーネントのStorybook
 * 各バリアント・サイズ・状態の表示確認
 */
import type { Meta, StoryObj } from '@storybook/react-vite';
import { Button } from './button';

const meta: Meta<typeof Button> = {
  title: 'UI/Button',
  component: Button,
  tags: ['autodocs'],
  argTypes: {
    variant: {
      control: 'select',
      options: ['default', 'destructive', 'outline', 'secondary', 'ghost', 'link'],
      description: 'ボタンのスタイルバリアント',
    },
    size: {
      control: 'select',
      options: ['default', 'sm', 'lg', 'icon'],
      description: 'ボタンのサイズ',
    },
    disabled: {
      control: 'boolean',
      description: '無効化状態',
    },
  },
};

export default meta;
type Story = StoryObj<typeof Button>;

/** デフォルトボタン */
export const Default: Story = {
  args: {
    children: 'ボタン',
    variant: 'default',
    size: 'default',
  },
};

/** 破壊的アクション用ボタン */
export const Destructive: Story = {
  args: {
    children: '削除',
    variant: 'destructive',
  },
};

/** アウトラインボタン */
export const Outline: Story = {
  args: {
    children: 'アウトライン',
    variant: 'outline',
  },
};

/** セカンダリボタン */
export const Secondary: Story = {
  args: {
    children: 'セカンダリ',
    variant: 'secondary',
  },
};

/** ゴーストボタン */
export const Ghost: Story = {
  args: {
    children: 'ゴースト',
    variant: 'ghost',
  },
};

/** リンクスタイルボタン */
export const Link: Story = {
  args: {
    children: 'リンク',
    variant: 'link',
  },
};

/** 小サイズボタン */
export const Small: Story = {
  args: {
    children: '小サイズ',
    size: 'sm',
  },
};

/** 大サイズボタン */
export const Large: Story = {
  args: {
    children: '大サイズ',
    size: 'lg',
  },
};

/** 無効化状態 */
export const Disabled: Story = {
  args: {
    children: '無効',
    disabled: true,
  },
};

/** 全バリアント一覧 */
export const AllVariants: Story = {
  render: () => (
    <div className="flex flex-wrap gap-4">
      <Button variant="default">Default</Button>
      <Button variant="destructive">Destructive</Button>
      <Button variant="outline">Outline</Button>
      <Button variant="secondary">Secondary</Button>
      <Button variant="ghost">Ghost</Button>
      <Button variant="link">Link</Button>
    </div>
  ),
};

/** 全サイズ一覧 */
export const AllSizes: Story = {
  render: () => (
    <div className="flex items-center gap-4">
      <Button size="sm">Small</Button>
      <Button size="default">Default</Button>
      <Button size="lg">Large</Button>
    </div>
  ),
};
