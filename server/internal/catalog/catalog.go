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

// CATALOG: daftar obat master (untuk pencarian & buat obat baru).
var CATALOG = []Item{
	{"Paracetamol 500mg", "Tablet"},
	{"Paracetamol sirup", "Sirup"},
	{"Amoxicillin 500mg", "Tablet"},
	{"Amoxicillin sirup", "Sirup"},
	{"Ibuprofen 400mg", "Tablet"},
	{"Asam Mefenamat 500mg", "Tablet"},
	{"Cetirizine 10mg", "Tablet"},
	{"CTM 4mg", "Tablet"},
	{"Antasida", "Tablet"},
	{"Vitamin B Complex", "Tablet"},
	{"Vitamin C 500mg", "Tablet"},
	{"Domperidone 10mg", "Tablet"},
	{"Dexamethasone 0.5mg", "Tablet"},
	{"Oralit", "Sachet"},
}
