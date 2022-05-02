package meme

import (
	"fmt"
	"io"

	exif "github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	"github.com/northwood-labs/golang-utils/exiterrorf"
)

func WriteImageDescription(r io.Reader, wr io.Writer, words string) error {
	// so := &exif.ScanOptions{}

	rawExif, err := exif.SearchAndExtractExifWithReader(r)
	if err != nil {
		return exiterrorf.Errorf(err, "failed to extract EXIF data from reader")
	}

	im, err := exifcommon.NewIfdMappingWithStandard()
	if err != nil {
		return exiterrorf.Errorf(err, "failed to identify EXIF mapping")
	}

	ti := exif.NewTagIndex()
	_, index, err := exif.Collect(im, ti, rawExif)
	ib := exif.NewIfdBuilderFromExistingChain(index.RootIfd)

	// Read the IFD whose tag we want to change.

	// Standard:
	// - "IFD0"
	// - "IFD0/Exif0"
	// - "IFD0/Exif0/Iop0"
	// - "IFD0/GPSInfo0"
	//
	// If the numeric indices are not included, (0) is the default. Note that
	// this isn't strictly necessary in our case since IFD0 is the first IFD anyway, but we're putting it here to show usage.
	ifdPath := "IFD0"

	childIb, err := exif.GetOrCreateIbFromRootIb(ib, ifdPath)
	if err != nil {
		return exiterrorf.Errorf(err, "failed to create IB from root")
	}

	// There are a few functions that allow you to surgically change the tags in an
	// IFD, but we're just gonna overwrite a tag that has an ASCII value.

	tagName := "ImageDescription"

	err = childIb.SetStandardWithName(tagName, words)

	// Encode the in-memory representation back down to bytes.

	ibe := exif.NewIfdByteEncoder()

	updatedRawExif, err := ibe.EncodeToExif(ib)

	// Reparse the EXIF to confirm that our value is there.

	_, index, err = exif.Collect(im, ti, updatedRawExif)

	// This isn't strictly necessary for the same reason as above, but it's here
	// for documentation.
	childIfd, err := exif.FindIfdFromRootIfd(index.RootIfd, ifdPath)

	results, err := childIfd.FindTagWithName(tagName)

	for _, ite := range results {
		valueRaw, err := ite.Value()

		stringValue := valueRaw.(string)
		fmt.Println(stringValue)
	}

}
