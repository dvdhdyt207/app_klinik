// Package catalog memuat aturan bisnis konversi satuan & katalog obat
// (dari README design handoff). qty obat SELALU disimpan dalam base unit.
package catalog

// Unit = satuan kemasan untuk menambah stok (mis. "Box (100)" -> mult 100).
type Unit struct {
	Label string `json:"label"`
	Mult  int    `json:"mult"`
}

// Rule = aturan per kategori obat: base unit, satuan kemasan, & ambang menipis.
type Rule struct {
	Base  string `json:"base"`
	Units []Unit `json:"units"`
	Thr   int    `json:"thr"`
}

// CAT: aturan per kategori. Kunci = kategori obat.
var CAT = map[string]Rule{
	"Tablet": {Base: "butir", Units: []Unit{{"Box (100)", 100}, {"Strip (10)", 10}, {"Butir", 1}}, Thr: 20},
	"Sirup":  {Base: "botol", Units: []Unit{{"Botol", 1}}, Thr: 5},
	"Sachet": {Base: "sachet", Units: []Unit{{"Box (100)", 100}, {"Sachet", 1}}, Thr: 20},
}

func rule(cat string) Rule {
	if r, ok := CAT[cat]; ok {
		return r
	}
	return CAT["Tablet"]
}

// BaseUnit -> base unit (butir/botol/sachet) untuk kategori.
func BaseUnit(cat string) string { return rule(cat).Base }

// Threshold -> ambang menipis untuk kategori.
func Threshold(cat string) int { return rule(cat).Thr }

// IsLow: menipis bila qty <= threshold.
func IsLow(cat string, qty int) bool { return qty <= Threshold(cat) }

// IsDanger (merah): qty <= threshold*0.5, selain itu warning (amber).
func IsDanger(cat string, qty int) bool { return float64(qty) <= float64(Threshold(cat))*0.5 }

// Item = entri katalog master obat untuk fitur "Cari Obat".
type Item struct {
	Name string `json:"name"`
	Cat  string `json:"cat"`
}

// CATALOG sengaja kosong. Dulu ia memuat 14 obat contoh dari design handoff,
// dan itu keliru: daftar obat adalah DATA milik klinik, bukan aturan bisnis.
// Obat karangan yang tidak pernah dimiliki klinik ikut muncul di pencarian dan
// tidak bisa dihapus tanpa deploy ulang.
//
// Sumber daftar obat sekarang tabel `medicines` — apa yang benar-benar
// distok bidan. Lihat PickModal.vue.
var CATALOG = []Item{}

// ValidCat: kategori yang dikenal. Dipakai sebelum menulis ke kolom `cat`,
// yang bertipe ENUM di MySQL — nilai di luar daftar ditolak database dengan
// pesan yang tidak menjelaskan apa-apa bagi pemakainya.
func ValidCat(cat string) bool {
	_, ok := CAT[cat]
	return ok
}
