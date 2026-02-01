import type { StorybookConfig } from '@storybook/react-vite';

const config: StorybookConfig = {
  stories: [
    // features配下のストーリーを対象
    '../src/features/**/*.mdx',
    '../src/features/**/*.stories.@(js|jsx|mjs|ts|tsx)',
    // pages配下のストーリーを対象（必要に応じて）
    '../src/pages/**/*.mdx',
    '../src/pages/**/*.stories.@(js|jsx|mjs|ts|tsx)',
  ],
  addons: [
    '@storybook/addon-a11y',
    '@storybook/addon-docs',
  ],
  framework: '@storybook/react-vite',
};

export default config;