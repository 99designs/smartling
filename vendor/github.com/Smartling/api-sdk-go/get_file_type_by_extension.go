package smartling

import "strings"

var (
	extensions = map[string]FileType{
		"yml":         FileTypeYAML,
		"yaml":        FileTypeYAML,
		"html":        FileTypeHTML,
		"htm":         FileTypeHTML,
		"xlf":         FileTypeXLIFF,
		"xliff":       FileTypeXLIFF,
		"json":        FileTypeJSON,
		"docx":        FileTypeDOCX,
		"pptx":        FileTypePPTX,
		"xlsx":        FileTypeXLSX,
		"txt":         FileTypePlaintext,
		"ts":          FileTypeQt,
		"idml":        FileTypeIDML,
		"resx":        FileTypeResx,
		"resw":        FileTypeResx,
		"csv":         FileTypeCSV,
		"stringsdict": FileTypeStringsdict,
		"strings":     FileTypeIOS,
		"po":          FileTypeGettext,
		"pot":         FileTypeGettext,
		"xml":         FileTypeXML,
		"properties":  FileTypeJavaProperties,
		"":            FileTypeUnknown,
	}
)

func GetFileTypeByExtension(ext string) FileType {
	return extensions[strings.TrimPrefix(ext, ".")]
}
