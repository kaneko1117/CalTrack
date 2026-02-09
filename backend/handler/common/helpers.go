package common

// ExtractErrorMessages はエラーリストからメッセージ文字列を抽出する
func ExtractErrorMessages(errs []error) []string {
	details := make([]string, len(errs))
	for i, err := range errs {
		details[i] = err.Error()
	}
	return details
}
