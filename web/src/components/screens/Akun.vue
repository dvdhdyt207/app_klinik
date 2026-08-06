<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../../stores/auth'
import { api } from '../../api/client'

const auth = useAuth()
const router = useRouter()

// Sama dengan SandiMinimal di server (internal/api/akun.go). Diperiksa di sini
// hanya supaya bidan tahu sebelum menekan tombol — server tetap memeriksanya
// sendiri, karena pemeriksaan di browser bisa dilewati siapa pun.
const SANDI_MINIMAL = 10

const lama = ref('')
const baru = ref('')
const ulang = ref('')
const galat = ref('')
const berhasil = ref('')
const mengirim = ref(false)

const namaPengguna = computed(() => auth.user?.username || 'bidan')

const siap = computed(() =>
  lama.value !== '' &&
  baru.value.length >= SANDI_MINIMAL &&
  baru.value === ulang.value &&
  baru.value !== lama.value
)

// Alasan tombol mati ditampilkan apa adanya. Tombol yang mati tanpa keterangan
// membuat orang mengetik ulang berkali-kali sambil menebak apa yang kurang.
const petunjuk = computed(() => {
  if (!baru.value) return ''
  if (baru.value.length < SANDI_MINIMAL) return `Sandi baru minimal ${SANDI_MINIMAL} karakter.`
  if (lama.value && baru.value === lama.value) return 'Sandi baru harus berbeda dari sandi lama.'
  if (ulang.value && baru.value !== ulang.value) return 'Ulangan sandi belum sama.'
  return ''
})

async function kirim() {
  if (mengirim.value || !siap.value) return
  galat.value = ''
  berhasil.value = ''
  mengirim.value = true
  try {
    await api.gantiSandi({ old_password: lama.value, new_password: baru.value })
    lama.value = ''
    baru.value = ''
    ulang.value = ''
    berhasil.value = 'Sandi berhasil diganti. Perangkat lain yang masih masuk sudah dikeluarkan.'
  } catch (e) {
    galat.value = e.message || 'Gagal mengganti sandi'
  } finally {
    mengirim.value = false
  }
}

async function keluar() {
  if (!window.confirm('Keluar dari aplikasi?')) return
  await auth.keluar()
  router.replace({ name: 'masuk' })
}
</script>

<template>
  <div class="screen">
    <div class="head">
      <div class="eyebrow">AKUN</div>
      <div class="judul">{{ namaPengguna }}</div>
    </div>

    <form class="kartu" @submit.prevent="kirim">
      <div class="k-judul">Ganti kata sandi</div>
      <p class="k-sub">
        Setelah diganti, perangkat lain yang masih masuk akan dikeluarkan.
        Perangkat ini tetap masuk.
      </p>

      <label class="label" for="lama">Sandi sekarang</label>
      <input id="lama" v-model="lama" class="kolom" type="password" autocomplete="current-password" required />

      <label class="label" for="baru">Sandi baru</label>
      <input id="baru" v-model="baru" class="kolom" type="password" autocomplete="new-password" required />

      <label class="label" for="ulang">Ulangi sandi baru</label>
      <input id="ulang" v-model="ulang" class="kolom" type="password" autocomplete="new-password" required />

      <p v-if="petunjuk" class="petunjuk">{{ petunjuk }}</p>
      <p v-if="galat" class="galat" role="alert">{{ galat }}</p>
      <p v-if="berhasil" class="berhasil" role="status">{{ berhasil }}</p>

      <button class="tombol" type="submit" :disabled="mengirim || !siap">
        {{ mengirim ? 'Menyimpan…' : 'Ganti sandi' }}
      </button>
    </form>

    <div class="kartu">
      <div class="k-judul">Keluar</div>
      <p class="k-sub">Kamu perlu mengetik sandi lagi untuk masuk kembali.</p>
      <button class="tombol keluar" type="button" @click="keluar()">Keluar dari aplikasi</button>
    </div>
  </div>
</template>

<style scoped>
.screen { padding: 20px 18px 24px; display: flex; flex-direction: column; gap: 16px; }
.head { padding: 2px; }
.eyebrow { font-size: 10.5px; font-weight: 800; letter-spacing: 1.3px; color: var(--muted); }
.judul { font-size: 21px; font-weight: 800; letter-spacing: -.02em; color: var(--ink); }

.kartu {
  background: var(--card); border: 1px solid var(--border);
  border-radius: 16px; padding: 18px 16px;
  display: flex; flex-direction: column;
}
.k-judul { font-size: 15px; font-weight: 800; color: var(--ink); }
.k-sub { font-size: 12.5px; color: var(--text-secondary); margin: 4px 0 16px; line-height: 1.45; }

.label { font-size: 12px; font-weight: 700; color: var(--label); margin-bottom: 6px; }
.kolom {
  border: 1px solid var(--input-border); border-radius: 10px;
  padding: 11px 12px; font: inherit; font-size: 15px; color: var(--ink);
  background: var(--fill2); margin-bottom: 14px; width: 100%;
}
.kolom:focus { outline: 2px solid var(--accent); outline-offset: 1px; background: var(--card); }

.petunjuk { font-size: 12.5px; color: var(--muted); margin: 0 0 12px; }
.galat {
  background: var(--danger-bg); color: var(--danger);
  font-size: 13px; font-weight: 600; border-radius: 10px;
  padding: 9px 11px; margin: 0 0 12px;
}
.berhasil {
  background: #e8f6ee; color: var(--green);
  font-size: 13px; font-weight: 600; border-radius: 10px;
  padding: 9px 11px; margin: 0 0 12px;
}

.tombol {
  background: var(--accent); color: #fff; border: 0;
  border-radius: 10px; padding: 12px; font: inherit; font-size: 15px; font-weight: 700;
  cursor: pointer;
}
.tombol:hover:not(:disabled) { background: var(--accent-hover); }
.tombol:disabled { background: var(--disabled); cursor: default; }
.tombol.keluar { background: var(--danger-bg); color: var(--danger); }
.tombol.keluar:hover { background: var(--danger); color: #fff; }

/* Layar lebar: kartu tidak ikut melar selebar monitor — baris sepanjang itu
   melelahkan dibaca, dan kolom sandi selebar 1200px terlihat salah. */
@media (min-width: 920px) {
  .screen { padding: 8px 4px 24px; max-width: 560px; }
}
</style>
