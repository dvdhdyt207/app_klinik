<script setup>
import { ref, computed, onMounted, useTemplateRef } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuth } from '../stores/auth'
import { api } from '../api/client'
import { monogram } from '../lib/format'

const router = useRouter()
const route = useRoute()
const auth = useAuth()

// Sengaja KOSONG. Dulu berisi 'bidan' supaya bidan tinggal mengisi sandi —
// tapi halaman ini terbuka ke internet, dan kolom yang sudah terisi berarti
// menyebutkan nama penggunanya kepada siapa pun yang membukanya; yang menebak
// tinggal mengurus sandinya saja.
//
// Kemudahannya tidak hilang: kolomnya bertanda autocomplete="username", jadi
// peramban bidan sendiri tetap mengisikannya setelah sekali masuk.
const username = ref('')
const password = ref('')
const galat = ref('')
const mengirim = ref(false)
const kolomNama = useTemplateRef('kolomNama')
const kolomSandi = useTemplateRef('kolomSandi')

// Nama klinik dulu ditulis mati "Bidan Pit" di sini. Diambil dari
// /api/public/status — satu-satunya sumber yang terbuka tanpa sesi, yang
// memang tujuannya dibaca sebelum siapa pun masuk. Gagal mengambilnya tidak
// menghalangi login: kolom sandi tetap bisa diisi, hanya namanya jadi umum.
const namaKlinik = ref('')
const judul = computed(() => namaKlinik.value || 'Klinik')
const inisial = computed(() => monogram(judul.value))

onMounted(async () => {
  // Kursor jatuh ke kolom yang memang masih kosong. Selama nama pengguna
  // terisi bawaan, langsung ke sandi selalu benar; sekarang tidak lagi — dan
  // memaksa bidan mengetuk kolom di atasnya tiap kali masuk adalah persis
  // gangguan yang tidak boleh ada di pintu rekam medisnya sendiri.
  if (username.value) kolomSandi.value?.focus()
  else kolomNama.value?.focus()
  try {
    const st = await api.publicStatus()
    namaKlinik.value = st?.clinic || ''
  } catch { /* nama klinik cuma hiasan di layar ini — jangan halangi login */ }
})

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
    <div class="kolom-tengah">
      <!-- Tanda pengenal klinik di luar kartu: sebelumnya layar ini tidak
           memuat satu pun penanda visual, sehingga terlihat sama saja dengan
           halaman masuk aplikasi mana pun. -->
      <header class="merek">
        <div class="logo" aria-hidden="true">{{ inisial }}</div>
        <h1 class="judul">{{ judul }}</h1>
        <p class="sub">Masuk untuk membuka catatan klinik.</p>
      </header>

      <form class="masuk-kartu" @submit.prevent="kirim">
        <label class="label" for="u">Nama pengguna</label>
        <input
          id="u"
          ref="kolomNama"
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

        <button class="tombol" type="submit" :disabled="mengirim || !username || !password">
          {{ mengirim ? 'Memeriksa…' : 'Masuk' }}
        </button>
      </form>

      <router-link class="tautan" to="/">← Lihat halaman status bidan</router-link>
    </div>
  </div>
</template>

<style scoped>
.masuk-bg {
  min-height: 100dvh;
  display: flex; align-items: center; justify-content: center;
  padding: 24px;
  background: var(--patient-bg);
  /* Cahaya teal samar di belakang kartu. Tanpa ini layar masuk di monitor
     lebar hanyalah satu kartu kecil di tengah bidang datar yang luas. */
  background-image:
    radial-gradient(circle at 50% 0%, rgba(18, 128, 106, .10), transparent 60%),
    radial-gradient(circle at 50% 100%, rgba(18, 128, 106, .06), transparent 55%);
}
.kolom-tengah {
  width: 100%; max-width: 400px;
  display: flex; flex-direction: column;
}

.merek { text-align: center; margin-bottom: 22px; }
.logo {
  width: 56px; height: 56px; margin: 0 auto 14px; border-radius: 18px;
  background: linear-gradient(145deg, var(--accent) 0%, var(--accent-press) 100%);
  color: #fff; font-weight: 800; font-size: 19px;
  display: flex; align-items: center; justify-content: center;
  box-shadow: 0 10px 24px rgba(18, 128, 106, .28);
}
.judul { font-size: 25px; font-weight: 800; letter-spacing: -.02em; color: var(--ink); margin: 0; line-height: 1.2; overflow-wrap: anywhere; }
.sub { font-size: 13.5px; color: var(--text-secondary); margin: 6px 0 0; }

.masuk-kartu {
  background: var(--card); border: 1px solid var(--border);
  border-radius: var(--r-lg); padding: 24px;
  box-shadow: var(--shadow-md);
  display: flex; flex-direction: column;
}
.label { font-size: 12px; font-weight: 700; color: var(--label); margin-bottom: 6px; }
.kolom {
  border: 1px solid var(--input-border); border-radius: 10px;
  padding: 12px; font: inherit; font-size: 15px; color: var(--ink);
  background: var(--fill2); margin-bottom: 14px; width: 100%;
}
.kolom:focus { outline: none; border-color: var(--accent); background: var(--card); box-shadow: 0 0 0 3px var(--accent-soft); }
.galat {
  background: var(--danger-bg); color: var(--danger);
  font-size: 13px; font-weight: 600;
  border-radius: 10px; padding: 10px 12px; margin: 0 0 14px;
}
/* Pesan galat datang apa adanya dari server dan diawali huruf kecil
   ("nama pengguna atau sandi salah"). Dibesarkan di sini, bukan di JavaScript:
   ini urusan tampilan, dan teks aslinya tetap utuh untuk log. */
.galat::first-letter { text-transform: uppercase; }
.tombol {
  background: var(--accent); color: #fff; border: 0;
  border-radius: var(--r-md); padding: 13px; font: inherit; font-size: 15px; font-weight: 700;
  cursor: pointer; transition: background .15s ease;
  margin-top: 2px;
}
.tombol:hover:not(:disabled) { background: var(--accent-hover); }
.tombol:active:not(:disabled) { background: var(--accent-press); }
.tombol:disabled { background: var(--disabled); cursor: default; }
.tautan {
  margin-top: 20px; text-align: center;
  font-size: 13px; font-weight: 600; color: var(--muted); text-decoration: none;
}
.tautan:hover { color: var(--accent); }
</style>
