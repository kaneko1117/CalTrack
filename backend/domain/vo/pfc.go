package vo

// Pfc はPFC(タンパク質・脂質・炭水化物)を表すValue Object
type Pfc struct {
	protein float64
	fat     float64
	carbs   float64
}

// NewPfc は新しいPfcを生成する
func NewPfc(protein, fat, carbs float64) Pfc {
	return Pfc{protein: protein, fat: fat, carbs: carbs}
}

// Protein はタンパク質(g)を返す
func (p Pfc) Protein() float64 {
	return p.protein
}

// Fat は脂質(g)を返す
func (p Pfc) Fat() float64 {
	return p.fat
}

// Carbs は炭水化物(g)を返す
func (p Pfc) Carbs() float64 {
	return p.carbs
}
