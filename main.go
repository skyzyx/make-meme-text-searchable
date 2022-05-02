package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/caarlos0/spin"
	"github.com/dustin/go-humanize"
	"github.com/northwood-labs/awsutils"
	"github.com/northwood-labs/golang-utils/exiterrorf"
	"github.com/northwood-labs/golang-utils/log"
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"

	"github.com/skyzyx/make-meme-text-searchable/meme"
)

const defaultRetries = 5

var errFilePattern = errors.New("a file pattern is required")

func main() {
	ctx := context.Background()
	logger := log.GetStdTextLogger()

	var flagProfile string
	var flagRegion string
	var flagRetries int
	var flagVerbose bool
	var flagQuiet bool

	flag.StringVarP(&flagProfile, "profile", "p", "", "The AWS CLI profile with which to perform the request.")
	flag.StringVarP(&flagRegion, "region", "r", "", "The AWS region in which to perform the request.")
	flag.IntVarP(&flagRetries, "retries", "t", defaultRetries, "The max number of times to retry failed AWS requests.")
	flag.BoolVarP(&flagVerbose, "verbose", "v", false, "Enable verbose logging. Opposite of quiet mode.")
	flag.BoolVarP(&flagQuiet, "quiet", "q", false, "Enable silent mode. Opposite of verbose mode.")
	flag.Parse()

	if len(os.Args) == 1 {
		fmt.Fprintf(os.Stderr, "%s\n", errFilePattern)
		os.Exit(1)
	}

	config, err := awsutils.GetAWSConfig(ctx, flagRegion, flagProfile, flagRetries, flagVerbose)
	if err != nil {
		exiterrorf.ExitErrorf(err)
	}

	// Iterate over the Bash inputs
	for i := range os.Args[1:] {
		input := os.Args[1:][i]

		spinPrompt := fmt.Sprintf("%s %%s ", input)
		s := spin.New(spinPrompt)
		s.Set(spin.Box2)
		s.Start()

		r, err := os.Open(input) // lint:allow_possible_insecure
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", exiterrorf.Errorf(err, "failed to open file"))
			os.Exit(1)
		}

		// exif

		buf, _, _, err := meme.ReadImage(r, meme.DefaultJPEGQuality)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", exiterrorf.Errorf(err, "failed to convert image data as jpeg"))
			os.Exit(1)
		}

		if flagVerbose {
			logger.Info().Str("Size", humanize.Bytes(uint64(len(buf.Bytes())))).Msg("Rekognition limit for images is 5 MB.")
		}

		results, err := meme.DetectText(ctx, &config, buf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", exiterrorf.Errorf(err, "failed to detect the text in the image"))
			os.Exit(1)
		}

		words := meme.GetSanitizedText(results)

		fmt.Println(words)

		s.Stop()

		output := updateFileName(input, ".exif")

		_, err = cp(input, output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", exiterrorf.Errorf(err, "failed to copy output file"))
			os.Exit(1)
		}

		wr, err := os.Open(output) // lint:allow_possible_insecure
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", exiterrorf.Errorf(err, "failed to open the output file pointer"))
			os.Exit(1)
		}

		// exif

		err = r.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", exiterrorf.Errorf(err, "failed to close the input file pointer"))
			os.Exit(1)
		}

		err = wr.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", exiterrorf.Errorf(err, "failed to close the output file pointer"))
			os.Exit(1)
		}
	}
}

func updateFileName(s, a string) string {
	dir := filepath.Dir(s)
	base := filepath.Base(s)
	ext := filepath.Ext(s)
	realbase := strings.ReplaceAll(base, ext, "")

	return filepath.Join(dir, realbase+a+ext)
}

func cp(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("could not stat %s", src))
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src) // lint:allow_errorf
	}

	source, err := os.Open(src) // lint:allow_possible_insecure
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("could not open %s", src))
	}

	defer func() {
		_ = source.Close()
	}()

	destination, err := os.Create(dst) // lint:allow_possible_insecure
	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("could not open %s", dst))
	}

	defer func() {
		_ = destination.Close()
	}()

	nBytes, err := io.Copy(destination, source)

	return nBytes, errors.Wrap(err, fmt.Sprintf("could not copy from %s → %s", src, dst))
}
