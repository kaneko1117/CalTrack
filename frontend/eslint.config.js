import js from '@eslint/js';
import globals from 'globals';
import reactHooks from 'eslint-plugin-react-hooks';
import reactRefresh from 'eslint-plugin-react-refresh';
import tseslint from 'typescript-eslint';

export default tseslint.config(
  // 除外設定（ビルド成果物、依存関係、キャッシュ）
  { ignores: ['dist', 'node_modules', '.vite'] },

  // 基本設定
  {
    extends: [js.configs.recommended, ...tseslint.configs.recommended],
    files: ['**/*.{ts,tsx}'],
    languageOptions: {
      ecmaVersion: 2020,
      globals: globals.browser,
    },
    plugins: {
      'react-hooks': reactHooks,
      'react-refresh': reactRefresh,
    },
    rules: {
      // React Hooksルール
      ...reactHooks.configs.recommended.rules,

      // React Refreshルール（shadcn/uiのvariantsエクスポートを許容）
      'react-refresh/only-export-components': [
        'warn',
        {
          allowConstantExport: true,
          allowExportNames: ['buttonVariants'],
        },
      ],

      // TypeScript厳格ルール: any禁止（CLAUDE.mdの規約対応）
      '@typescript-eslint/no-explicit-any': 'error',

      // 未使用変数の警告
      '@typescript-eslint/no-unused-vars': [
        'error',
        {
          argsIgnorePattern: '^_',
          varsIgnorePattern: '^_',
        },
      ],
    },
  }
);
