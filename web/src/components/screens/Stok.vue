<script setup>
import { computed } from 'vue'
import { useKlinik } from '../../stores/klinik'
const k = useKlinik()
const d = computed(() => k.derived)
</script>

<template>
  <div class="screen">
    <div class="top">
      <h1 class="title">Stok Obat</h1>
      <button class="add" @click="k.openPickStock()"><span class="plus">+</span> Tambah</button>
    </div>
    <div class="subtitle">{{ d.medCount }} jenis obat terdaftar</div>

    <div class="list">
      <div v-for="m in d.medsList" :key="m.id" class="line">
        <span class="dot" :style="{ background: m.dot }" />
        <div class="grow">
          <div class="ln-name">{{ m.name }}</div>
          <div class="ln-cat">{{ m.cat }}</div>
        </div>
        <div class="qty"><span class="ln-qty" :style="{ color: m.qtyColor }">{{ m.qty }}</span><span class="ln-unit">{{ m.unit }}</span></div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.screen { padding: 20px 18px 24px; }
.top { display: flex; align-items: center; justify-content: space-between; margin-bottom: 6px; }
.title { margin: 0; font-size: 22px; font-weight: 800; letter-spacing: -.4px; color: var(--ink); }
.add { height: 40px; padding: 0 15px; border-radius: 20px; background: var(--accent); color: #fff; font-size: 13.5px; font-weight: 700; display: flex; align-items: center; gap: 6px; box-shadow: 0 6px 16px rgba(18,128,106,.3); }
.add:hover { background: var(--accent-hover); }
.plus { font-size: 19px; }
.subtitle { font-size: 13px; color: var(--muted); margin-bottom: 18px; }
.list { display: flex; flex-direction: column; gap: 9px; }
.line { display: flex; align-items: center; gap: 12px; background: #fff; border: 1px solid var(--border); border-radius: 14px; padding: 12px 14px; }
.grow { flex: 1; min-width: 0; }
.qty { flex-shrink: 0; white-space: nowrap; }
.ln-name { font-size: 14.5px; font-weight: 600; color: var(--ink); }
.ln-cat { font-size: 12px; color: var(--muted); }
.ln-qty { font-size: 16px; font-weight: 800; }
/* Jaraknya diatur di sini, bukan lewat spasi di dalam template: kompilator
   Vue merapatkan spasi tepi, dan hasilnya terbaca "8butir". */
.ln-unit { font-size: 11px; font-weight: 600; color: var(--muted); margin-left: 4px; }

@media (min-width: 920px) {
  .screen { padding: 8px 4px 24px; }
  .title { font-size: 26px; }
  .list { display: grid; grid-template-columns: repeat(auto-fill, minmax(320px, 1fr)); gap: 10px; }
}
</style>
