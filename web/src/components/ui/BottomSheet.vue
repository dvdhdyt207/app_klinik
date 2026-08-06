<script setup>
// Bottom sheet: slide up dari bawah + backdrop gelap. Menutup frame HP.
defineProps({ visible: Boolean })
defineEmits(['close'])
</script>

<template>
  <div v-if="visible" class="sheet-backdrop" @click="$emit('close')">
    <div class="sheet" @click.stop>
      <div class="handle" />
      <slot />
    </div>
  </div>
</template>

<style scoped>
/* Lapisan gelap sheet ini satu-satunya untuk modal jenis sheet — lihat catatan
   di BidanApp.vue soal kenapa lapisan milik layar tidak boleh ikut dipasang. */
.sheet-backdrop {
  position: absolute; inset: 0; z-index: 30;
  background: rgba(20, 33, 28, .45);
  display: flex; flex-direction: column; justify-content: flex-end;
}
.sheet {
  background: #fff; border-top-left-radius: 16px; border-top-right-radius: 16px;
  padding: 18px 18px 24px; animation: sheet-up .25s ease;
  max-height: 90%; overflow-y: auto;
}
.handle { width: 40px; height: 4px; border-radius: 2px; background: var(--input-border); margin: 0 auto 16px; }

/* desktop: sheet jadi dialog terpusat */
@media (min-width: 920px) {
  .sheet-backdrop { justify-content: center; align-items: center; padding: 24px; }
  .sheet {
    width: min(440px, 92vw); max-height: 86vh;
    border-radius: var(--ra-lg); padding: 22px;
    animation: dialog-pop .2s ease;
  }
  .handle { display: none; }
}
@keyframes dialog-pop { from { opacity: .5; transform: scale(.97); } to { opacity: 1; transform: scale(1); } }
</style>
