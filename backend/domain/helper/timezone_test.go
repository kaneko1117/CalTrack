package helper

import (
	"testing"
	"time"
)

func TestJST(t *testing.T) {
	loc := JST()

	t.Run("nilでないこと", func(t *testing.T) {
		if loc == nil {
			t.Fatal("JST() returned nil")
		}
	})

	t.Run("タイムゾーン名がAsia/Tokyoであること", func(t *testing.T) {
		if loc.String() != "Asia/Tokyo" {
			t.Errorf("JST().String() = %q, want %q", loc.String(), "Asia/Tokyo")
		}
	})

	t.Run("UTCからのオフセットが+9時間であること", func(t *testing.T) {
		utcTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		jstTime := utcTime.In(loc)

		_, offset := jstTime.Zone()
		wantOffset := 9 * 60 * 60
		if offset != wantOffset {
			t.Errorf("JST offset = %d, want %d", offset, wantOffset)
		}
	})

	t.Run("同一インスタンスを返すこと", func(t *testing.T) {
		loc1 := JST()
		loc2 := JST()
		if loc1 != loc2 {
			t.Error("JST() should return the same instance")
		}
	})
}
