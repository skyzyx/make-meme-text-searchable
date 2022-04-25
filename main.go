package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dustin/go-humanize"
	"github.com/northwood-labs/awsutils"
	"github.com/northwood-labs/golang-utils/debug"
	"github.com/northwood-labs/golang-utils/exiterrorf"
	"github.com/northwood-labs/golang-utils/log"
	flag "github.com/spf13/pflag"

	"github.com/skyzyx/make-meme-text-searchable/meme"
)

const defaultRetries = 5

func main() {
	pp := debug.GetSpew()
	ctx := context.Background()
	logger := log.GetStdTextLogger()

	var flagProfile string
	var flagRegion string
	var flagRetries int
	var flagVerbose bool

	flag.StringVarP(&flagProfile, "profile", "p", "", "The AWS CLI profile with which to perform the request.")
	flag.StringVarP(&flagRegion, "region", "r", "", "The AWS region in which to perform the request.")
	flag.IntVarP(&flagRetries, "retries", "t", defaultRetries, "The max number of times to retry failed AWS requests.")
	flag.BoolVarP(&flagVerbose, "verbose", "v", false, "Enable verbose logging.")
	flag.Parse()

	config, err := awsutils.GetAWSConfig(ctx, flagRegion, flagProfile, flagRetries, flagVerbose)
	if err != nil {
		exiterrorf.ExitErrorf(err)
	}

	r, err := os.Open("/Library/WebServer/Documents/make-meme-text-searchable/images/paris-airport.heic")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", exiterrorf.Errorf(err, "failed to open file"))
		os.Exit(1)
	}

	buf, _, _, err := meme.ReadImage(r, meme.DefaultJPEGQuality)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	if flagVerbose {
		logger.Info().Str("Size", humanize.Bytes(uint64(len(buf.Bytes())))).Msg("Rekognition limit for images is 5 MB.")
	}

	results, err := meme.DetectText(ctx, &config, buf)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	words := meme.GetSanitizedText(results)

	pp.Dump(words)
}
