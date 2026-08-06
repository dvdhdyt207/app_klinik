<script setup>
import { computed } from 'vue'
import { useKlinik } from '../../stores/klinik'
import { KATEGORI, baseUnit } from '../../lib/catalog'
import ModalHeader from '../ui/ModalHeader.vue'

const k = useKlinik()
const m = computed(() => k.medDraft)

// Nilai asli dibaca ulang dari data server, bukan disalin ke dalam draft:
// draft-nya ikut berubah saat diketik, jadi ia tidak bisa jadi pembanding.
const asli = computed(() => (k.data.medicines || []).find((x) => x.id === m.value?.id) || null)

// Mengubah kategori mengubah satuan dasarnya — butir jadi botol, dan angka
// stok yang sudah ada ikut berpindah arti tanpa berubah nilainya. Peringatan
// hanya muncul pada obat yang memang sudah punya stok tercatat.
const gantiKategori = computed(() =>
  !!asli.value && !!m.value && m.value.cat !== asli.value.cat)
const satuanBaru = computed(() => (m.value ? baseUnit(m.value.cat) : ''))
const satuanLama = computed(() => (asli.value ? baseUnit(asli.value.cat) : ''))

const canSave = computed(() => !!m.value && !!(m.value.name || '').trim() && !k.medBusy)
</script>

<template>
  <div v-if="k.modal === 'med'" class="fullmodal">
    <ModalHeader title="Ubah Obat" @back="k.closeMedEdit()" />
    <div v-if="m" class="body">
      <label class="label">Nama obat</label>
      <input class="field mt mb" v-model="m.name" placeholder="mis. Paracetamol 500mg" />

      <label class="label">Kategori</label>
      <div class="segs mt mb">
        <button v-for="c in KATEGORI" :key="c" class="seg" :class="{ on: m.cat === c }" @click="m.cat = c">
          {{ c }}
        </button>
      </div>

      <p v-if="gantiKategori" class="peringatan">
        Satuan ikut berubah dari <b>{{ satuanLama }}</b> ke <b>{{ satuanBaru }}</b>.
        Jumlah di bawah tidak dihitung ulang — periksa lagi sebelum menyimpan.
      </p>

      <label class="label">Jumlah stok</label>
      <div class="qrow mt mb">
        <input class="field grow" type="number" inputmode="numeric" min="0" v-model="m.qty" />
        <span class="unit">{{ satuanBaru }}</span>
      </div>
      <p class="bantuan">
        Isi dengan hitungan fisik yang sebenarnya. Angka ini menimpa stok tercatat,
        bukan menambahnya — untuk menambah, pakai <b>Tambah stok</b>.
      </p>

      <p v-if="k.medErr" class="galat">{{ k.medErr }}</p>

      <template v-if="!k.medHapusConfirm">
        <button class="del" @click="k.medHapusConfirm = true">Hapus obat ini</button>
      </template>
      <div v-else class="konfirmasi">
        <p class="k-teks">
          Hapus <b>{{ asli ? asli.name : m.name }}</b> dari daftar stok?
          Riwayat kunjungan yang pernah memakainya tetap utuh.
        </p>
        <div class="k-aksi">
          <button class="k-batal" @click="k.medHapusConfirm = false">Batal</button>
          <button class="k-ya" :disabled="k.medBusy" @click="k.hapusMed()">Ya, hapus</button>
        </div>
      </div>
    </div>
    <div class="footer">
      <button class="btn" :class="{ 'is-disabled': !canSave }" :disabled="!canSave" @click="k.saveMed()">
        Simpan Perubahan
      </button>
    </div>
  </div>
</template>

<style scoped>
.body { flex: 1; overflow-y: auto; padding: 18px; }
.mt { margin-top: 7px; } .mb { margin-bottom: 16px; }

/* Tiga kategori saja, dan pilihannya saling meniadakan — deretan tombol lebih
   cepat dibaca sekali lihat daripada <select> yang isinya harus dibuka dulu. */
.segs { display: flex; gap: 8px; }
.seg {
  flex: 1; min-height: 40px; border-radius: var(--ra-md);
  border: 1px solid var(--input-border); background: #fff;
  font-size: 13.5px; font-weight: 600; color: var(--muted);
  transition: background .15s ease, color .15s ease, border-color .15s ease;
}
.seg:hover { border-color: var(--accent); }
.seg.on { background: var(--accent); border-color: var(--accent); color: #fff; }

.qrow { display: flex; align-items: center; gap: 10px; }
.grow { flex: 1; min-width: 0; }
.unit { font-size: 13px; color: var(--muted); flex-shrink: 0; }

.peringatan {
  margin: 0 0 16px; padding: 11px 13px; border-radius: var(--ra-md);
  background: var(--warning-bg); color: var(--ink);
  font-size: 12.5px; line-height: 1.5;
}
.bantuan { margin: -8px 0 18px; font-size: 12.5px; color: var(--muted); line-height: 1.5; }
.galat {
  margin: 0 0 14px; padding: 11px 13px; border-radius: var(--ra-md);
  background: var(--del-bg); color: var(--del); font-size: 13px; line-height: 1.45;
}

.del {
  width: 100%; padding: 13px; border-radius: 12px;
  background: var(--del-bg); color: var(--del); font-size: 14px; font-weight: 700; margin-top: 6px;
}
/* Penghapusan minta konfirmasi di tempat, bukan lewat dialog sistem: dialog
   bawaan browser muncul di luar bahasa visual aplikasi dan pada layar sentuh
   tombolnya berdempetan. */
.konfirmasi {
  margin-top: 6px; padding: 14px; border-radius: 12px;
  border: 1px solid var(--del); background: var(--del-bg);
}
.k-teks { margin: 0 0 12px; font-size: 13px; color: var(--ink); line-height: 1.5; }
.k-aksi { display: flex; gap: 8px; }
.k-batal, .k-ya { flex: 1; min-height: 40px; border-radius: var(--ra-md); font-size: 13.5px; font-weight: 700; }
.k-batal { background: #fff; border: 1px solid var(--input-border); color: var(--ink); }
.k-ya { background: var(--del); color: #fff; }
.k-ya:disabled { opacity: .6; }

.footer { padding: 16px; background: #fff; border-top: 1px solid var(--border); }
</style>
