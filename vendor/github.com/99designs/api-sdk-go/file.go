package smartling

// FileStatus represents current file status in the Smartling system.
type File struct {
	// FileURI is a unique path to file in Smartling system.
	FileURI string

	// FileType is a file type identifier.
	FileType FileType

	// LastUploaded refers to time when file was uploaded.
	LastUploaded UTC

	// HasInstructions specifies does files have instructions or not.
	HasInstructions bool
}
