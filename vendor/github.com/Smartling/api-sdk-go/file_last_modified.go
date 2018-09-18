package smartling

type FileLastModifiedLocales struct {
	Items []FileLastModified
}

type FileLastModified struct {
	LocaleID     string
	LastModified UTC
}
