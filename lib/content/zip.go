package content

import (
	"archive/zip"
	"os"
	"path/filepath"
)

// Path: lib/content/zip.go

// ZipDirectory zips the directory and returns a reader
func ZipDirectory(directory, zipPath string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || (info.Mode()&os.ModeSocket) != 0 {
			return err
		}

		relPath, _ := filepath.Rel(directory, path)
		zf, _ := zipWriter.Create(relPath)
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		_, err = zf.Write(content)
		return err
	})
}
