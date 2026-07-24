<script setup>
import { computed } from 'vue'
import { useKlinik } from '../../stores/klinik'
import { CAT } from '../../lib/catalog'
import BottomSheet from '../ui/BottomSheet.vue'

const k = useKlinik()
const st = computed(() => k.stockTarget)
const cfg = computed(() => (st.value ? (CAT[st.value.cat] || CAT.Tablet) : CAT.Tablet))
const mult = computed(() => (cfg.value.units[k.stockUnitIdx] || cfg.value.units[0]).mult)
const n = computed(() => parseInt(k.stockCount, 10) || 0)
const preview = computed(() =>
  '= +' + (n.value * mult.value) + ' ' + cfg.value.base + ' (total jadi ' + ((st.value?.qty || 0) + n.value * mult.value) + ')')

function dec() { k.stockCount = String(Math.max(1, (parseInt(k.stockCount, 10) || 0) - 1)) }
function inc() { k.stockCount = String((parseInt(k.stockCount, 10) || 0) + 1) }
function countInput(e) { k.stockCount = e.target.value.replace(/[^0-9]/g, '') }
</script>

<template>
  <BottomSheet :visible="k.modal === 'addstock'" @close="k.closeAddStock()">
    <template v-if="st">
      <div class="eyebrow">TAMBAH STOK</div>
      <div class="name">{{ st.name }}</div>
      <div class="now">Stok sekarang: {{ (st.qty || 0) + ' ' + cfg.base }}</div>

      <label class="label">Satuan</label>
      <div class="chips">
        <button v-for="(u, i) in cfg.units" :key="u.label" class="chip" :class="{ sel: i === k.stockUnitIdx }" @click="k.stockUnitIdx = i">
          {{ u.label }}
        </button>
      </div>

      <label class="label">Jumlah</label>
      <div class="countrow">
        <button class="cbtn" @click="dec">−</button>
        <input class="field count" :value="k.stockCount" @input="countInput" inputmode="numeric" />
        <button class="cbtn" @click="inc">+</button>
      </div>
      <div class="preview">{{ preview }}</div>

      <div class="actions">
        <button class="btn-ghost" @click="k.closeAddStock()">Batal</button>
        <button class="btn grow2" @click="k.confirmAddStock()">Simpan Stok</button>
      </div>
    </template>
  </BottomSheet>
</template>

<style scoped>
.eyebrow { font-size: 12px; font-weight: 700; color: var(--muted); letter-spacing: 1px; }
.name { font-size: 20px; font-weight: 800; color: var(--ink); margin-top: 3px; }
.now { font-size: 12.5px; color: var(--muted); margin-bottom: 18px; }
.chips { display: flex; gap: 8px; margin: 8px 0 16px; flex-wrap: wrap; }
.chip { padding: 9px 14px; border-radius: 12px; border: 1.5px solid var(--border2); background: #fff; color: var(--label); font-size: 13px; font-weight: 700; }
.chip.sel { border-color: var(--accent); background: var(--tablet-tint); color: var(--accent); }
.countrow { display: flex; align-items: center; gap: 12px; margin: 8px 0; }
.cbtn { width: 46px; height: 46px; border-radius: 12px; background: var(--fill); color: #3b4a5e; font-size: 24px; display: flex; align-items: center; justify-content: center; }
.count { flex: 1; text-align: center; font-size: 20px; font-weight: 800; background: var(--fill2); }
.preview { font-size: 13px; color: var(--accent); font-weight: 700; margin-bottom: 20px; }
.actions { display: flex; gap: 10px; }
.btn-ghost { flex: 1; padding: 15px; border-radius: 14px; background: var(--fill); color: #3b4a5e; font-size: 15px; font-weight: 700; }
.grow2 { flex: 2; }
</style>
