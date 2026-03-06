package types

var Rules = map[string]string{
	// Images
	".jpg":  "images",
	".jpeg": "images",
	".png":  "images",
	".webp": "images",
	".tiff": "images",
	".tif":  "images",
	".psd":  "images",
	".raw":  "images",
	".avif": "images",
	".svg":  "images",
	".gif":  "images",

	// Videos
	".mp4":  "videos",
	".mkv":  "videos",
	".avi":  "videos",
	".webm": "videos",
	".mov":  "videos",
	".flv":  "videos",
	".wmv":  "videos",

	// Documents
	".pdf":  "documents",
	".doc":  "documents",
	".xls":  "documents",
	".ppt":  "documents",
	".docx": "documents",
	".xlsx": "documents",
	".pptx": "documents",
	".csv":  "documents",
	".odt":  "documents",
	".odp":  "documents",
	".ods":  "documents",
	".txt":  "documents",

	// Archives
	".zip": "archives",
	".rar": "archives",
	".7z":  "archives",
	".tar": "archives",
}
