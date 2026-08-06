<script setup>
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '../../stores/auth'
import { useKlinik } from '../../stores/klinik'
import { api } from '../../api/client'

const auth = useAuth()
const router = useRouter()
const k = useKlinik()

// ---------- profil klinik ----------
const profil = ref({ clinic: '', alamat: '', whatsapp: '' })
const profilSiap = ref(false)
const profilGalat = ref('')
const profilBerhasil = ref('')
const menyimpanProfil = ref(false)

// Isi form sekali saja, saat data server pertama tiba. Menyalin ulang setiap
// kali store menyegarkan dirinya (tiap 15 detik) akan menimpa ketikan bidan di
// tengah jalan.
watch(() => k.profil, (p) => {
  if (profilSiap.value || !p) return
  if (!p.clinic && !p.alamat && !p.whatsapp && k.loading) return
  profil.value = { clinic: p.clinic || '', alamat: p.alamat || '', whatsapp: p.whatsapp || '' }
  profilSiap.value = true
}, { immediate: true, deep: true })

async function simpanProfil() {
  if (menyimpanProfil.value) return
  profilGalat.value = ''
  profilBerhasil.value = ''
  menyimpanProfil.value = true
  try {
    const baru = await k.simpanProfil({ ...profil.value })
    // Nomor ditampilkan dalam bentuk yang benar-benar tersimpan, bukan yang
    // tadi diketik — supaya bidan melihat sendiri 0812… berubah jadi 62812…
    // alih-alih mengira nomornya tidak jadi disimpan.
    profil.value = { clinic: baru.clinic, alamat: baru.alamat, whatsapp: baru.whatsapp }
    profilBerhasil.value = 'Tersimpan. Halaman pasien ikut berubah.'
  } catch (e) {
    profilGalat.value = e.message || 'Gagal menyimpan profil'
  } finally {
    menyimpanProfil.value = false
  }
}

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
    <header class="head">
      <h1 class="judul">Setelan</h1>
      <p class="sub">{{ profil.clinic || 'Klinik' }}</p>
    </header>

    <form class="kartu" @submit.prevent="simpanProfil">
      <div class="k-judul">Profil klinik</div>
      <p class="k-sub">Tampil di halaman yang dibuka pasien.</p>

      <label class="label" for="nama-klinik">Nama klinik</label>
      <input id="nama-klinik" v-model="profil.clinic" class="kolom" maxlength="120" required />

      <label class="label" for="alamat">Alamat</label>
      <textarea id="alamat" v-model="profil.alamat" class="kolom area" maxlength="255"
        rows="3" placeholder="Kosongkan bila belum ingin ditampilkan" />

      <label class="label" for="wa">Nomor WhatsApp</label>
      <input id="wa" v-model="profil.whatsapp" class="kolom" inputmode="tel"
        placeholder="081234567890" />
      <p class="petunjuk">
        Boleh ditulis 0812…, +62 812…, atau pakai tanda hubung — akan
        diseragamkan sendiri. Dikosongkan berarti tombol WhatsApp tidak
        ditampilkan ke pasien.
      </p>

      <p v-if="profilGalat" class="galat" role="alert">{{ profilGalat }}</p>
      <p v-if="profilBerhasil" class="berhasil" role="status">{{ profilBerhasil }}</p>

      <button class="tombol" type="submit" :disabled="menyimpanProfil || !profil.clinic.trim()">
        {{ menyimpanProfil ? 'Menyimpan…' : 'Simpan profil' }}
      </button>
    </form>

    <form class="kartu" @submit.prevent="kirim">
      <div class="k-judul">Ganti kata sandi <span class="k-akun">akun {{ namaPengguna }}</span></div>
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
.screen { padding: 16px 16px 24px; display: flex; flex-direction: column; gap: 16px; }
.head { padding: 0 2px; }
.judul { margin: 0; font-size: 20px; font-weight: 700; letter-spacing: -.02em; color: var(--ink); }
.sub { margin: 3px 0 0; font-size: 13px; color: var(--muted); }

.kartu {
  background: var(--card); border: 1px solid var(--line);
  border-radius: var(--ra-lg); padding: 18px;
  display: flex; flex-direction: column;
}
.k-judul { font-size: 14.5px; font-weight: 700; color: var(--ink); }
.k-akun { font-size: 12px; font-weight: 500; color: var(--muted); }
.area { resize: vertical; min-height: 76px; line-height: 1.5; }
.k-sub { font-size: 12.5px; color: var(--muted); margin: 5px 0 16px; line-height: 1.5; }

.label { font-size: 12px; font-weight: 600; color: var(--label); margin-bottom: 6px; }
.kolom {
  border: 1px solid var(--input-border); border-radius: var(--ra-md);
  padding: 11px 12px; font: inherit; font-size: 14.5px; color: var(--ink);
  background: var(--fill2); margin-bottom: 14px; width: 100%;
  transition: border-color .15s ease, background .15s ease;
}
.kolom:focus { outline: none; border-color: var(--accent); background: var(--card); box-shadow: 0 0 0 3px var(--accent-soft); }

.petunjuk { font-size: 12.5px; color: var(--muted); margin: 0 0 12px; line-height: 1.5; }
.galat {
  background: var(--danger-bg); color: var(--danger);
  font-size: 13px; font-weight: 600; border-radius: var(--ra-md);
  padding: 10px 12px; margin: 0 0 12px;
}
/* Sama dengan Login.vue: pesan server diawali huruf kecil, dibesarkan lewat
   tampilan supaya teks aslinya tetap utuh. */
.galat::first-letter { text-transform: uppercase; }
.berhasil {
  background: var(--accent-soft); color: var(--accent-ink);
  font-size: 13px; font-weight: 600; border-radius: var(--ra-md);
  padding: 10px 12px; margin: 0 0 12px;
}

.tombol {
  background: var(--accent); color: #fff; border: 0;
  border-radius: var(--ra-md); min-height: 40px; padding: 0 14px;
  font: inherit; font-size: 14px; font-weight: 600;
  cursor: pointer; transition: background .15s ease;
}
.tombol:hover:not(:disabled) { background: var(--accent-hover); }
.tombol:disabled { background: var(--disabled); cursor: default; }
.tombol.keluar { background: var(--danger-bg); color: var(--danger); }
.tombol.keluar:hover { background: var(--danger); color: #fff; }

/* Layar lebar: kartu berdampingan, bukan satu kolom sempit yang menyisakan
   separuh layar kosong. Tiap kartu tetap dibatasi lebarnya oleh jumlah kolom —
   kolom sandi selebar 1200px terlihat salah, dan baris sepanjang itu
   melelahkan dibaca. */
@media (min-width: 920px) {
  .screen {
    padding: 8px 4px 24px;
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(380px, 1fr));
    max-width: 1000px;
    align-items: start;
  }
  .head { grid-column: 1 / -1; }
  .judul { font-size: 23px; }
}
</style>
