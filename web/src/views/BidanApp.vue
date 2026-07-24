<script setup>
import { useKlinik } from '../stores/klinik'
import Beranda from '../components/screens/Beranda.vue'
import Kunjungan from '../components/screens/Kunjungan.vue'
import Stok from '../components/screens/Stok.vue'
import Rekap from '../components/screens/Rekap.vue'
import TabBar from '../components/TabBar.vue'
import VisitModal from '../components/modals/VisitModal.vue'
import PickModal from '../components/modals/PickModal.vue'
import JadwalModal from '../components/modals/JadwalModal.vue'
import EventFormModal from '../components/modals/EventFormModal.vue'
import AddStockSheet from '../components/modals/AddStockSheet.vue'
import AwaySheet from '../components/modals/AwaySheet.vue'

const k = useKlinik()
</script>

<template>
  <div class="app-frame">
    <div v-if="k.loading" class="loading">
      <div class="spinner" />
      <div class="loading-text">Memuat…</div>
    </div>

    <template v-else>
      <button v-if="k.error" class="errbar" @click="k.refresh()">
        Gagal terhubung ke server ({{ k.error }}). Ketuk untuk coba lagi.
      </button>

      <div class="scroll-area">
        <div class="content-wrap">
          <Beranda v-if="k.screen === 'beranda'" />
          <Kunjungan v-else-if="k.screen === 'kunjungan'" />
          <Stok v-else-if="k.screen === 'stok'" />
          <Rekap v-else-if="k.screen === 'rekap'" />
        </div>
      </div>

      <TabBar />

      <VisitModal />
      <PickModal />
      <JadwalModal />
      <EventFormModal />
      <AddStockSheet />
      <AwaySheet />
    </template>
  </div>
</template>

<style scoped>
.loading { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: center; gap: 12px; }
.spinner { width: 38px; height: 38px; border-radius: 50%; border: 3px solid #dbe6f7; border-top-color: var(--accent); animation: spin 1s linear infinite; }
.loading-text { color: var(--muted); font-size: 13px; }
@keyframes spin { to { transform: rotate(360deg); } }
.errbar { width: 100%; background: var(--danger-bg); color: var(--danger); padding: 10px; font-size: 12.5px; font-weight: 600; text-align: center; flex-shrink: 0; }
</style>
