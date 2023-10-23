package content

import (
	"fmt"
	"github.com/bariiss/transfer.sh-cli/lib"
	c "github.com/bariiss/transfer.sh-cli/lib/config"
	"github.com/cheggaaa/pb/v3"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Path: lib/content/content.go

// PrepareContent prepares the content for uploading
func PrepareContent(filePath string) (fileName string, reader io.Reader, size int64, err error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return
	}

	if fileInfo.IsDir() {
		fileName = filepath.Base(filePath) + ".zip"
		zipPath := filepath.Join(os.TempDir(), fileName)
		err = ZipDirectory(filePath, zipPath)
		if err != nil {
			return
		}
		reader, err = os.Open(zipPath)
		if err != nil {
			return
		}
		defer func(name string) {
			removeErr := os.Remove(name)
			if removeErr != nil {
				fmt.Println("Error removing zip file:", removeErr)
			}
		}(zipPath) // ensure zip file is removed after uploading
		var zippedFileInfo os.FileInfo
		zippedFileInfo, err = os.Stat(zipPath) // Removed the inner declaration of err
		if err != nil {
			return
		}
		size = zippedFileInfo.Size()
		return
	}

	fileName = filepath.Base(filePath)
	reader, err = os.Open(filePath)
	size = fileInfo.Size()
	return
}

func UploadContent(filePath string, config *c.Config, maxDays string, maxDownloads string) (*http.Response, error) {
	client := &http.Client{}

	fileName, reader, size, err := PrepareContent(filePath)
	if err != nil {
		fmt.Println("Error preparing content:", err)
	}

	req, err := http.NewRequest("PUT", config.BaseURL+"/"+fileName, reader)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Max-Days", maxDays)
	req.Header.Add("Max-Downloads", maxDownloads)

	bar := pb.Full.Start64(size)
	barReader := bar.NewProxyReader(reader)
	req.Body = io.NopCloser(barReader)
	req.ContentLength = size

	req.SetBasicAuth(config.User, config.Pass)
	resp, err := client.Do(req)
	bar.Finish()

	if err != nil {
		return nil, fmt.Errorf("error uploading: %w", err)
	}

	err = lib.PrintResponse(resp, size, config, fileName)
	if err != nil {
		fmt.Println(err)
	}

	return resp, nil
}
