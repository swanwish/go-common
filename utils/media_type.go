package utils

import "path/filepath"

var (
	MediaTypeHtml     = "application/xhtml+xml"
	MediaTypeImage    = "image/jpeg"
	MediaTypeMarkdown = "text/markdown"
	MediaTypeOther    = "other"
)

func GetMediaType(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".html", ".htm", ".xhtml":
		return MediaTypeHtml
	case ".jpg", ".png", ".jpeg":
		return MediaTypeImage
	case ".md", ".markdown":
		return MediaTypeMarkdown
	default:
		return MediaTypeOther
	}
}
