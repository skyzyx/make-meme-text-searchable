package meme

import (
	"context"
	"image"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
)

// GetSanitizedText is a higher-level wrapper around the Amazon Rekognition API and the related sanitizing functions.
//
// * https://pkg.go.dev/context#Context
// * https://pkg.go.dev/image#Image
// * https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws#Config
func GetSanitizedText(ctx context.Context, awsConfig *aws.Config, imageData image.Image) ([]string, error) {
	return []string{}, nil
}

// DetectText wraps the low-level SDK interface, handling the AWS stuff.
//
// * https://pkg.go.dev/context#Context
// * https://pkg.go.dev/image#Image
// * https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws#Config
// * https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/rekognition#Client.DetectText
// * https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/rekognition#DetectTextInput
// * https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/rekognition/types#Image
// * https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/rekognition#DetectTextOutput
// * https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/rekognition/types#TextDetection
func DetectText(
	ctx context.Context,
	awsConfig *aws.Config,
	imageData image.Image,
) (*rekognition.DetectTextOutput, error) {
	// rekogClient := rekognition.NewFromConfig(*awsConfig)
	o := rekognition.DetectTextOutput{}

	return &o, nil
}
