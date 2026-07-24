<script setup>
import { computed } from 'vue'
import { useKlinik } from '../../stores/klinik'
const k = useKlinik()
const d = computed(() => k.derived)
</script>

<template>
  <div class="screen">
    <div class="top">
      <h1 class="title">Kunjungan</h1>
      <button class="fab" @click="k.openVisit()">+</button>
    </div>
    <div class="subtitle">{{ d.totalVisits }} total · {{ d.todayCount }} hari ini</div>

    <div class="list">
      <div v-for="v in d.visitCards" :key="v.id" class="vcard">
        <div class="vhead">
          <div class="vavatar">{{ v.initial }}</div>
          <div class="grow">
            <div class="vname">{{ v.name }}</div>
            <div class="vmeta">{{ v.meta }}</div>
          </div>
          <div class="vdate">{{ v.dateLabel }}</div>
        </div>
        <div v-if="v.gejala" class="vgejala"><span class="glabel">Gejala: </span>{{ v.gejala }}</div>
        <div v-if="v.chips.length" class="vchips">
          <span v-for="(c, i) in v.chips" :key="i" class="vchip">{{ c }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.screen { padding: 20px 18px 24px; }
.top { display: flex; align-items: center; justify-content: space-between; margin-bottom: 6px; }
.title { margin: 0; font-size: 22px; font-weight: 800; letter-spacing: -.4px; color: var(--ink); }
.fab { width: 40px; height: 40px; border-radius: 20px; background: var(--accent); color: #fff; font-size: 24px; display: flex; align-items: center; justify-content: center; box-shadow: 0 6px 16px rgba(47,108,224,.3); }
.subtitle { font-size: 13px; color: var(--muted); margin-bottom: 18px; }
.list { display: flex; flex-direction: column; gap: 10px; }
.vcard { background: #fff; border: 1px solid var(--border); border-radius: 16px; padding: 15px; }
.vhead { display: flex; align-items: center; gap: 12px; }
.vavatar { width: 40px; height: 40px; border-radius: 12px; background: #dbe6f7; color: var(--accent); font-weight: 800; font-size: 16px; display: flex; align-items: center; justify-content: center; }
.grow { flex: 1; min-width: 0; }
.vname { font-size: 15px; font-weight: 700; color: var(--ink); }
.vmeta { font-size: 12.5px; color: var(--muted); }
.vdate { font-size: 11.5px; color: var(--muted); font-weight: 600; }
.vgejala { margin-top: 9px; font-size: 12.5px; color: var(--text-secondary); line-height: 1.4; }
.glabel { color: var(--muted2); font-weight: 700; }
.vchips { margin-top: 11px; padding-top: 11px; border-top: 1px solid #eef2f7; display: flex; flex-wrap: wrap; gap: 6px; }
.vchip { background: var(--fill); border-radius: 8px; padding: 5px 9px; font-size: 12px; font-weight: 600; color: #3b4a5e; }

@media (min-width: 920px) {
  .screen { padding: 8px 4px 24px; }
  .title { font-size: 26px; }
  .list { display: grid; grid-template-columns: repeat(auto-fill, minmax(330px, 1fr)); gap: 12px; align-items: start; }
}
</style>
