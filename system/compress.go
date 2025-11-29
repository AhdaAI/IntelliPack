package system

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"

	"github.com/AhdaAI/IntelliPack/system/file"
	"github.com/klauspost/compress/zstd"
)

func CompressFolder(folderPath, outputFile string) error {
	out, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer out.Close()

	encoder, err := zstd.NewWriter(
		out,
		zstd.WithEncoderLevel(zstd.SpeedBetterCompression),
	) // ZStandard Writer
	if err != nil {
		return err
	}
	defer encoder.Close()

	tw := tar.NewWriter(encoder) // Tar Writer
	defer tw.Close()

	return filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, _ := filepath.Rel(folderPath, path)

		if info.IsDir() {
			hdr := &tar.Header{
				Name:     relPath + "/",
				Mode:     0755,
				Typeflag: tar.TypeDir,
			}
			return tw.WriteHeader(hdr)
		}

		hdr := &tar.Header{
			Name: relPath,
			Mode: 0644,
			Size: info.Size(),
		}

		if file.IsCompressed(path) {
			hdr.PAXRecords = map[string]string{
				"compressed": "false",
			}

			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}

			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			_, err = io.Copy(tw, f)
			return err
		}

		hdr.PAXRecords = map[string]string{
			"compressed": "true",
		}

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		// Compress file data manually
		// (Zstd compresses the whole tar stream anyway,
		// but you can pipe file-by-file too)
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(tw, f)
		return err
	})
}
