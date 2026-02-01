/**
 * RegisterForm コンポーネントのStorybook
 * 新規登録フォームの各状態の表示確認
 */
import type { Meta, StoryObj } from '@storybook/react-vite';
import { RegisterForm } from './RegisterForm';

const meta: Meta<typeof RegisterForm> = {
  title: 'Auth/RegisterForm',
  component: RegisterForm,
  tags: ['autodocs'],
  parameters: {
    layout: 'centered',
  },
  decorators: [
    (Story) => (
      <div className="w-[450px]">
        <Story />
      </div>
    ),
  ],
};

export default meta;
type Story = StoryObj<typeof RegisterForm>;

/** デフォルト状態 */
export const Default: Story = {
  args: {},
};

/** 成功コールバック付き */
export const WithOnSuccess: Story = {
  args: {
    onSuccess: () => {
      console.log('登録成功');
      alert('登録成功！');
    },
  },
};
