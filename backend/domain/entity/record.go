package entity

import (
	"time"

	domainErrors "caltrack/domain/errors"
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
	itemInputs []RecordItemInput,
) (*Record, []error) {
	var errs []error

	if len(itemInputs) == 0 {
		errs = append(errs, domainErrors.ErrRecordItemsRequired)
		return nil, errs
	}

	eatenAt, err := vo.NewEatenAt(eatenAtTime)
	errs = appendIfErr(errs, err)

	recordID := vo.NewRecordID()

	var items []RecordItem
	for _, input := range itemInputs {
		item, itemErrs := newRecordItem(recordID, input)
		if len(itemErrs) > 0 {
			errs = append(errs, itemErrs...)
		} else {
			items = append(items, *item)
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return &Record{
		id:        recordID,
		userID:    userID,
		eatenAt:   eatenAt,
		items:     items,
		createdAt: time.Now(),
	}, nil
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
