# Handoff: Sistem Manajemen Klinik Bidan Pit

## Overview
A mobile web application for a solo midwife's clinic ("Klinik Bidan Pit"). It has **two apps**:

1. **Bidan App** (`Klinik Bidan Pit App.dc.html`) — private, used only by the midwife on an Android phone. Records patient visits, manages medicine stock, shows recaps, and controls her live availability status + schedule.
2. **Patient Web Page** (`Status Bidan Web.dc.html`) — a public, read-only landing page. Patients open it to see whether the midwife is currently at the clinic, her estimated return time if she's out, and her upcoming schedule (e.g. "out of town for 3 days").

The two communicate: when the midwife flips her status or edits her schedule in the Bidan App, the Patient Web Page updates in real time. In this prototype that "API" is **simulated** with `localStorage` + `BroadcastChannel` on one origin. **In production this must be replaced by a real backend** (the app POSTs status/schedule; the web page GETs it, ideally with websockets/polling for live updates).

- **Language:** Bahasa Indonesia (all UI copy).
- **Primary device:** Android phone only (single-column, ~412px design width). Patient page is mobile-first but fine on desktop (centered, max-width 460px).

## About the Design Files
The files in this bundle are **design references created in HTML** — interactive prototypes showing the intended look and behavior. They are **not production code to copy directly**. The task is to **recreate these designs in the target codebase's environment** (React, Vue, Flutter, native Android, etc.) using its established patterns, then wire the real backend.

Notes on the prototype tech (ignore when porting):
- The prototype is authored as "Design Components" — logic lives in a `class Component` with `renderVals()` returning values consumed by a `{{ }}` template. Treat this as a normal component with local state.
- `android-frame.jsx` is only a **device bezel/status-bar mock** for previewing on a fake Android frame. **Do not port it** — the real app runs full-screen on a real device.
- Styling in the prototype is all inline. Exact values are documented under Design Tokens.

## Fidelity
**High-fidelity (hifi).** Final colors, typography, spacing, and interactions. Recreate the UI faithfully using the codebase's existing component library, then apply the tokens below.

---

## Global Layout & Navigation (Bidan App)
- Single phone-width column. Content area scrolls; a fixed **bottom tab bar** (height 66px) stays pinned.
- **4 tabs:** Beranda (Home), Kunjungan (Visits), Stok (Stock), Rekap (Recap). Active tab = accent blue `#2f6ce0`; inactive = `#9aa7b8`. Each tab: 22px line icon + 10.5px/700 label, stacked, centered.
- **Sub-screens/modals** overlay the content area (position absolute, inset 0) with their own back button: Catat Kunjungan, Cari Obat, Jadwal & Agenda, Form Agenda. Two are **bottom sheets** (slide up from bottom, dimmed backdrop `rgba(16,22,30,.4)`): Tambah Stok, Atur Status Keluar.

## Screens / Views (Bidan App)

### 1. Beranda (Home)
- **Purpose:** at-a-glance status + quick actions.
- **Header:** eyebrow "KLINIK" (11px/600, letter-spacing .12em, `#8a97a8`, uppercase) + "Bidan Pit" (22px/800, `-.02em`). Right: 42px circle avatar `#dbe6f7`, letter "P" in `#2f6ce0`.
- **Status hero card** (two states):
  - **Present (hadir):** background `#2f6ce0`, white text, radius 20, shadow `0 12px 26px rgba(47,108,224,.28)`. Row: "Status Bidan" (12.5/600, .85 opacity) + an ON toggle (52×30 track `rgba(255,255,255,.32)`, 24px white knob on right). Title "Sedang di klinik" (20/800). Subtitle "Ketuk untuk atur status keluar". **Tapping the card opens the "Atur Status Keluar" bottom sheet.**
  - **Away (tidak):** white card, border `#e2e7ef`, radius 20. Row: "Status Bidan" (`#6b7688`) + an OFF toggle (track `#d6dbe4`, knob left) — **tapping the toggle sets status back to present**. Title "Tidak di klinik" (`#16202e`). Optional note line (13/600 `#5a6675`). A pill (`#f6f8fc`, radius 12) showing either "Perkiraan kembali HH.MM" (13/700) + countdown "± X menit/jam lagi" (12/600 `#8a97a8`), OR "Waktu kembali belum dipastikan". Row of 3 buttons: **+15 mnt**, **+30 mnt**, **Ubah** (each `#f0f4f9`/`#3b4a5e`, 12.5/700, radius 10). Full-width green button "Saya sudah kembali" (`#1f9d57`, white, radius 12).
- **Web-status link:** pill `#eef4ff` border `#d9e5fb` radius 14: green dot + "Status ini tampil di halaman web pasien" (12.5/600 `#2f6ce0`) + chevron. Opens the Patient Web Page.
- **Quick stats:** two cards side by side (flex, gap 10). Each white, border `#e7edf5`, radius 14: big number (22/800) + label (11.5/600 `#8a97a8`). Left = "Kunjungan hari ini" (count). Right = "Obat menipis" (count, number in `#d64545`).
- **Stok menipis (low stock) list:** section header "Stok menipis" (15/700) + "Lihat stok" link (`#2f6ce0`). Cards (white, border `#e7edf5`, radius 14): colored dot + name (14/600) + category (12 `#8a97a8`) + qty (15/800, colored) + unit. Dot/qty color: `#d64545` (danger) or `#e0a52a` (warning). **Sorted ascending by quantity.**
- **Jadwal & agenda card:** header "Jadwal & agenda" (15/700) + "Kelola" link. Up to 3 upcoming events as cards (dot + title 14/600 + when-string 12 `#8a97a8`). Empty state: dashed box "Belum ada agenda mendatang". "Kelola" opens the Jadwal screen.
- **Catat Kunjungan button:** full-width `#16202e`, white, radius 16, "+ Catat Kunjungan" (15/700). Opens the visit form.

### 2. Catat Kunjungan (Add Visit — full-screen modal)
- Header: back button (36px circle `#f0f4f9`, "‹") + "Catat Kunjungan" (17/800).
- Fields:
  - **Nama pasien** — text input (radius 12, border `#dde4ee`, 15px). Required.
  - **Umur** — numeric input, 90px wide, + "tahun" suffix. Digits only.
  - **Gejala / keluhan · opsional** — multiline textarea (3 rows, resize none). Optional.
  - **Obat yang diberikan** — label + "+ Tambah obat" link (opens Cari Obat). Added items render as rows: name (14/600) + unit (11.5 `#8a97a8`) + a stepper (−/qty/+, bordered `#e2e7ef`, radius 10) + remove "×". Empty state: dashed box "Belum ada obat. Ketuk 'Tambah obat'."
- Footer: "Simpan Kunjungan" button. Enabled (blue `#2f6ce0`) only when name is non-empty AND ≥1 medicine; otherwise disabled color `#b9c6d8`.
- **On save:** create visit with timestamp; **decrement stock** of each given medicine by its qty (clamped at 0); navigate to Kunjungan tab.

### 3. Cari Obat (Medicine Search — full-screen modal)
- Header: back button + "Cari Obat". Back returns to the visit form (if opened from a visit) or closes (if opened from stock).
- Search input (`#f6f8fc`) filters the master catalog by substring.
- Results: rows (white, border `#e7edf5`, radius 12): 36px rounded-square avatar with first letter, tinted by category (Tablet `#e7eefb`/`#2f6ce0`, Sirup `#e6f3ee`/`#2f8a63`, Sachet `#f3edfb`/`#7a54c0`), name (14.5/600), sub = category + stock note ("· stok N butir" or "· belum ada stok"), and a "+".
- When invoked from **stock** context and the query (≥3 chars) doesn't match anything, show a dashed "Tambah '<query>' sebagai obat baru" row (creates a new Tablet medicine).
- Selecting a medicine: from a visit → adds it to the visit draft (qty 1, or +1 if already added) and returns to the form. From stock → opens the "Tambah Stok" bottom sheet for that medicine.

### 4. Tambah Stok (Add Stock — bottom sheet)
- Drag handle, title "Tambah stok", medicine name (20/800), "Stok sekarang: N <baseUnit>".
- **Satuan (unit) chips** — depend on medicine category (see Unit Conversion). Selected chip = `#16202e`/white or `#e7eefb`/`#2f6ce0` border `#2f6ce0`; unselected = white/`#6b7688` border `#e2e7ef`.
- **Jumlah** — stepper: −  [numeric input, 20/800, centered]  +.
- **Live preview line** (13/700 `#2f6ce0`): "= +<count×multiplier> <baseUnit> (total jadi <new total>)".
- Buttons: "Batal" (`#f0f4f9`) and "Simpan Stok" (`#2f6ce0`, flex 2). On save, add the converted amount to the medicine's quantity (creating the medicine if new); go to Stok tab.

### 5. Stok Obat (Stock list tab)
- Header "Stok Obat" (22/800) + "+ Tambah" pill button (`#2f6ce0`, opens Cari Obat in stock context). Subtitle "N jenis obat terdaftar".
- List, **sorted alphabetically**: cards (white, border `#e7edf5`, radius 14): status dot (`#d64545` if low, else `#3fbf8f`) + name (14.5/600) + category (12 `#8a97a8`) + qty (16/800, `#d64545` if low else `#16202e`) + base unit.

### 6. Rekap (Recap tab)
- Title "Rekap" + a segmented control (`#e9eef5`, radius 14, active segment white): **Kunjungan** / **Stok Menipis**.
- **Kunjungan:** two stat cards — "Total kunjungan" (`#16202e` bg) and "Obat dikeluarkan" (`#2f6ce0` bg), each 26/800 number. Then "Ringkasan per hari": rows grouped by day (label "Hari ini"/"Kemarin"/"D Month", first-name list, and count in `#2f6ce0` 18/800), newest first.
- **Stok Menipis:** subtitle "N obat perlu segera dibeli · urut dari terkecil". Rows: a 52px tinted square showing qty (20/800) + unit stacked (tint `#fdeaea`+`#d64545` danger, or `#fdf3df`+`#e0a52a` warning) + name (14.5/700) + "category · status" where status = "segera dibeli" (danger) / "perlu diperhatikan" (warning). **Ascending by qty.**

### 7. Kunjungan (Visits tab)
- Header "Kunjungan" + round "+" FAB (`#2f6ce0`, opens the visit form). Subtitle "N total · M hari ini".
- Visit cards (white, border `#e7edf5`, radius 16), **newest first**: 40px rounded-square avatar `#dbe6f7`/`#2f6ce0` with patient initial + name (15/700) + age "N tahun" (12.5 `#8a97a8`) + date label (right, 11.5/600). If gejala present: line "Gejala: <text>" (12.5 `#5a6675`). Divider, then medicine chips ("<name> · <qty> <unit>", `#f0f4f9` pills).
- Date label format: "Hari ini HH.MM", "Kemarin HH.MM", or "D Mon".

### 8. Jadwal & Agenda (Schedule — full-screen modal)
- Header: back + "Jadwal & Agenda" + "+ Tambah" pill (opens Form Agenda blank).
- Event rows (tappable → edit): 10px dot (past `#c0c8d4`, all-day `#e0a52a`, timed `#2f6ce0`) + title (14.5/700) + when-string (12.5 `#8a97a8`) + chevron. **Sorted ascending by start.**
- When-string formats: all-day same day → "Sen, 12 Jul · seharian"; all-day range → "Sen, 12 Jul – Rab, 14 Jul · N hari"; timed same day → "Sen, 12 Jul · 13.00–15.00"; timed range → full "date HH.MM – date HH.MM".

### 9. Form Agenda (Add/Edit Event — full-screen modal)
- Header: back (returns to Jadwal) + "Agenda Baru" / "Ubah Agenda".
- Fields: **Judul agenda** (text, required). **"Seharian / beberapa hari"** toggle (blue when on). **Mulai** = date input (+ time input when NOT all-day). **Selesai** = date input (+ time input when NOT all-day). If editing: red "Hapus agenda" button (`#fbeceb`/`#c25146`).
- Footer: "Simpan Agenda" (enabled blue when title + start date present, else `#b9c6d8`). Save upserts the event; delete removes it. Both re-broadcast to the web page.

## Screens / Views (Patient Web Page)
Single scrolling column, max-width 460px, background `#eef1f6`.
- **Brand header:** 46px blue rounded-square "P" + "Klinik Bidan Pit" (19/800) + "Praktik Bidan Mandiri · melayani 24 jam" (12.5 `#8a97a8`).
- **Status hero (two states):**
  - **Present:** green card `#1f9d57`, white text, shadow `0 18px 40px rgba(31,157,87,.28)`. Badge "● SEDANG BUKA" (white dot with a pulsing box-shadow animation, 1.8s infinite). "Bidan sedang di klinik" (32/800). Reassurance line. "Diperbarui <relative time>".
  - **Away:** white card. Badge "● SEDANG TIDAK ADA" (grey). "Bidan sedang tidak di tempat" (32/800). Optional note line. Inner box `#f5f8fc`: either "PERKIRAAN KEMBALI / Pukul HH.MM / ± X lagi" or "Waktu kembali belum dapat dipastikan". "Diperbarui <relative time>".
- **Jadwal Bidan card:** "Jadwal Bidan" + "Agenda mendatang saat bidan berhalangan". Upcoming events (dot + title 14.5/700 + when 12.5 `#8a97a8`). Empty → "Tidak ada agenda mendatang — bidan siap melayani." Shows only events whose end ≥ now, ascending, max 4.
- **Lokasi & Kontak card:** address (placeholder) + green "Hubungi via WhatsApp" button (`#25a15a`). **Wire the real WhatsApp deep link / phone number.**
- Footer note.

## Interactions & Behavior
- **Status toggle:** present→away opens the sheet (must set note + optional estimate); away→present clears note & estimate instantly ("kill"). +15/+30 extend the estimate (or start from now if none).
- **Live sync:** any status/schedule change is broadcast; the web page reflects it near-instantly and also re-reads periodically to refresh the relative time / countdown.
- **Stock is the source of truth for low-stock alerts and recaps**; saving a visit reduces it.
- **Forms:** disabled save buttons until required fields valid; numeric-only enforced on age & stock count.
- **Animations:** bottom sheets slide up; the "SEDANG BUKA" badge dot pulses (CSS keyframes). Keep transitions subtle.

## State Management
Bidan App state (persist locally in prototype; back with real API + DB in production):
- `bidanHadir: boolean`
- `awayNote: string`, `awayUntil: number|null` (epoch ms)
- `events: [{ id, title, allDay:boolean, startTs:number, endTs:number }]`
- `medicines: [{ id, name, cat:'Tablet'|'Sirup'|'Sachet', qty:number }]` — qty is always in the **base unit**.
- `visits: [{ id, name, age:number, gejala:string, ts:number, items:[{ name, qty, unit }] }]`
- Ephemeral: current screen/tab, open modal, search query, draft objects for visit/away/event/add-stock.

**Broadcast payload → patient page** (`klinik_bidan_status`):
```json
{ "hadir": true, "awayNote": "", "awayUntil": null, "ts": 1720000000000,
  "clinic": "Klinik Bidan Pit",
  "events": [{ "id": "e2", "title": "Keluar kota", "allDay": true, "startTs": 0, "endTs": 0 }] }
```
Replace the localStorage/BroadcastChannel transport with real endpoints; keep this shape.

## Unit Conversion (important business rule)
Medicines are counted in a **base unit**; stock is added in packaging units that convert to base:
- **Tablet** → base `butir`. Add by **Box (100 butir)**, **Strip (10 butir)**, or **Butir (1)**. Low-stock threshold ≤ 20.
- **Sirup** → base `botol`. Add by **Botol (1)** only. Threshold ≤ 5.
- **Sachet** → base `sachet`. Add by **Box (100 sachet)** or **Sachet (1)**. Threshold ≤ 20.
"Danger" (red) when qty ≤ threshold × 0.5, else "warning" (amber).

## Design Tokens
**Colors**
- Accent blue `#2f6ce0` (hover `#1f52b8`); ink/dark `#16202e`; near-black surface `#16202e`.
- Text: primary `#16202e`, secondary `#5a6675`, muted `#8a97a8` / `#9aa7b8`, label `#6b7688`.
- Surfaces: app bg `#f5f8fc`, card white `#ffffff`, card border `#e7edf5` / `#e2e7ef`, subtle fill `#f0f4f9` / `#f6f8fc`, input border `#dde4ee`.
- Success/present green `#1f9d57`; healthy dot `#3fbf8f`; WhatsApp green `#25a15a`.
- Danger `#d64545` (bg `#fdeaea`), delete `#c25146` (bg `#fbeceb`); warning `#e0a52a` (bg `#fdf3df`).
- Category tints: Tablet `#e7eefb`/`#2f6ce0`, Sirup `#e6f3ee`/`#2f8a63`, Sachet `#f3edfb`/`#7a54c0`.
- Patient page bg `#eef1f6`.

**Typography** — Plus Jakarta Sans (Google Fonts), weights 400–800.
- Screen title 22/800 (-.02em); big numbers 26–32/800 (-.02/-.03em); section header 15/700; body 14–15/500–600; label/eyebrow 11–12.5/600–700; tab label 10.5/700.

**Radii** — cards 14–24; inputs/buttons 12–16; pills/toggles 16–30; sheets 22 (top corners).

**Shadows** — blue hero `0 12px 26px rgba(47,108,224,.28)`; green hero `0 18px 40px rgba(31,157,87,.28)`; FAB `0 6px 16px rgba(47,108,224,.3)`.

**Toggle** — 52×30 track, 24px knob (3px inset). Event toggle 46×27, 21px knob.

## Assets
- **Font:** Plus Jakarta Sans via Google Fonts.
- **Icons:** simple inline line SVGs for tab bar (home, clipboard, box, bar-chart) and a filled WhatsApp glyph. Replace with the codebase's icon set.
- **Avatar/logo:** the "P" monogram is a placeholder — swap for the clinic's real logo if available.
- No raster image assets. No copyrighted/brand assets used.

## Files
- `Klinik Bidan Pit App.dc.html` — the midwife's app (all screens above).
- `Status Bidan Web.dc.html` — the public patient status/schedule page.
- `android-frame.jsx` — device-bezel preview mock only; **do not port**.
- Open either HTML file directly in a browser to see the live prototype. To see the live sync, open both in two tabs and toggle status in the app.
