package smartling

type FileStatusTranslation struct {
	LocaleID string

	AuthorizedStringCount int
	AuthorizedWordCount   int
	CompletedStringCount  int
	CompletedWordCount    int
	ExcludedStringCount   int
	ExcludedWordCount     int
}

func (fst FileStatusTranslation) AwaitingAuthorizationStringCount(totalStringCount int) int {
	return totalStringCount - fst.AuthorizedStringCount - fst.ExcludedStringCount - fst.CompletedStringCount
}
