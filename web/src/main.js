import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'
import { setPenanganTanpaSesi } from './api/client'
import { useAuth } from './stores/auth'

// Font Plus Jakarta Sans di-bundle (offline-friendly, bukan CDN).
import '@fontsource/plus-jakarta-sans/400.css'
import '@fontsource/plus-jakarta-sans/500.css'
import '@fontsource/plus-jakarta-sans/600.css'
import '@fontsource/plus-jakarta-sans/700.css'
import '@fontsource/plus-jakarta-sans/800.css'

import './styles/tokens.css'

const app = createApp(App).use(createPinia()).use(router)

// Sesi bisa habis sementara aplikasi dibiarkan terbuka. Begitu server menjawab
// 401, bidan langsung diantar ke halaman masuk — tanpa ini yang terlihat
// hanyalah data yang diam-diam berhenti diperbarui.
setPenanganTanpaSesi(() => {
  useAuth().lupakan()
  if (router.currentRoute.value.name !== 'masuk') {
    router.replace({ name: 'masuk', query: { next: router.currentRoute.value.fullPath } })
  }
})

app.mount('#app')
