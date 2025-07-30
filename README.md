# ScribeShot

A command-line tool for extracting text from images and PDFs using AWS Textract.

## Features

- **Document Text Extraction**: Extract text from images (JPG, JPEG, PNG, TIFF) and PDF files
- **AWS Profile Support**: Use specific AWS profiles from your `~/.aws/credentials` file
- **Cost Tracking**: Real-time cost estimation for AWS Textract API calls
- **File Validation**: Automatic validation of file existence, type, and size
- **Error Handling**: Comprehensive error handling with user-friendly messages
- **Progress Feedback**: Shows file processing information and extraction statistics
- **Help System**: Built-in help and usage information

## Usage

### Basic Usage

```bash
scribeshot document.pdf
scribeshot image.jpg
```

### Using AWS Profiles and Regions

```bash
scribeshot --profile dev-account document.pdf
scribeshot --profile production --region us-west-2 image.png
scribeshot --region eu-west-1 document.pdf
```

### Help and Version

```bash
scribeshot --help
scribeshot --version
```

## Installation

### Prerequisites

- Go 1.19 or later
- AWS CLI configured with appropriate credentials
- AWS account with Textract permissions

### Build from Source; Add to path

```bash
go mod init scribeshot
go get github.com/aws/aws-sdk-go
go build -o scribeshot main.go

chmod +x ./scribeshot
sudo mv ./scribeshot /usr/bin/
which scribeshot
```

## Dependencies

- **AWS SDK for Go**: `github.com/aws/aws-sdk-go`
- **Go Standard Library**: `flag`, `fmt`, `os`, `path/filepath`

## AWS Configuration

### Credentials

The tool uses the AWS SDK's default credential and region resolution chain:

**Credentials:**

1. Environment variables (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
2. Shared credentials file (`~/.aws/credentials`)
3. IAM roles (when running on EC2)

**Region:**

1. `--region` flag
2. `AWS_REGION` environment variable
3. Profile-specific region in `~/.aws/config`
4. Default region in `~/.aws/config`

### Required IAM Permissions

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": ["textract:DetectDocumentText"],
      "Resource": "*"
    }
  ]
}
```

### AWS Profile Configuration

Add profiles to `~/.aws/credentials`:

```ini
[default]
aws_access_key_id = YOUR_ACCESS_KEY
aws_secret_access_key = YOUR_SECRET_KEY

[dev-account]
aws_access_key_id = DEV_ACCESS_KEY
aws_secret_access_key = DEV_SECRET_KEY
```

Add regions to `~/.aws/config`:

```ini
[default]
region = us-east-1

[profile dev-account]
region = us-west-2
```

## Supported File Types

- **Images**: JPG, JPEG, PNG, TIFF, TIF
- **Documents**: PDF
- **Size Limit**: 10MB (AWS Textract synchronous limit)

## Example Output

```bash
$ scribeshot --profile dev --region us-west-2 sample.pdf
Using AWS region: us-west-2
Processing file: sample.pdf (2048576 bytes)
This is the first line of text from the document.
Here is another line with some important information.
The document contains multiple paragraphs of text.
Extracted 23 lines of text from 2 page(s)
Estimated cost: $0.0030 USD (DetectDocumentText @ $0.0015 per page)
```

## Error Handling

The tool provides clear error messages for common issues:

- File not found
- Unsupported file types
- File size exceeding limits
- AWS authentication errors
- Textract API errors

## TODOs

### High Priority

- [ ] **Output Formats**: Add support for JSON and structured output formats
- [ ] **Batch Processing**: Support for processing multiple files at once
- [ ] **S3 Integration**: Direct processing of files from S3 buckets

### Medium Priority

- [ ] **Confidence Scores**: Display confidence scores for extracted text
- [ ] **Table Detection**: Add support for table extraction using `AnalyzeDocument`
- [ ] **Form Detection**: Extract key-value pairs from forms
- [ ] **Configuration File**: Support for YAML/JSON configuration files
- [ ] **Logging**: Add structured logging with different verbosity levels

### Low Priority

- [ ] **Progress Bar**: Visual progress indicator for large files
- [ ] **Output Filtering**: Filter output by confidence threshold
- [ ] **Text Post-processing**: Clean up extracted text (remove extra whitespace, etc.)
- [ ] **Advanced Cost Tracking**: Track monthly usage and tier-based pricing
- [ ] **Cost Budgets**: Warning when approaching spending thresholds

## Nice to Haves

### Developer Experience

- [ ] **Docker Container**: Containerized version for easy deployment
- [ ] **GitHub Actions**: CI/CD pipeline for automated testing and releases
- [ ] **Homebrew Formula**: Easy installation on macOS
- [ ] **Cross-compilation**: Pre-built binaries for different platforms

### Advanced Features

- [ ] **Async Processing**: Support for large documents using async Textract jobs
- [ ] **OCR Comparison**: Compare results with other OCR engines
- [ ] **Text Analysis**: Integration with AWS Comprehend for sentiment/entity analysis
- [ ] **Cost Tracking**: ~~Estimate AWS costs for processing operations~~ ✅ **COMPLETED**

### Infrastructure as Code

- [ ] **CDK Module**: AWS CDK constructs for Textract resources
- [ ] **Terragrunt Configuration**: Example Terragrunt setup for multi-environment deployments
- [ ] **CloudFormation Templates**: Alternative deployment option
- [ ] **AWS SAM Template**: Serverless deployment option

### Monitoring & Observability

- [ ] **CloudWatch Integration**: Send metrics to CloudWatch
- [ ] **X-Ray Tracing**: Distributed tracing for performance monitoring
- [ ] **Prometheus Metrics**: Export metrics in Prometheus format

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the GNU GENERAL PUBLIC LICENSE - see the LICENSE file for details.

## Support

For issues related to:

- **AWS Textract**: Check the [AWS Textract documentation](https://docs.aws.amazon.com/textract/)
- **AWS SDK**: Refer to the [AWS SDK for Go documentation](https://docs.aws.amazon.com/sdk-for-go/)
- **This tool**: Open an issue in the GitHub repository
