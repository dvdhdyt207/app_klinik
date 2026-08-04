package auth

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// KeyedLimiter memberi satu jatah laju terpisah untuk tiap kunci (alamat IP,
// nama pengguna, …).
//
// Entri lama dibuang berkala dan jumlahnya dibatasi. Tanpa itu, penyerang yang
// terus berganti kunci membuat peta ini tumbuh tanpa henti — pembatas laju
// yang seharusnya melindungi justru berubah menjadi cara menghabiskan memori.
type KeyedLimiter struct {
	mu    sync.Mutex
	entri map[string]*entri

	r     rate.Limit
	burst int
	ttl   time.Duration
	maks  int

	berhenti chan struct{}
	sekali   sync.Once
}

type entri struct {
	lim      *rate.Limiter
	terakhir time.Time
}

// NewKeyedLimiter membuat pembatas dengan laju r dan lonjakan burst.
// Entri yang tidak disentuh selama ttl dibuang; jumlahnya dibatasi maks.
func NewKeyedLimiter(r rate.Limit, burst int, ttl time.Duration, maks int) *KeyedLimiter {
	l := &KeyedLimiter{
		entri: make(map[string]*entri),
		r:     r, burst: burst, ttl: ttl, maks: maks,
		berhenti: make(chan struct{}),
	}
	go l.sapuBerkala()
	return l
}

// Allow memotong satu jatah untuk key dan melaporkan apakah masih boleh lewat.
func (l *KeyedLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	e, ada := l.entri[key]
	if !ada {
		if len(l.entri) >= l.maks {
			l.buangTertuaTerkunci()
		}
		e = &entri{lim: rate.NewLimiter(l.r, l.burst)}
		l.entri[key] = e
	}
	e.terakhir = time.Now()
	return e.lim.Allow()
}

// Reset menghapus jatah sebuah kunci — dipakai setelah login berhasil, supaya
// percobaan gagal sebelumnya tidak terus membebani pemakai yang sah.
func (l *KeyedLimiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.entri, key)
}

// Stop menghentikan penyapu latar.
func (l *KeyedLimiter) Stop() {
	l.sekali.Do(func() { close(l.berhenti) })
}

func (l *KeyedLimiter) sapuBerkala() {
	t := time.NewTicker(l.ttl)
	defer t.Stop()
	for {
		select {
		case <-l.berhenti:
			return
		case now := <-t.C:
			l.mu.Lock()
			for k, e := range l.entri {
				if now.Sub(e.terakhir) > l.ttl {
					delete(l.entri, k)
				}
			}
			l.mu.Unlock()
		}
	}
}

// buangTertuaTerkunci membuang satu entri paling lama tidak dipakai.
// Pemanggil harus sudah memegang l.mu.
func (l *KeyedLimiter) buangTertuaTerkunci() {
	var kunciTua string
	var waktuTua time.Time
	for k, e := range l.entri {
		if kunciTua == "" || e.terakhir.Before(waktuTua) {
			kunciTua, waktuTua = k, e.terakhir
		}
	}
	if kunciTua != "" {
		delete(l.entri, kunciTua)
	}
}
