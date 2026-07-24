import express from 'express';
import cors from 'cors';
import { fileURLToPath } from 'url';
import path from 'path';
import 'dotenv/config';
import { pool } from './db.js';
import { CAT, CATALOG, baseUnit } from './catalog.js';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const app = express();
app.use(cors());
app.use(express.json());

const uid = (p) => p + Date.now().toString(36) + Math.random().toString(36).slice(2, 6);

// ---------- helper query ----------
async function getStatus() {
  const [rows] = await pool.query('SELECT * FROM clinic_status WHERE id = 1');
  if (!rows.length) {
    const ts = Date.now();
    await pool.query(
      'INSERT INTO clinic_status (id, bidan_hadir, away_note, away_until, updated_ts, clinic) VALUES (1,1,"",NULL,?,?)',
      [ts, 'Klinik Bidan Pit']
    );
    return { bidanHadir: true, awayNote: '', awayUntil: null, updatedTs: ts, clinic: 'Klinik Bidan Pit' };
  }
  const r = rows[0];
  return {
    bidanHadir: !!r.bidan_hadir,
    awayNote: r.away_note || '',
    awayUntil: r.away_until != null ? Number(r.away_until) : null,
    updatedTs: Number(r.updated_ts),
    clinic: r.clinic,
  };
}

async function getEvents() {
  const [rows] = await pool.query('SELECT * FROM events ORDER BY start_ts ASC');
  return rows.map((e) => ({
    id: e.id, title: e.title, allDay: !!e.all_day,
    startTs: Number(e.start_ts), endTs: Number(e.end_ts),
  }));
}

async function getMedicines() {
  const [rows] = await pool.query('SELECT * FROM medicines ORDER BY name ASC');
  return rows.map((m) => ({ id: m.id, name: m.name, cat: m.cat, qty: m.qty, baseUnit: baseUnit(m.cat) }));
}

async function getVisits() {
  const [visits] = await pool.query('SELECT * FROM visits ORDER BY ts DESC');
  const [items] = await pool.query('SELECT * FROM visit_items');
  const byVisit = {};
  for (const it of items) {
    (byVisit[it.visit_id] ||= []).push({ name: it.name, qty: it.qty, unit: it.unit });
  }
  return visits.map((v) => ({
    id: v.id, name: v.name, age: v.age, gejala: v.gejala || '',
    ts: Number(v.ts), items: byVisit[v.id] || [],
  }));
}

// payload untuk halaman web pasien (bentuk sesuai README: klinik_bidan_status)
async function statusPayload() {
  const s = await getStatus();
  const events = await getEvents();
  return {
    hadir: s.bidanHadir, awayNote: s.awayNote, awayUntil: s.awayUntil,
    ts: s.updatedTs, clinic: s.clinic, events,
  };
}

// ---------- API ----------
app.get('/api/health', (req, res) => res.json({ ok: true }));

// katalog + aturan satuan (untuk layar Cari Obat / Tambah Stok)
app.get('/api/catalog', (req, res) => res.json({ catalog: CATALOG, cat: CAT }));

// snapshot penuh untuk bootstrap Bidan App
app.get('/api/state', async (req, res, next) => {
  try {
    const [status, events, medicines, visits] = await Promise.all([
      getStatus(), getEvents(), getMedicines(), getVisits(),
    ]);
    res.json({ status, events, medicines, visits });
  } catch (e) { next(e); }
});

// ubah status bidan (hadir / keluar / extend)
app.put('/api/status', async (req, res, next) => {
  try {
    const { bidanHadir, awayNote = '', awayUntil = null } = req.body;
    const ts = Date.now();
    await pool.query(
      `INSERT INTO clinic_status (id, bidan_hadir, away_note, away_until, updated_ts)
       VALUES (1,?,?,?,?)
       ON DUPLICATE KEY UPDATE bidan_hadir=VALUES(bidan_hadir), away_note=VALUES(away_note),
         away_until=VALUES(away_until), updated_ts=VALUES(updated_ts)`,
      [bidanHadir ? 1 : 0, awayNote || '', awayUntil ?? null, ts]
    );
    res.json(await getStatus());
  } catch (e) { next(e); }
});

// upsert agenda
app.post('/api/events', async (req, res, next) => {
  try {
    const { title, allDay = false, startTs, endTs } = req.body;
    if (!title || !title.trim() || !startTs) return res.status(400).json({ error: 'title & startTs wajib' });
    const id = uid('e');
    const end = endTs && endTs >= startTs ? endTs : startTs;
    await pool.query('INSERT INTO events (id,title,all_day,start_ts,end_ts) VALUES (?,?,?,?,?)',
      [id, title.trim(), allDay ? 1 : 0, startTs, end]);
    res.status(201).json({ id, title: title.trim(), allDay: !!allDay, startTs, endTs: end });
  } catch (e) { next(e); }
});

app.put('/api/events/:id', async (req, res, next) => {
  try {
    const { title, allDay = false, startTs, endTs } = req.body;
    if (!title || !title.trim() || !startTs) return res.status(400).json({ error: 'title & startTs wajib' });
    const end = endTs && endTs >= startTs ? endTs : startTs;
    await pool.query('UPDATE events SET title=?, all_day=?, start_ts=?, end_ts=? WHERE id=?',
      [title.trim(), allDay ? 1 : 0, startTs, end, req.params.id]);
    res.json({ id: req.params.id, title: title.trim(), allDay: !!allDay, startTs, endTs: end });
  } catch (e) { next(e); }
});

app.delete('/api/events/:id', async (req, res, next) => {
  try {
    await pool.query('DELETE FROM events WHERE id=?', [req.params.id]);
    res.json({ ok: true });
  } catch (e) { next(e); }
});

// tambah stok (menambah obat baru bila belum ada). amount = jumlah dalam base unit.
app.post('/api/medicines/stock', async (req, res, next) => {
  try {
    const { id, name, cat = 'Tablet', amount } = req.body;
    const add = parseInt(amount, 10) || 0;
    if (!name || !name.trim() || add <= 0) return res.status(400).json({ error: 'name & amount(>0) wajib' });
    // cari berdasarkan id atau nama
    const [found] = await pool.query('SELECT * FROM medicines WHERE id=? OR LOWER(name)=LOWER(?) LIMIT 1',
      [id || '', name.trim()]);
    if (found.length) {
      await pool.query('UPDATE medicines SET qty = qty + ? WHERE id=?', [add, found[0].id]);
      return res.json({ id: found[0].id });
    }
    const newId = id || uid('m');
    await pool.query('INSERT INTO medicines (id,name,cat,qty) VALUES (?,?,?,?)',
      [newId, name.trim(), cat, add]);
    res.status(201).json({ id: newId });
  } catch (e) { next(e); }
});

// catat kunjungan + kurangi stok (transaksi)
app.post('/api/visits', async (req, res, next) => {
  const conn = await pool.getConnection();
  try {
    const { name, age = 0, gejala = '', items = [] } = req.body;
    if (!name || !name.trim() || !Array.isArray(items) || items.length === 0)
      return res.status(400).json({ error: 'name & minimal 1 obat wajib' });
    await conn.beginTransaction();
    const id = uid('v');
    await conn.query('INSERT INTO visits (id,name,age,gejala,ts) VALUES (?,?,?,?,?)',
      [id, name.trim(), parseInt(age, 10) || 0, (gejala || '').trim(), Date.now()]);
    for (const it of items) {
      await conn.query('INSERT INTO visit_items (visit_id,name,qty,unit) VALUES (?,?,?,?)',
        [id, it.name, it.qty, it.unit]);
      // kurangi stok berdasarkan nama, clamp di 0
      await conn.query('UPDATE medicines SET qty = GREATEST(0, qty - ?) WHERE LOWER(name)=LOWER(?)',
        [it.qty, it.name]);
    }
    await conn.commit();
    res.status(201).json({ id });
  } catch (e) {
    await conn.rollback();
    next(e);
  } finally {
    conn.release();
  }
});

// ---------- Halaman web pasien ----------
app.get('/api/public/status', async (req, res, next) => {
  try { res.json(await statusPayload()); } catch (e) { next(e); }
});

app.use(express.static(path.join(__dirname, '..', 'public')));
app.get('/', (req, res) => res.sendFile(path.join(__dirname, '..', 'public', 'index.html')));

// ---------- error handler ----------
app.use((err, req, res, next) => {
  console.error(err);
  res.status(500).json({ error: err.message || 'server error' });
});

const PORT = Number(process.env.PORT) || 4000;
app.listen(PORT, () => console.log(`Klinik Bidan Pit API di http://localhost:${PORT}`));
