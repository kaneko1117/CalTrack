/// <reference types="vitest" />
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

// allowedHosts を環境変数から取得
// VITE_ALLOWED_HOSTS: カンマ区切りで複数指定可能
function getAllowedHosts(): string[] | true {
  const hostsEnv = process.env.VITE_ALLOWED_HOSTS
  if (!hostsEnv) {
    return ['localhost']
  }
  if (hostsEnv.toLowerCase() === 'all') {
    return true
  }
  return hostsEnv.split(',').map(host => host.trim()).filter(host => host !== '')
}

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    host: '0.0.0.0',
    port: 5173,
    allowedHosts: getAllowedHosts(),
    watch: {
      usePolling: true,
    },
  },
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: ['./src/test/setup.ts'],
  },
})
