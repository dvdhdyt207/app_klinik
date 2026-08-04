<script setup>
import { ref, onMounted, useTemplateRef } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuth } from '../stores/auth'

const router = useRouter()
const route = useRoute()
const auth = useAuth()

const username = ref('bidan')
const password = ref('')
const galat = ref('')
const mengirim = ref(false)
const kolomSandi = useTemplateRef('kolomSandi')

onMounted(() => kolomSandi.value?.focus())

async function kirim() {
  if (mengirim.value) return
  galat.value = ''
  mengirim.value = true
  try {
    await auth.masuk(username.value, password.value)
    // Kembali ke halaman yang tadi dituju, kalau ada. Hanya jalur internal
    // yang diterima — nilai ?next= datang dari URL dan bisa diisi siapa saja;
    // tanpa saringan ini, tautan "masuk" bisa dipakai melempar bidan ke situs
    // palsu tepat setelah ia mengetikkan sandinya.
    const next = String(route.query.next || '')
    router.replace(next.startsWith('/') && !next.startsWith('//') ? next : '/app')
  } catch (e) {
    galat.value = e.message || 'Gagal masuk'
    password.value = ''
    kolomSandi.value?.focus()
  } finally {
    mengirim.value = false
  }
}
</script>

<template>
  <div class="masuk-bg">
    <form class="masuk-kartu" @submit.prevent="kirim">
      <div class="eyebrow">KLINIK</div>
      <h1 class="judul">Bidan Pit</h1>
      <p class="sub">Masuk untuk membuka catatan klinik.</p>

      <label class="label" for="u">Nama pengguna</label>
      <input
        id="u"
        v-model.trim="username"
        class="kolom"
        autocomplete="username"
        autocapitalize="none"
        spellcheck="false"
        required
      />

      <label class="label" for="p">Kata sandi</label>
      <input
        id="p"
        ref="kolomSandi"
        v-model="password"
        class="kolom"
        type="password"
        autocomplete="current-password"
        required
      />

      <p v-if="galat" class="galat" role="alert">{{ galat }}</p>

      <button class="tombol" type="submit" :disabled="mengirim || !password">
        {{ mengirim ? 'Memeriksa…' : 'Masuk' }}
      </button>

      <router-link class="tautan" to="/">← Lihat halaman status bidan</router-link>
    </form>
  </div>
</template>

<style scoped>
.masuk-bg {
  min-height: 100dvh;
  display: flex; align-items: center; justify-content: center;
  background: var(--app-bg); padding: 24px;
}
.masuk-kartu {
  width: 100%; max-width: 380px;
  background: var(--card); border: 1px solid var(--border);
  border-radius: 18px; padding: 28px 24px;
  box-shadow: 0 8px 28px rgba(22, 32, 46, .06);
  display: flex; flex-direction: column;
}
.eyebrow { font-size: 11px; font-weight: 800; letter-spacing: 1.2px; color: var(--muted); }
.judul { font-size: 24px; font-weight: 800; color: var(--ink); margin: 2px 0 4px; }
.sub { font-size: 13px; color: var(--text-secondary); margin: 0 0 20px; }
.label { font-size: 12px; font-weight: 700; color: var(--label); margin-bottom: 6px; }
.kolom {
  border: 1px solid var(--input-border); border-radius: 10px;
  padding: 11px 12px; font: inherit; font-size: 15px; color: var(--ink);
  background: var(--fill2); margin-bottom: 14px; width: 100%;
}
.kolom:focus { outline: 2px solid var(--accent); outline-offset: 1px; background: var(--card); }
.galat {
  background: var(--danger-bg); color: var(--danger);
  font-size: 13px; font-weight: 600;
  border-radius: 10px; padding: 9px 11px; margin: 0 0 14px;
}
.tombol {
  background: var(--accent); color: #fff; border: 0;
  border-radius: 10px; padding: 12px; font: inherit; font-size: 15px; font-weight: 700;
  cursor: pointer;
}
.tombol:hover:not(:disabled) { background: var(--accent-hover); }
.tombol:disabled { background: var(--disabled); cursor: default; }
.tautan {
  margin-top: 18px; text-align: center;
  font-size: 13px; color: var(--muted); text-decoration: none;
}
.tautan:hover { color: var(--accent); }
</style>
