# qrclip

`qrclip` is a small command-line tool for decoding QR codes from an image in the clipboard or from an image file.

It is useful when you see a QR code on your desktop and want to extract the URL or text without using your phone.

## Features
- Read QR codes from a clipboard image
- Read QR codes from an image file
- Print decoded text to the terminal
- Optionally copy decoded text back to the clipboard
- Works well with screenshot tools such as Windows Snipping Tool

## Installation
You can download the latest release from the [releases page](https://github.com/sean1832/qrclip/releases/latest). Make sure to choose the correct binary for your operating system.

| Platform | Status      | Notes                                                                                                                           |
| -------- | ----------- | ------------------------------------------------------------------------------------------------------------------------------- |
| Windows  | Supported   | Recommended platform. Works with Win + Shift + S.                                                                               |
| macOS    | Supported   | Works with screenshot-to-clipboard workflows. The release binary may be unsigned.                                               |
| Linux    | Best effort | Clipboard mode requires X11 or XWayland. Wayland-only sessions may not support clipboard image input. Use `--file` as fallback. |

## Usage
```
Usage:
  qrclip [options]

Options:
  -c, --copy        Copy decoded text to clipboard
  -f, --file PATH   Read QR code from image file instead of clipboard
  -h, --help        Show help
```
Running `qrclip` without any options will read the QR code from the clipboard and print the decoded text to the console.

### Examples:
```bash
./qrclip                # Read QR code from clipboard and print to console
./qrclip -c             # Read QR code from clipboard and copy decoded text to clipboard
./qrclip -f qr.png      # Read QR code from image file and print to console
./qrclip -f qr.png -c   # Read QR code from image file and copy decoded text to clipboard
```

## Building from source
Make sure you have Go installed on your system. Then, clone the repository and build the project:

```bash
git clone https://github.com/sean1832/qrclip.git
cd qrclip
go build -o qrclip.exe main.go
```

> Build with options for smaller binary size:
> ```bash
> go build -ldflags="-s -w" -o qrclip.exe main.go
> ```

> To include version information in the binary, use:
> ```bash
> go build -ldflags="-s -w -X 'main.appVersion=1.0.0'" -o qrclip.exe main.go
> ```

### Linux build dependencies

On Linux, clipboard support may require X11 development libraries.

For Ubuntu/Debian:

```bash
sudo apt-get update
sudo apt-get install -y libx11-dev
```

## License
This project is licensed under Apache License 2.0. See the [LICENSE](LICENSE) file for details.