package smartling

// FileStatus describes file translation status obtained by GetFileStatus
// method.
type FileStatus struct {
	File

	TotalStringCount int
	TotalWordCount   int
	TotalCount       int

	Items []FileStatusTranslation
}
