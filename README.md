# transfersh

`transfersh` is a Command-Line Interface (CLI) tool to easily upload files or directories to [transfer.sh](https://github.com/dutchcoders/transfer.sh). The app is specifically designed around the transfer.sh infrastructure.

## Features
- Upload files or directories directly from the command line.
- Automatically compresses directories into a `.zip` file for uploading.
- Provides a direct download link for the uploaded content.
- Returns a `curl` command for content deletion.
- Simple configuration stored in `~/.config/transfersh-cli/.config`.

## Installation

You need to have [Go](https://golang.org/) installed. Then you can get the tool via:

```bash
go install github.com/bariiss/transfer.sh-cli/cli/transfersh-cli@latest
```

## Usage
```bash
transfersh-cli [file|directory]
```

- file: The path to the file you want to upload.
- directory: The path to the directory you want to compress and upload.

The tool will return the URL for your uploaded content. If the content is a directory, it will be compressed as a .zip file and then uploaded.

## Configuration

On the first run, the CLI will prompt you for the transfer.sh base URL and your optional user and pass (for basic authentication if enabled on your transfer.sh instance). The configuration will be saved at ~/.config/transfersh/.config.

## License

MIT

## Author Information

This tool was created in 2023 by [bariiss](https://github.com/bariiss)