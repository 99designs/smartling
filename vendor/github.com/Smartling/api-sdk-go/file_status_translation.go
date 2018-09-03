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
