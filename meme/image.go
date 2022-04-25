package meme

import (
	"bytes"
	"image"
	_ "image/gif" // lint:allow_blank_imports
	jpeg "image/jpeg"
	_ "image/png" // lint:allow_blank_imports
	"io"

	_ "github.com/jdeng/goheif" // lint:allow_blank_imports
	"github.com/northwood-labs/golang-utils/exiterrorf"
	_ "golang.org/x/image/bmp"  // lint:allow_blank_imports
	_ "golang.org/x/image/tiff" // lint:allow_blank_imports
	_ "golang.org/x/image/webp" // lint:allow_blank_imports
)

// Format is a string type represents the image format.
type Format string

const (
	// FormatUnknown represents an unknown format.
	FormatUnknown Format = ""
	// FormatBMP represents the BMP format.
	FormatBMP Format = "bmp"
	// FormatGIF represents the GIF format.
	FormatGIF Format = "gif"
	// FormatHEIC represents the HEIC format.
	FormatHEIC Format = "heic"
	// FormatJPEG represents the JPEG format.
	FormatJPEG Format = "jpeg"
	// FormatPNG represents the PNG format.
	FormatPNG Format = "png"
	// FormatTIFF represents the TIFF format.
	FormatTIFF Format = "tiff"
	// FormatWEBP represents the WebP format.
	FormatWEBP Format = "webp"

	// DefaultJPEGQuality is the default JPEG compression quality.
	DefaultJPEGQuality = 80
)

// ReadImage reads the image into a PNG bytestream.
func ReadImage(r io.Reader, jpegQuality int) (bytes.Buffer, image.Image, string, error) {
	var img image.Image
	var buf bytes.Buffer

	img, format, err := image.Decode(r)
	if err != nil {
		return buf, img, "", exiterrorf.Errorf(err, "decoding error")
	}

	err = jpeg.Encode(&buf, img, &jpeg.Options{
		Quality: jpegQuality,
	})
	if err != nil {
		return buf, img, "", exiterrorf.Errorf(err, "png encoding error")
	}

	return buf, img, format, nil
}
