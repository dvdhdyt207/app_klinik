<script setup>
// Toggle 52x30 (knob 24) atau 'small' 46x27 (knob 21). Presentasional.
defineProps({
  on: Boolean,
  variant: { type: String, default: 'big' },
  trackOn: { type: String, default: 'var(--accent)' },
  trackOff: { type: String, default: 'var(--disabled)' },
})
</script>

<template>
  <span class="toggle" :class="variant" :style="{ background: on ? trackOn : trackOff }">
    <span class="knob" :style="on ? { right: '3px' } : { left: '3px' }" />
  </span>
</template>

<style scoped>
.toggle { position: relative; display: inline-block; border-radius: 16px; flex-shrink: 0; }
.toggle.big { width: 52px; height: 30px; }
.toggle.small { width: 46px; height: 27px; }
.knob { position: absolute; top: 3px; background: #fff; border-radius: 50%; transition: all .15s ease; }
.toggle.big .knob { width: 24px; height: 24px; }
.toggle.small .knob { width: 21px; height: 21px; }
</style>
