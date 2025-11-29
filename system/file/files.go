package file

import (
	"log"
	"path/filepath"
	"strings"
)

var compressedExt = map[string]bool{
	".zip": true, ".rar": true, ".7z": true,
	".jpg": true, ".jpeg": true, ".png": true,
	".gif": true, ".mp4": true, ".mkv": true,
	".mp3": true, ".ogg": true, ".pdf": true,
	".zst": true, ".gz": true, ".bz2": true,
	".xz": true, ".lz4": true,
}

func IsCompressed(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	compress := compressedExt[ext]
	if compress {
		log.Printf("Skipping %s ...\n", path)
	} else {
		log.Printf("Compressing %s ...\n", path)
	}
	return compress
}
