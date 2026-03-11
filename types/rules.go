package types

const (
	Images    = "images"
	Videos    = "videos"
	Audios    = "audios"
	Documents = "documents"
	Archives  = "archives"
	Other     = "other"
)

var Rules = map[string]string{
	// Images
	".jpg":  Images,
	".jpeg": Images,
	".png":  Images,
	".webp": Images,
	".tiff": Images,
	".tif":  Images,
	".psd":  Images,
	".raw":  Images,
	".avif": Images,
	".svg":  Images,
	".gif":  Images,

	// Videos
	".mp4":  Videos,
	".mkv":  Videos,
	".avi":  Videos,
	".webm": Videos,
	".mov":  Videos,
	".flv":  Videos,
	".wmv":  Videos,

	// Audios
	".mp3":  Audios,
	".aac":  Audios,
	".wav":  Audios,
	".flac": Audios,
	".aiff": Audios,
	".ogg":  Audios,

	// Documents
	".pdf":  Documents,
	".doc":  Documents,
	".xls":  Documents,
	".ppt":  Documents,
	".docx": Documents,
	".xlsx": Documents,
	".pptx": Documents,
	".csv":  Documents,
	".odt":  Documents,
	".odp":  Documents,
	".ods":  Documents,
	".txt":  Documents,

	// Archives
	".zip": Archives,
	".rar": Archives,
	".7z":  Archives,
	".tar": Archives,
}

func BuildCategorySet(rules map[string]string) map[string]bool {
	categories := map[string]bool{}

	for _, category := range rules {
		categories[category] = true
	}

	categories["other"] = true

	return categories
}
