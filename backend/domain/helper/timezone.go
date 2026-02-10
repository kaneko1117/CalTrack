package helper

import "time"

// jst は日本標準時（UTC+9）のロケーション
var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

// JST は日本標準時の *time.Location を返す
func JST() *time.Location {
	return jst
}
