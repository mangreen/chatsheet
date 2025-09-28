import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'

// https://vite.dev/config/
export default defineConfig({
  plugins: [svelte()],
  server: {
    // 確保這裡設定與後端 API 伺服器同步
    // 這樣前端才能正確地向後端發送請求 (http://localhost:8080)
    proxy: {
        '/api': {
            target: 'http://localhost:8080',
            changeOrigin: true,
            secure: false,
        }
    }
  }
})
