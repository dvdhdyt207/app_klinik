// Seed data awal (meniru prototipe design handoff). Idempotent: mengosongkan tabel dulu.
// Jalankan: npm run seed
import { pool } from '../src/db.js';

const now = Date.now();
const DAY = 86400000;
const sot = new Date(); sot.setHours(0, 0, 0, 0);
const at = (base, h, m) => { const d = new Date(base); d.setHours(h, m, 0, 0); return d.getTime(); };

const events = [
  { id: 'e1', title: 'Posyandu Desa Sukamaju', allDay: 0, startTs: at(sot, 13, 0), endTs: at(sot, 15, 0) },
  { id: 'e2', title: 'Keluar kota — acara keluarga', allDay: 1, startTs: sot.getTime() + DAY, endTs: sot.getTime() + 3 * DAY },
];

const medicines = [
  { id: 'm1', name: 'Paracetamol sirup', cat: 'Sirup', qty: 2 },
  { id: 'm2', name: 'Amoxicillin 500mg', cat: 'Tablet', qty: 6 },
  { id: 'm3', name: 'Vitamin B Complex', cat: 'Tablet', qty: 14 },
  { id: 'm4', name: 'Paracetamol 500mg', cat: 'Tablet', qty: 240 },
  { id: 'm5', name: 'Ibuprofen 400mg', cat: 'Tablet', qty: 85 },
  { id: 'm6', name: 'Cetirizine 10mg', cat: 'Tablet', qty: 40 },
  { id: 'm7', name: 'Amoxicillin sirup', cat: 'Sirup', qty: 9 },
  { id: 'm8', name: 'Antasida', cat: 'Tablet', qty: 120 },
  { id: 'm9', name: 'Oralit', cat: 'Sachet', qty: 60 },
];

const visits = [
  { id: 'v1', name: 'Siti Aminah', age: 27, ts: now - 2 * 3600000, gejala: 'Demam 2 hari, batuk kering', items: [{ name: 'Paracetamol sirup', qty: 1, unit: 'botol' }] },
  { id: 'v2', name: 'Budi Santoso', age: 45, ts: now - 5 * 3600000, gejala: 'Nyeri lambung, mual', items: [{ name: 'Amoxicillin 500mg', qty: 10, unit: 'butir' }, { name: 'Vitamin B Complex', qty: 10, unit: 'butir' }] },
  { id: 'v3', name: 'Rina Wati', age: 5, ts: now - DAY - 3600000, gejala: 'Panas tinggi, rewel', items: [{ name: 'Paracetamol 500mg', qty: 10, unit: 'butir' }] },
  { id: 'v4', name: 'Joko Priyono', age: 60, ts: now - 2 * DAY, gejala: '', items: [{ name: 'Antasida', qty: 20, unit: 'butir' }, { name: 'Vitamin C 500mg', qty: 10, unit: 'butir' }] },
];

async function main() {
  const conn = await pool.getConnection();
  try {
    await conn.query('SET FOREIGN_KEY_CHECKS = 0');
    await conn.query('TRUNCATE TABLE visit_items');
    await conn.query('TRUNCATE TABLE visits');
    await conn.query('TRUNCATE TABLE medicines');
    await conn.query('TRUNCATE TABLE events');
    await conn.query('TRUNCATE TABLE clinic_status');
    await conn.query('SET FOREIGN_KEY_CHECKS = 1');

    await conn.query(
      'INSERT INTO clinic_status (id,bidan_hadir,away_note,away_until,updated_ts,clinic) VALUES (1,1,"",NULL,?,?)',
      [now, 'Klinik Bidan Pit']
    );
    for (const e of events)
      await conn.query('INSERT INTO events (id,title,all_day,start_ts,end_ts) VALUES (?,?,?,?,?)',
        [e.id, e.title, e.allDay, e.startTs, e.endTs]);
    for (const m of medicines)
      await conn.query('INSERT INTO medicines (id,name,cat,qty) VALUES (?,?,?,?)',
        [m.id, m.name, m.cat, m.qty]);
    for (const v of visits) {
      await conn.query('INSERT INTO visits (id,name,age,gejala,ts) VALUES (?,?,?,?,?)',
        [v.id, v.name, v.age, v.gejala, v.ts]);
      for (const it of v.items)
        await conn.query('INSERT INTO visit_items (visit_id,name,qty,unit) VALUES (?,?,?,?)',
          [v.id, it.name, it.qty, it.unit]);
    }
    console.log('Seed selesai: ' + events.length + ' agenda, ' + medicines.length + ' obat, ' + visits.length + ' kunjungan.');
  } finally {
    conn.release();
    await pool.end();
  }
}

main().catch((e) => { console.error(e); process.exit(1); });
