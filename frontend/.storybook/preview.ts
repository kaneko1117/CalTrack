import type { Preview } from '@storybook/react-vite';
import { MemoryRouter } from 'react-router-dom';
import React from 'react';

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
  // React Router対応のデコレータ
  decorators: [
    (Story) => (
      React.createElement(MemoryRouter, null, React.createElement(Story))
    ),
  ],
};

export default preview;