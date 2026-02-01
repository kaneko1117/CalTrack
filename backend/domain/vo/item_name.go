package vo

import (
	domainErrors "caltrack/domain/errors"
)

// ItemName は食品名を表すValue Object
type ItemName struct {
	value string
}

// NewItemName は新しいItemNameを生成する
// 空文字の場合はエラーを返す
func NewItemName(value string) (ItemName, error) {
	if value == "" {
		return ItemName{}, domainErrors.ErrItemNameRequired
	}
	return ItemName{value: value}, nil
}

// ReconstructItemName はDBから復元する際に使用する
// バリデーションをスキップする
func ReconstructItemName(value string) ItemName {
	return ItemName{value: value}
}

// String は食品名を文字列として返す
func (i ItemName) String() string {
	return i.value
}
