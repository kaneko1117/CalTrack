package vo

import "time"

// SetNowFunc はテスト用に現在時刻関数を差し替える
func SetNowFunc(f func() time.Time) {
	nowFunc = f
}

// ResetNowFunc は現在時刻関数をデフォルトに戻す
func ResetNowFunc() {
	nowFunc = time.Now
}
