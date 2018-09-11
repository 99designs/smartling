package smartling

// FileStatusExtended describes file translation status obtained by GetFileStatusExtended
// method.
type FileStatusExtended struct {
	File

	AuthorizedStringCount int
	AuthorizedWordCount   int
	CompletedStringCount  int
	CompletedWordCount    int
	ExcludedStringCount   int
	ExcludedWordCount     int
	TotalStringCount      int
	TotalWordCount        int
}

func (fs FileStatusExtended) AwaitingAuthorizationStringCount() int {
	return fs.TotalStringCount - fs.AuthorizedStringCount - fs.ExcludedStringCount - fs.CompletedStringCount
}

func (fs FileStatusExtended) InProgressStringCount() int {
	return fs.NotCompletedStringCount() - fs.AwaitingAuthorizationStringCount()
}

func (fs FileStatusExtended) NotCompletedStringCount() int {
	return fs.TotalStringCount - fs.CompletedStringCount
}
