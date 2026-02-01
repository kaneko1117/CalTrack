package vo_test

import (
	"testing"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
)

func TestNewItemName(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantItemName string
		wantErr      error
	}{
		// 正常系
		{"正常な食品名", "りんご", "りんご", nil},
		{"英語の食品名", "Apple", "Apple", nil},
		{"スペース付き食品名", "グリーン サラダ", "グリーン サラダ", nil},
		{"長い食品名", "特製チーズハンバーグステーキ定食（ライス大盛り）", "特製チーズハンバーグステーキ定食（ライス大盛り）", nil},
		// 異常系
		{"空文字の場合はエラー", "", "", domainErrors.ErrItemNameRequired},
		// 境界値
		{"最小長1文字は有効", "a", "a", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vo.NewItemName(tt.input)

			if err != tt.wantErr {
				t.Errorf("NewItemName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if err == nil && got.String() != tt.wantItemName {
				t.Errorf("NewItemName(%q).String() = %v, want %v", tt.input, got.String(), tt.wantItemName)
			}
		})
	}
}

func TestReconstructItemName(t *testing.T) {
	t.Run("DBから復元した食品名", func(t *testing.T) {
		value := "りんご"
		got := vo.ReconstructItemName(value)

		if got.String() != value {
			t.Errorf("ReconstructItemName(%q).String() = %v, want %v", value, got.String(), value)
		}
	})

	t.Run("空文字でもバリデーションエラーにならない", func(t *testing.T) {
		// ReconstructはDBからの復元用なのでバリデーションをスキップ
		got := vo.ReconstructItemName("")

		if got.String() != "" {
			t.Errorf("ReconstructItemName(\"\").String() = %v, want \"\"", got.String())
		}
	})
}
