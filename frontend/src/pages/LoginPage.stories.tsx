/**
 * LoginPage コンポーネントのStorybook
 * ログインページ全体の表示確認
 */
import type { Meta, StoryObj } from '@storybook/react-vite';
import { LoginPage } from './LoginPage';

const meta: Meta<typeof LoginPage> = {
  title: 'Pages/LoginPage',
  component: LoginPage,
  tags: ['autodocs'],
  parameters: {
    layout: 'fullscreen',
  },
};

export default meta;
type Story = StoryObj<typeof LoginPage>;

/** デフォルト状態 */
export const Default: Story = {
  args: {},
};

/** カスタムリダイレクト先 */
export const WithCustomRedirect: Story = {
  args: {
    redirectTo: '/dashboard',
  },
};
