/**
 * RegisterPage コンポーネントのStorybook
 * 新規登録ページ全体の表示確認
 */
import type { Meta, StoryObj } from '@storybook/react-vite';
import { RegisterPage } from './RegisterPage';

const meta: Meta<typeof RegisterPage> = {
  title: 'Pages/RegisterPage',
  component: RegisterPage,
  tags: ['autodocs'],
  parameters: {
    layout: 'fullscreen',
  },
};

export default meta;
type Story = StoryObj<typeof RegisterPage>;

/** デフォルト状態 */
export const Default: Story = {
  args: {},
};

/** カスタムリダイレクト先 */
export const WithCustomRedirect: Story = {
  args: {
    redirectTo: '/login',
  },
};
