package entity

import "caltrack/domain/vo"

// RecordPfc は食事記録のPFC情報（グラム）を表すEntity
type RecordPfc struct {
	id       vo.RecordPfcID
	recordID vo.RecordID
	pfc      vo.Pfc
}

// NewRecordPfc は新しいRecordPfcを生成する
func NewRecordPfc(recordID vo.RecordID, protein, fat, carbs float64) *RecordPfc {
	return &RecordPfc{
		id:       vo.NewRecordPfcID(),
		recordID: recordID,
		pfc:      vo.NewPfc(protein, fat, carbs),
	}
}

// ReconstructRecordPfc はDBからRecordPfcを復元する
func ReconstructRecordPfc(idStr, recordIDStr string, protein, fat, carbs float64) *RecordPfc {
	return &RecordPfc{
		id:       vo.ReconstructRecordPfcID(idStr),
		recordID: vo.ReconstructRecordID(recordIDStr),
		pfc:      vo.NewPfc(protein, fat, carbs),
	}
}

// ID はRecordPfcIDを返す
func (r *RecordPfc) ID() vo.RecordPfcID {
	return r.id
}

// RecordID はRecordIDを返す
func (r *RecordPfc) RecordID() vo.RecordID {
	return r.recordID
}

// Pfc はPfc情報を返す
func (r *RecordPfc) Pfc() vo.Pfc {
	return r.pfc
}

// Protein はタンパク質(g)を返す
func (r *RecordPfc) Protein() float64 {
	return r.pfc.Protein()
}

// Fat は脂質(g)を返す
func (r *RecordPfc) Fat() float64 {
	return r.pfc.Fat()
}

// Carbs は炭水化物(g)を返す
func (r *RecordPfc) Carbs() float64 {
	return r.pfc.Carbs()
}
