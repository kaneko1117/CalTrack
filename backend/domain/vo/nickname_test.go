package vo_test

import (
	"strings"
	"testing"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
)

func TestNewNickname(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantNickname string
		wantErr      error
	}{
		// 正常系
		{"正常なニックネーム", "John", "John", nil},
		{"スペース付きニックネーム", "John Doe", "John Doe", nil},
		// 異常系
		{"空文字の場合はエラー", "", "", domainErrors.ErrNicknameRequired},
		{"51文字以上はエラー", strings.Repeat("a", 51), "", domainErrors.ErrNicknameTooLong},
		// 境界値
		{"最小長1文字は有効", "A", "A", nil},
		{"最大長50文字は有効", strings.Repeat("a", 50), strings.Repeat("a", 50), nil},
		{"51文字は無効", strings.Repeat("a", 51), "", domainErrors.ErrNicknameTooLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := vo.NewNickname(tt.input)

			if err != tt.wantErr {
				t.Errorf("NewNickname(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if err == nil && got.String() != tt.wantNickname {
				t.Errorf("NewNickname(%q).String() = %v, want %v", tt.input, got.String(), tt.wantNickname)
			}
		})
	}
}
