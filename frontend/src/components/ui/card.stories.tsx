/**
 * Card コンポーネントのStorybook
 * カード全体および各サブコンポーネントの表示確認
 */
import type { Meta, StoryObj } from '@storybook/react-vite';
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
  CardFooter,
} from './card';
import { Button } from './button';

const meta: Meta<typeof Card> = {
  title: 'UI/Card',
  component: Card,
  tags: ['autodocs'],
};

export default meta;
type Story = StoryObj<typeof Card>;

/** 基本的なカード */
export const Default: Story = {
  render: () => (
    <Card className="w-[350px]">
      <CardHeader>
        <CardTitle>カードタイトル</CardTitle>
        <CardDescription>カードの説明文がここに入ります</CardDescription>
      </CardHeader>
      <CardContent>
        <p>カードの本文コンテンツです。</p>
      </CardContent>
      <CardFooter>
        <Button>アクション</Button>
      </CardFooter>
    </Card>
  ),
};

/** ヘッダーのみ */
export const HeaderOnly: Story = {
  render: () => (
    <Card className="w-[350px]">
      <CardHeader>
        <CardTitle>タイトルのみ</CardTitle>
      </CardHeader>
    </Card>
  ),
};

/** ヘッダー + 説明文 */
export const WithDescription: Story = {
  render: () => (
    <Card className="w-[350px]">
      <CardHeader>
        <CardTitle>タイトル</CardTitle>
        <CardDescription>
          これはカードの説明文です。補足情報を記載します。
        </CardDescription>
      </CardHeader>
    </Card>
  ),
};

/** コンテンツのみ */
export const ContentOnly: Story = {
  render: () => (
    <Card className="w-[350px]">
      <CardContent className="pt-6">
        <p>コンテンツのみのカードです。</p>
      </CardContent>
    </Card>
  ),
};

/** フッター付きカード */
export const WithFooter: Story = {
  render: () => (
    <Card className="w-[350px]">
      <CardHeader>
        <CardTitle>確認</CardTitle>
        <CardDescription>この操作を実行しますか？</CardDescription>
      </CardHeader>
      <CardContent>
        <p>操作の詳細説明がここに入ります。</p>
      </CardContent>
      <CardFooter className="flex justify-end gap-2">
        <Button variant="outline">キャンセル</Button>
        <Button>確認</Button>
      </CardFooter>
    </Card>
  ),
};

/** フォームカード例 */
export const FormCard: Story = {
  render: () => (
    <Card className="w-[400px]">
      <CardHeader>
        <CardTitle>ログイン</CardTitle>
        <CardDescription>アカウントにログインしてください</CardDescription>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <label className="text-sm font-medium">メールアドレス</label>
          <input
            type="email"
            placeholder="example@example.com"
            className="w-full h-10 px-3 border rounded-md"
          />
        </div>
        <div className="space-y-2">
          <label className="text-sm font-medium">パスワード</label>
          <input
            type="password"
            placeholder="パスワード"
            className="w-full h-10 px-3 border rounded-md"
          />
        </div>
      </CardContent>
      <CardFooter>
        <Button className="w-full">ログイン</Button>
      </CardFooter>
    </Card>
  ),
};
