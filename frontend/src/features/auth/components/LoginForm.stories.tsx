/**
 * LoginForm コンポーネントのStorybook
 * ログインフォームの各状態の表示確認
 */
import type { Meta, StoryObj } from '@storybook/react-vite';
import { BrowserRouter } from 'react-router-dom';
import { LoginForm } from './LoginForm';

const meta: Meta<typeof LoginForm> = {
  title: 'Features/Auth/LoginForm',
  component: LoginForm,
  tags: ['autodocs'],
  parameters: {
    layout: 'centered',
  },
  decorators: [
    (Story) => (
      <BrowserRouter>
        <div className="w-[400px]">
          <Story />
        </div>
      </BrowserRouter>
    ),
  ],
};

export default meta;
type Story = StoryObj<typeof LoginForm>;

/** デフォルト状態 */
export const Default: Story = {
  args: {},
};

/** 成功コールバック付き */
export const WithOnSuccess: Story = {
  args: {
    onSuccess: (response) => {
      console.log('ログイン成功:', response);
      alert('ログイン成功！');
    },
  },
};
