package lib

import (
	"fmt"
	"github.com/atotto/clipboard"
	c "github.com/bariiss/transfer.sh-cli/lib/config"
	"github.com/fatih/color"
	"io"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"
)

// Path: lib/response.go

func PrintResponse(resp *http.Response, size int64, config *c.Config, fileName string) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body:", err)
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to upload: %s", body)
	}

	hiBlue := color.New(color.FgHiBlue).SprintFunc()

	// Copy to clipboard if possible, don't fail if it's not possible
	if err := clipboard.WriteAll(string(body)); err != nil {
		fmt.Println(hiBlue("Failed to copy to clipboard:"), hiBlue(err))
	}

	blue := color.New(color.FgBlue).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)

	sizeStr := fmtSize(size)
	_, err = fmt.Fprintf(w, "\n%s (%s):\t", blue(fileName), yellow(sizeStr))
	if err != nil {
		return fmt.Errorf("error writing to stdout: %w", err)
	}
	_, err = fmt.Fprintf(w, "\n%s\t%s\t", green("File URL:"), body)
	if err != nil {
		return fmt.Errorf("error writing to stdout: %w", err)
	}

	directDownloadLink := strings.Replace(string(body), config.BaseURL+"/", config.BaseURL+"/get/", 1)
	_, err = fmt.Fprintf(w, "\n%s\t%s\t", green("Direct Download Link:"), directDownloadLink)
	if err != nil {
		return fmt.Errorf("error writing to stdout: %w", err)
	}

	for name, values := range resp.Header {
		if strings.ToLower(name) == "x-url-delete" {
			_, err = fmt.Fprintf(w, "\n%s\tcurl -X DELETE \"%s\"\t", red("To delete use:"), values[0])
			if err != nil {
				return fmt.Errorf("error writing to stdout: %w", err)
			}
		}
	}

	_, err = fmt.Fprintln(w)
	if err != nil {
		return fmt.Errorf("error writing to stdout: %w", err)
	}
	err = w.Flush()
	if err != nil {
		return fmt.Errorf("error flushing stdout: %w", err)
	}

	return nil
}

func fmtSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%d B", size)
	}
	if size < 1024*1024 {
		return fmt.Sprintf("%d KB", size/1024)
	}
	if 1024*1024 <= size && size < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(size)/float64(1024*1024))
	}
	return fmt.Sprintf("%.2f GB", float64(size)/float64(1024*1024*1024))
}
