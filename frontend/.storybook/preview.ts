import type { Preview } from '@storybook/react-vite';

// グローバルCSSの読み込み（Tailwind CSS）
import '../src/index.css';

const preview: Preview = {
  parameters: {
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i,
      },
    },
  },
};

export default preview;