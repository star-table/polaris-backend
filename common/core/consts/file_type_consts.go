package consts

//资源类型
const (
	FileTypeOthers = 0
	FileTypeDoc    = 1
	FileTypePhoto  = 2
	FileTypeVideo  = 3
	FileTypeMusic  = 4
)

var FileTypes = map[string]int{
	"JPG":       FileTypePhoto,
	"JPEG":      FileTypePhoto,
	"JFIF":      FileTypePhoto,
	"JPEG 2000": FileTypePhoto,
	"GIF":       FileTypePhoto,
	"BMP":       FileTypePhoto,
	"PNG":       FileTypePhoto,
	"WEBP":      FileTypePhoto,
	"HDR":       FileTypePhoto,
	"HEIF":      FileTypePhoto,
	"BDP":       FileTypePhoto,
	"SVG":       FileTypePhoto,

	"PPT":  FileTypeDoc,
	"XLS":  FileTypeDoc,
	"XLSX": FileTypeDoc,
	"CSV":  FileTypeDoc,
	"DOCX": FileTypeDoc,
	"DOC":  FileTypeDoc,
	"PSD":  FileTypeDoc,
	"TXT":  FileTypeDoc,
	"PDF":  FileTypeDoc,
	"ZIP":  FileTypeDoc,
	"RAR":  FileTypeDoc,

	"MP4": FileTypeVideo,
	"AVI": FileTypeVideo,
	"MPG": FileTypeVideo,
	"RMV": FileTypeVideo,
	"MOV": FileTypeVideo,

	"MP3": FileTypeMusic,
}
