package meme

import (
	"bytes"
	"context"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"github.com/northwood-labs/golang-utils/exiterrorf"
)

type void struct{}

var (
	// DefaultPunctuation represents the default punctuation characters to strip.
	DefaultPunctuation = []byte("`~!@#$%^&*()_+[]\\{}|;':\",./<>?¡™£¢∞§¶•ªº–—=±“”‘’…æÆ≤¯≥˘÷¿Œ∑„´‰†¥ˆøπƒ©˚¬≈∫")
	member             void
)

// DetectText wraps the low-level SDK interface, handling the AWS stuff, and
// returning a list of results.
func DetectText(
	ctx context.Context,
	awsConfig *aws.Config,
	imageData bytes.Buffer,
) ([]*string, error) {
	rekogClient := rekognition.NewFromConfig(*awsConfig)
	var emptyOut []*string

	response, err := rekogClient.DetectText(ctx, &rekognition.DetectTextInput{
		Image: &types.Image{
			Bytes: imageData.Bytes(),
		},
	})
	if err != nil {
		return emptyOut, exiterrorf.Errorf(err, "")
	}

	collect := []*string{}

	// Collect lines of text (raw).
	for i := range response.TextDetections {
		textDetection := &response.TextDetections[i]
		collect = append(collect, textDetection.DetectedText)
	}

	return collect, nil
}

// GetSanitizedText sanitizes and de-dupes the resulting words.
func GetSanitizedText(lines []*string) []string {
	words := make(map[string]void)
	outWords := []string{}

	for i := range lines {
		line := *lines[i]
		line = strings.TrimSpace(line)
		ws := strings.Split(line, " ")

		for j := range ws {
			word := ws[j]

			word = strings.ToLower(word)
			word = strings.Trim(word, string(DefaultPunctuation))

			words[word] = member
		}
	}

	// I hate having to loop twice, but Go lacks a Set datatype.
	for k := range words {
		outWords = append(outWords, k)
	}

	// Deterministic output.
	sort.Strings(outWords)

	return outWords
}
