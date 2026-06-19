import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

export default defineConfig({
  plugins: [react(), tailwindcss()],
  test: {
    environment: 'jsdom',
    setupFiles: './src/test/setup.js',
    environmentOptions: {
      jsdom: {
        // jsdom 28 requires a URL for Storage API (localStorage/sessionStorage) to work
        url: 'http://localhost:5173',
      },
    },
  },
})
