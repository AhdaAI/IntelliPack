package system

import (
	"archive/tar"
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

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
		zstd.WithEncoderConcurrency(runtime.NumCPU()/2),
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

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		// ===========================================================
		// CASE 1: FILE IS ALREADY COMPRESSED → Write RAW (unmodified)
		// ===========================================================
		if file.IsCompressed(path) {
			log.Printf("Skipping %s ...\n", relPath)
			hdr := &tar.Header{
				Name: relPath,
				Mode: 0644,
				Size: info.Size(),
				PAXRecords: map[string]string{
					"compressed": "false",
				},
			}

			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}

			_, err = io.Copy(tw, f)
			return err
		}

		// ===========================================================
		// CASE 2: FILE IS NOT COMPRESSED → Compress it now
		// ===========================================================
		log.Printf("Compressing %s ...\n", relPath)
		var buf bytes.Buffer

		z, err := zstd.NewWriter(&buf)
		if err != nil {
			return err
		}

		_, err = io.Copy(z, f)
		if err != nil {
			z.Close()
			return err
		}

		z.Close()
		compressedBytes := buf.Bytes()

		hdr := &tar.Header{
			Name:   relPath,
			Mode:   0644,
			Size:   int64(len(compressedBytes)),
			Format: tar.FormatPAX,
			PAXRecords: map[string]string{
				"compressed": "true",
			},
		}

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		_, err = tw.Write(compressedBytes)
		return err
	})
}
