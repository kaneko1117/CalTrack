package entity

import (
	"time"

	"caltrack/domain/vo"
)

// Record はカロリー記録を表すエンティティ
type Record struct {
	id        vo.RecordID
	userID    vo.UserID
	eatenAt   vo.EatenAt
	items     []RecordItem
	createdAt time.Time
}

// NewRecord は新しいRecordを生成する
func NewRecord(
	userID vo.UserID,
	eatenAtTime time.Time,
) (*Record, error) {
	eatenAt, err := vo.NewEatenAt(eatenAtTime)
	if err != nil {
		return nil, err
	}

	return &Record{
		id:        vo.NewRecordID(),
		userID:    userID,
		eatenAt:   eatenAt,
		items:     []RecordItem{},
		createdAt: time.Now(),
	}, nil
}

// AddItem はRecordにRecordItemを追加する
func (r *Record) AddItem(nameStr string, caloriesVal int) error {
	item, errs := NewRecordItem(r.id, nameStr, caloriesVal)
	if len(errs) > 0 {
		return errs[0]
	}
	r.items = append(r.items, *item)
	return nil
}

// ReconstructRecord はDBからRecordを復元する
func ReconstructRecord(
	idStr string,
	userIDStr string,
	eatenAtTime time.Time,
	createdAt time.Time,
	items []RecordItem,
) *Record {
	return &Record{
		id:        vo.ReconstructRecordID(idStr),
		userID:    vo.ReconstructUserID(userIDStr),
		eatenAt:   vo.ReconstructEatenAt(eatenAtTime),
		items:     items,
		createdAt: createdAt,
	}
}

// ID はRecordIDを返す
func (r *Record) ID() vo.RecordID {
	return r.id
}

// UserID はユーザーIDを返す
func (r *Record) UserID() vo.UserID {
	return r.userID
}

// EatenAt は食事日時を返す
func (r *Record) EatenAt() vo.EatenAt {
	return r.eatenAt
}

// Items は記録明細リストを返す
func (r *Record) Items() []RecordItem {
	return r.items
}

// CreatedAt は作成日時を返す
func (r *Record) CreatedAt() time.Time {
	return r.createdAt
}

// TotalCalories は記録の合計カロリーを返す
func (r *Record) TotalCalories() int {
	total := 0
	for _, item := range r.items {
		total += item.Calories().Value()
	}
	return total
}
