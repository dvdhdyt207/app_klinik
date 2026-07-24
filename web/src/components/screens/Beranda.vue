<script setup>
import { computed } from 'vue'
import { useKlinik } from '../../stores/klinik'
import { hhmm, untilLabel } from '../../lib/format'
import Toggle from '../ui/Toggle.vue'

const k = useKlinik()
const st = computed(() => k.status)
const away = computed(() => !st.value.bidanHadir)
const d = computed(() => k.derived)
// bergantung pada k.now agar countdown ikut ter-refresh berkala
const untilText = computed(() => { void k.now; return st.value.awayUntil ? untilLabel(st.value.awayUntil) : '' })
</script>

<template>
  <div class="screen">
    <!-- header -->
    <div class="head">
      <div>
        <div class="eyebrow">KLINIK</div>
        <div class="clinic">Bidan Pit</div>
      </div>
      <div class="avatar">P</div>
    </div>

    <!-- status hero: hadir -->
    <button v-if="!away" class="hero-present" @click="k.openAway()">
      <div class="row-between">
        <span class="hp-label">Status Bidan</span>
        <Toggle :on="true" trackOn="rgba(255,255,255,.32)" />
      </div>
      <div class="hp-title">Sedang di klinik</div>
      <div class="hp-sub">Ketuk untuk atur status keluar</div>
    </button>

    <!-- status hero: tidak di klinik -->
    <div v-else class="hero-away">
      <div class="row-between">
        <span class="ha-label">Status Bidan</span>
        <Toggle :on="false" @click="k.setHadir()" style="cursor:pointer" />
      </div>
      <div class="ha-title">Tidak di klinik</div>
      <div v-if="st.awayNote" class="ha-note">{{ st.awayNote }}</div>

      <div v-if="st.awayUntil" class="ha-pill">
        <span class="pill-strong">Perkiraan kembali {{ hhmm(st.awayUntil) }}</span>
        <span class="pill-muted">{{ untilText }}</span>
      </div>
      <div v-else class="ha-pill">
        <span class="pill-muted">Waktu kembali belum dipastikan</span>
      </div>

      <div class="ha-actions">
        <button class="chip" @click="k.extend(15)">+15 mnt</button>
        <button class="chip" @click="k.extend(30)">+30 mnt</button>
        <button class="chip" @click="k.openAway()">Ubah</button>
      </div>
      <button class="back-btn" @click="k.setHadir()">Saya sudah kembali</button>
    </div>

    <!-- web status link -->
    <router-link to="/" class="weblink">
      <span class="wl-left"><span class="dot" style="background:var(--healthy)" /> Status ini tampil di halaman web pasien</span>
      <span class="chev">›</span>
    </router-link>

    <!-- quick stats -->
    <div class="stats">
      <div class="card stat">
        <div class="stat-num">{{ d.todayCount }}</div>
        <div class="stat-lbl">Kunjungan hari ini</div>
      </div>
      <div class="card stat">
        <div class="stat-num danger">{{ d.lowCount }}</div>
        <div class="stat-lbl">Obat menipis</div>
      </div>
    </div>

    <!-- stok menipis -->
    <div class="section-head">
      <span class="sh-title">Stok menipis</span>
      <button class="sh-action" @click="k.goScreen('stok')">Lihat stok</button>
    </div>
    <div class="list">
      <div v-for="m in d.low" :key="m.id" class="line">
        <span class="dot" :style="{ background: m.color }" />
        <div class="grow">
          <div class="ln-name">{{ m.name }}</div>
          <div class="ln-cat">{{ m.cat }}</div>
        </div>
        <div class="ln-qty" :style="{ color: m.color }">{{ m.qty }} <span class="ln-unit">{{ m.unit }}</span></div>
      </div>
    </div>

    <!-- jadwal & agenda -->
    <div class="section-head">
      <span class="sh-title">Jadwal & agenda</span>
      <button class="sh-action" @click="k.openJadwal()">Kelola</button>
    </div>
    <div v-if="d.upcoming.length" class="list">
      <div v-for="e in d.upcoming.slice(0, 3)" :key="e.id" class="line">
        <span class="dot" :style="{ background: e.dot }" />
        <div class="grow">
          <div class="ln-name">{{ e.title }}</div>
          <div class="ln-cat">{{ e.when }}</div>
        </div>
      </div>
    </div>
    <div v-else class="dashed">Belum ada agenda mendatang</div>

    <button class="visit-btn" @click="k.openVisit()">
      <span style="font-size:20px">+</span> Catat Kunjungan
    </button>
  </div>
</template>

<style scoped>
.screen { padding: 20px 18px 24px; }
.head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 20px; }
.eyebrow { font-size: 11px; font-weight: 600; letter-spacing: 1.3px; color: var(--muted); }
.clinic { font-size: 22px; font-weight: 800; letter-spacing: -.4px; margin-top: 2px; color: var(--ink); }
.avatar { width: 42px; height: 42px; border-radius: 21px; background: #dbe6f7; color: var(--accent); font-weight: 700; display: flex; align-items: center; justify-content: center; }

.row-between { display: flex; align-items: center; justify-content: space-between; }
.hero-present { display: block; width: 100%; text-align: left; background: var(--accent); border-radius: 20px; padding: 20px; box-shadow: 0 12px 26px rgba(47,108,224,.28); }
.hp-label { font-size: 12.5px; font-weight: 600; color: #fff; opacity: .85; }
.hp-title { font-size: 20px; font-weight: 800; color: #fff; margin-top: 14px; }
.hp-sub { font-size: 12.5px; color: #fff; opacity: .8; margin-top: 3px; }

.hero-away { background: #fff; border: 1px solid var(--border2); border-radius: 20px; padding: 20px; }
.ha-label { font-size: 12.5px; font-weight: 600; color: var(--label); }
.ha-title { font-size: 20px; font-weight: 800; color: var(--ink); margin-top: 14px; }
.ha-note { font-size: 13px; color: var(--text-secondary); margin-top: 5px; font-weight: 600; }
.ha-pill { margin-top: 12px; background: var(--fill2); border-radius: 12px; padding: 11px 13px; display: flex; flex-wrap: wrap; align-items: baseline; gap: 8px; }
.pill-strong { font-size: 13px; color: var(--ink); font-weight: 700; }
.pill-muted { font-size: 12.5px; color: var(--muted); font-weight: 600; }
.ha-actions { display: flex; gap: 8px; margin-top: 12px; }
.chip { flex: 1; text-align: center; padding: 9px 0; border-radius: 10px; background: var(--fill); color: #3b4a5e; font-size: 12.5px; font-weight: 700; }
.back-btn { width: 100%; margin-top: 10px; padding: 13px; border-radius: 12px; background: var(--green); color: #fff; font-size: 14px; font-weight: 700; }

.weblink { display: flex; align-items: center; justify-content: space-between; margin-top: 12px; background: #eef4ff; border: 1px solid #d9e5fb; border-radius: 14px; padding: 11px 14px; }
.wl-left { display: flex; align-items: center; gap: 9px; font-size: 12.5px; font-weight: 600; color: var(--accent); }
.chev { font-size: 18px; color: var(--accent); }

.stats { display: flex; gap: 10px; margin-top: 14px; }
.stat { flex: 1; }
.stat-num { font-size: 22px; font-weight: 800; letter-spacing: -.4px; color: var(--ink); }
.stat-num.danger { color: var(--danger); }
.stat-lbl { font-size: 11.5px; color: var(--muted); font-weight: 600; }

.section-head { display: flex; align-items: center; justify-content: space-between; margin: 24px 2px 12px; }
.sh-title { font-size: 15px; font-weight: 700; color: var(--ink); }
.sh-action { font-size: 12.5px; font-weight: 600; color: var(--accent); }

.list { display: flex; flex-direction: column; gap: 9px; }
.line { display: flex; align-items: center; gap: 12px; background: #fff; border: 1px solid var(--border); border-radius: 14px; padding: 12px 14px; }
.grow { flex: 1; min-width: 0; }
.ln-name { font-size: 14px; font-weight: 600; color: var(--ink); }
.ln-cat { font-size: 12px; color: var(--muted); }
.ln-qty { font-size: 15px; font-weight: 800; }
.ln-unit { font-size: 11px; font-weight: 600; color: var(--muted); }

.visit-btn { display: flex; align-items: center; justify-content: center; gap: 9px; width: 100%; margin-top: 22px; background: var(--ink); color: #fff; border-radius: 16px; padding: 15px; font-size: 15px; font-weight: 700; }
</style>
