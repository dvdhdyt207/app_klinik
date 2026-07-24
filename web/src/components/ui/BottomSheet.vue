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
.sheet-backdrop {
  position: absolute; inset: 0; z-index: 30;
  background: rgba(16, 22, 30, .4);
  display: flex; flex-direction: column; justify-content: flex-end;
}
.sheet {
  background: #fff; border-top-left-radius: 22px; border-top-right-radius: 22px;
  padding: 20px 20px 26px; animation: sheet-up .25s ease;
  max-height: 90%; overflow-y: auto;
}
.handle { width: 40px; height: 4px; border-radius: 2px; background: var(--input-border); margin: 0 auto 16px; }

/* desktop: sheet jadi dialog terpusat */
@media (min-width: 920px) {
  .sheet-backdrop { justify-content: center; align-items: center; padding: 24px; }
  .sheet {
    width: min(460px, 92vw); max-height: 86vh;
    border-radius: 22px; padding: 24px 24px 26px;
    animation: dialog-pop .2s ease;
  }
  .handle { display: none; }
}
@keyframes dialog-pop { from { opacity: .5; transform: scale(.97); } to { opacity: 1; transform: scale(1); } }
</style>
