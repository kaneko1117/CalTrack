package entity

// appendIfErr はerrがnilでない場合にスライスに追加するヘルパー関数
func appendIfErr(errs []error, err error) []error {
	if err != nil {
		return append(errs, err)
	}
	return errs
}
