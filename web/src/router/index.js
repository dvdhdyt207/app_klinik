import { createRouter, createWebHistory } from 'vue-router'
import Patient from '../views/Patient.vue'
import BidanApp from '../views/BidanApp.vue'
import Login from '../views/Login.vue'
import { useAuth } from '../stores/auth'

const routes = [
  { path: '/', name: 'patient', component: Patient },       // halaman publik pasien
  { path: '/masuk', name: 'masuk', component: Login },
  { path: '/app', name: 'app', component: BidanApp, meta: { wajibLogin: true } },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Penjaga rute — ini KENYAMANAN, bukan keamanan.
//
// Yang benar-benar menjaga data pasien adalah middleware di server: tanpa
// cookie sesi yang sah, /api/state menjawab 401 dan tidak ada satu baris pun
// yang keluar. Penjaga di sini hanya supaya bidan diantar ke halaman masuk,
// bukan disuguhi layar yang gagal memuat. Jangan pernah memindahkan keputusan
// "boleh lihat atau tidak" ke sini — semua kode frontend bisa dilewati.
router.beforeEach(async (to) => {
  if (!to.meta.wajibLogin) return true
  const auth = useAuth()
  if (await auth.cek()) return true
  return { name: 'masuk', query: { next: to.fullPath } }
})

export default router
