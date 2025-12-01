package file

import (
	"path/filepath"
	"strings"
)

var compressedExt = map[string]bool{
	".zip": true, ".rar": true, ".7z": true,
	".jpg": true, ".jpeg": true, ".png": true,
	".gif": true, ".mp4": true, ".mkv": true,
	".mp3": true, ".ogg": true, ".pdf": true,
	".zst": true, ".gz": true, ".bz2": true,
	".xz": true, ".lz4": true, ".pak": true,
	".assets": true, ".sharedAssets": true, ".forge": true,
	".utoc": true, ".ucas": true,
}

func IsCompressed(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	compress := compressedExt[ext]
	return compress
}
