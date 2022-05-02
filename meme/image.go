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

const (
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
