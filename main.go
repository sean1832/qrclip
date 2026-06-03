package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"os"
	"strings"

	_ "image/jpeg"
	_ "image/png"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"golang.design/x/clipboard"
)

const appName = "qrclip"

type Options struct {
	Copy     bool
	FilePath string
}

func main() {
	opts := parseArgs()

	if opts.FilePath == "" || opts.Copy {
		if err := clipboard.Init(); err != nil {
			fail("failed to initialize clipboard", err)
		}
	}

	img, err := loadImage(opts)
	if err != nil {
		fail(err.Error(), nil)
	}

	text, err := decodeQRCode(img)
	if err != nil {
		fail("no QR code found in the image", err)
	}

	text = strings.TrimSpace(text)
	if text == "" {
		fail("decoded QR code is empty", nil)
	}
	fmt.Println(text)

	if opts.Copy {
		if err := writeTextToClipboard(text); err != nil {
			fail("failed to copy text to clipboard", err)
		}
	}
}

func parseArgs() Options {
	var opts Options
	fs := flag.NewFlagSet(appName, flag.ExitOnError)

	fs.BoolVar(&opts.Copy, "c", false, "copy decoded text to clipboard")
	fs.BoolVar(&opts.Copy, "copy", false, "copy decoded text to clipboard")

	fs.StringVar(&opts.FilePath, "f", "", "read QR code from image file")
	fs.StringVar(&opts.FilePath, "file", "", "read QR code from image file")

	fs.Usage = func() {
		fmt.Fprintf(fs.Output(), "Usage:\n")
		fmt.Fprintf(fs.Output(), "  %s [options]\n", appName)
		fmt.Fprintf(fs.Output(), "\n")
		fmt.Fprintf(fs.Output(), "Default:\n")
		fmt.Fprintf(fs.Output(), "  Reads an image from the clipboard and prints the decoded QR text.\n")
		fmt.Fprintf(fs.Output(), "\n")
		fmt.Fprintf(fs.Output(), "Options:\n")
		fmt.Fprintf(fs.Output(), "  -c, --copy        Copy decoded text to clipboard\n")
		fmt.Fprintf(fs.Output(), "  -f, --file PATH   Read QR code from image file instead of clipboard\n")
		fmt.Fprintf(fs.Output(), "  -h, --help        Show help\n")
	}

	if err := fs.Parse(os.Args[1:]); err != nil {
		fail("failed to parse arguments", err)
	}

	if fs.NArg() > 0 {
		fail("unexpected positional argument: "+fs.Arg(0), nil)
	}

	return opts
}

// Load images from either a file or the clipboard, depending on the provided options.
func loadImage(opts Options) (image.Image, error) {
	if opts.FilePath != "" {
		return loadImageFromFile(opts.FilePath)
	}

	return loadImageFromClipboard()
}

func loadImageFromFile(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image file: %w", err)
	}

	return img, nil
}

func loadImageFromClipboard() (image.Image, error) {
	imgBytes := clipboard.Read(clipboard.FmtImage)
	if len(imgBytes) == 0 {
		return nil, fmt.Errorf("no image found in clipboard")
	}

	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image from clipboard: %w", err)
	}

	return img, nil
}

// Write the given text to the clipboard.
func writeTextToClipboard(text string) error {
	clipboard.Write(clipboard.FmtText, []byte(text))
	return nil
}

// Decode the QR code from the given image and return the decoded text.
func decodeQRCode(img image.Image) (string, error) {
	bitmap, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}

	reader := qrcode.NewQRCodeReader()
	result, err := reader.Decode(bitmap, nil)
	if err != nil {
		return "", err
	}

	return result.GetText(), nil
}

func fail(message string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
	} else {
		fmt.Fprintln(os.Stderr, message)
	}

	os.Exit(1)
}
