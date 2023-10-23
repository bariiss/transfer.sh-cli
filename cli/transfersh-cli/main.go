package main

import (
	"fmt"
	c "github.com/bariiss/transfer.sh-cli/lib/config"
	ct "github.com/bariiss/transfer.sh-cli/lib/content"
	"github.com/spf13/cobra"
	"os"
)

// Path: main.go

var version = "0.1.7"
var maxDays string
var maxDownloads string

var rootCmd = &cobra.Command{
	Use:   "transfersh-cli [file|directory]",
	Short: "transfersh-cli files or directories.",
	Long: `transfer.sh-cli is a CLI tool for uploading files or directories.
Given a file or directory path, it will upload the content and 
provide a URL for download. If a directory path is provided,
it will be compressed as a .zip file and then uploaded.`,
	Version: version,
	Args:    cobra.ExactArgs(1),
	Run:     executeTransfer,
}

func init() {
	rootCmd.Flags().StringVar(&maxDays, "max-days", "", "Maximum number of days before the file is deleted")
	rootCmd.Flags().StringVar(&maxDownloads, "max-downloads", "", "Maximum number of downloads before the file is deleted")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func executeTransfer(cmd *cobra.Command, args []string) {
	loadConfig, err := c.LoadConfig()
	if err != nil {
		fmt.Println("Error loading loadConfig:", err)
		return
	}

	filePath := args[0]                                  // file or directory path
	if _, err := os.Stat(filePath); os.IsNotExist(err) { // check if file or directory exists
		fmt.Printf("%s: No such file or directory\n", filePath)
		return
	}

	_, err = ct.UploadContent(filePath, loadConfig, maxDays, maxDownloads)
	if err != nil {
		fmt.Println(err)
		return
	}
}
