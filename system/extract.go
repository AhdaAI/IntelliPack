package system

import (
	"archive/tar"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
)

func ExtractArchive(path, output string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer f.Close()

	zr, err := zstd.NewReader(f)
	if err != nil {
		return fmt.Errorf("failed to start reader: %w", err)
	}
	defer zr.Close()

	tr := tar.NewReader(zr)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read header: %w", err)
		}

		targetPath := filepath.Join(output, hdr.Name)

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, os.FileMode(hdr.Mode)); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}

		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}
			outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.FileMode(hdr.Mode))
			if err != nil {
				return fmt.Errorf("failed to create file: %w", err)
			}

			if hdr.PAXRecords != nil {
				if val, ok := hdr.PAXRecords["compressed"]; ok && val == "true" {
					log.Printf("Extracting %s ...", filepath.Base(hdr.Name))

					inner, err := zstd.NewReader(tr)
					if err != nil {
						outFile.Close()
						return fmt.Errorf("failed to inner zstd: %w", err)
					}

					_, err = io.Copy(outFile, inner)
					inner.Close()

					if err != nil {
						outFile.Close()
						return fmt.Errorf("failed extracting compressed file: %w", err)
					}
				} else {
					log.Printf("Copying %s ...", filepath.Base(hdr.Name))
					if _, err := io.Copy(outFile, tr); err != nil {
						outFile.Close()
						return fmt.Errorf("failed writing raw file: %w", err)
					}
				}
			}
			outFile.Close()

		default:
			fmt.Printf("Skipping unsupported tar entry: %s\n", hdr.Name)
		}
	}

	return nil
}
