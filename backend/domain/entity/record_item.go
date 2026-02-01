package entity

import (
	"caltrack/domain/vo"
)

// RecordItem はカロリー記録の明細を表すエンティティ
type RecordItem struct {
	id       vo.RecordItemID
	recordID vo.RecordID
	name     vo.ItemName
	calories vo.Calories
}

// NewRecordItem は新しいRecordItemを生成する
func NewRecordItem(recordID vo.RecordID, nameStr string, caloriesVal int) (*RecordItem, []error) {
	var errs []error

	name, err := vo.NewItemName(nameStr)
	errs = appendIfErr(errs, err)

	calories, err := vo.NewCalories(caloriesVal)
	errs = appendIfErr(errs, err)

	if len(errs) > 0 {
		return nil, errs
	}

	return &RecordItem{
		id:       vo.NewRecordItemID(),
		recordID: recordID,
		name:     name,
		calories: calories,
	}, nil
}

// ReconstructRecordItem はDBからRecordItemを復元する
func ReconstructRecordItem(
	idStr string,
	recordIDStr string,
	nameStr string,
	caloriesVal int,
) *RecordItem {
	return &RecordItem{
		id:       vo.ReconstructRecordItemID(idStr),
		recordID: vo.ReconstructRecordID(recordIDStr),
		name:     vo.ReconstructItemName(nameStr),
		calories: vo.ReconstructCalories(caloriesVal),
	}
}

// ID はRecordItemIDを返す
func (ri *RecordItem) ID() vo.RecordItemID {
	return ri.id
}

// RecordID は親のRecordIDを返す
func (ri *RecordItem) RecordID() vo.RecordID {
	return ri.recordID
}

// Name は食品名を返す
func (ri *RecordItem) Name() vo.ItemName {
	return ri.name
}

// Calories はカロリーを返す
func (ri *RecordItem) Calories() vo.Calories {
	return ri.calories
}
