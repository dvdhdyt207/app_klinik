import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// Saat development, Vite dev server (:5173) mem-proxy /api ke server Go (:4000),
// sehingga frontend memakai URL relatif '/api/...' baik di dev maupun produksi
// (di produksi Go men-serve build ini di origin yang sama).
export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    proxy: {
      '/api': 'http://localhost:4000',
    },
  },
})
