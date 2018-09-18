package smartling

// FileType represents file type format used in Smartling API.
type FileType string

// Android and next are types that are supported by Smartling API.
const (
	FileTypeUnknown        FileType = ""
	FileTypeAndroid                 = "android"
	FileTypeIOS                     = "ios"
	FileTypeGettext                 = "gettext"
	FileTypeHTML                    = "html"
	FileTypeJavaProperties          = "javaProperties"
	FileTypeYAML                    = "yaml"
	FileTypeXLIFF                   = "xliff"
	FileTypeXML                     = "xml"
	FileTypeJSON                    = "json"
	FileTypeDOCX                    = "docx"
	FileTypePPTX                    = "pptx"
	FileTypeXLSX                    = "xlsx"
	FileTypeIDML                    = "idml"
	FileTypeQt                      = "qt"
	FileTypeResx                    = "resx"
	FileTypePlaintext               = "plaintext"
	FileTypeCSV                     = "csv"
	FileTypeStringsdict             = "stringsdict"
)
