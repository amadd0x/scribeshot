package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/textract"
)

const (
	appName    = "scribeshot"
	appVersion = "0.2.0"

	// AWS Textract pricing (as of 2025)
	// DetectDocumentText: $0.0015 per page for first 1M pages per month
	textractCostPerPage = 0.0015
)

var (
	sess *session.Session
	svc  *textract.Textract
)

func main() {
	var (
		help    = flag.Bool("help", false, "Show help message")
		profile = flag.String("profile", "", "AWS profile to use (defaults to AWS SDK default behavior)")
		region  = flag.String("region", "", "AWS region to use (defaults to profile/environment configuration)")
		version = flag.Bool("version", false, "Show version information")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s - AWS Textract CLI Tool v%s\n\n", appName, appVersion)
		fmt.Fprintf(os.Stderr, "USAGE:\n")
		fmt.Fprintf(os.Stderr, "  %s [OPTIONS] <filename>\n\n", appName)
		fmt.Fprintf(os.Stderr, "DESCRIPTION:\n")
		fmt.Fprintf(os.Stderr, "  Extract text from images and PDFs using AWS Textract\n\n")
		fmt.Fprintf(os.Stderr, "OPTIONS:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEXAMPLES:\n")
		fmt.Fprintf(os.Stderr, "  %s document.pdf\n", appName)
		fmt.Fprintf(os.Stderr, "  %s --profile dev-account image.jpg\n", appName)
		fmt.Fprintf(os.Stderr, "  %s --help\n", appName)
	}

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *version {
		fmt.Printf("v%s\n", appVersion)
		os.Exit(0)
	}

	// Check for filename argument
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Error: Please provide exactly one filename\n")
		os.Exit(1)
	}

	filename := flag.Arg(0)

	// Validate file exists and is readable
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: File '%s' does not exist\n", filename)
		os.Exit(1)
	}

	// Initialize AWS session with optional profile
	if err := initAWS(*profile, *region); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing AWS session: %v\n", err)
		os.Exit(1)
	}

	// Process the document
	if err := processDocument(filename); err != nil {
		fmt.Fprintf(os.Stderr, "Error processing document: %v\n", err)
		os.Exit(1)
	}
}

func initAWS(profileName, regionName string) error {
	opts := session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}

	// If profile is specified, use it
	if profileName != "" {
		opts.Profile = profileName
	}

	// If region is specified, use it
	if regionName != "" {
		opts.Config = aws.Config{
			Region: aws.String(regionName),
		}
	}

	var err error
	sess, err = session.NewSessionWithOptions(opts)
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %w", err)
	}

	// Verify we have a region configured
	if sess.Config.Region == nil || *sess.Config.Region == "" {
		return fmt.Errorf("no AWS region configured. Please set a region via:\n"+
			"  1. --region flag\n"+
			"  2. AWS_REGION environment variable\n"+
			"  3. [%s] section in ~/.aws/config with region setting\n"+
			"  4. [default] section in ~/.aws/config with region setting",
			getProfileName(profileName))
	}

	svc = textract.New(sess)
	fmt.Fprintf(os.Stderr, "Using AWS region: %s\n", *sess.Config.Region)
	return nil
}

func getProfileName(profile string) string {
	if profile == "" {
		return "default"
	}
	return profile
}

func processDocument(filename string) error {
	// Validate file type (basic check)
	ext := filepath.Ext(filename)
	validExts := map[string]bool{
		".pdf": true, ".jpg": true, ".jpeg": true, ".png": true, ".tiff": true, ".tif": true,
	}
	if !validExts[ext] {
		return fmt.Errorf("unsupported file type: %s (supported: PDF, JPG, JPEG, PNG, TIFF)", ext)
	}

	file, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Check file size (Textract has limits)
	const maxFileSize = 10 * 1024 * 1024 // 10MB for synchronous operations
	if len(file) > maxFileSize {
		return fmt.Errorf("file size (%d bytes) exceeds Textract limit of %d bytes", len(file), maxFileSize)
	}

	fmt.Fprintf(os.Stderr, "Processing file: %s (%d bytes)\n", filename, len(file))

	resp, err := svc.DetectDocumentText(&textract.DetectDocumentTextInput{
		Document: &textract.Document{
			Bytes: file,
		},
	})
	if err != nil {
		return fmt.Errorf("Textract API error: %w", err)
	}

	// Count blocks for user feedback
	lineCount := 0
	pageCount := 0
	for _, block := range resp.Blocks {
		if *block.BlockType == "PAGE" {
			pageCount++
		}
		if *block.BlockType == "LINE" {
			lineCount++
			fmt.Println(*block.Text)
		}
	}

	// Calculate cost
	cost := float64(pageCount) * textractCostPerPage

	fmt.Fprintf(os.Stderr, "Extracted %d lines of text from %d page(s)\n", lineCount, pageCount)
	fmt.Fprintf(os.Stderr, "Estimated cost: $%.4f USD (DetectDocumentText @ $%.4f per page)\n", cost, textractCostPerPage)

	return nil
}
