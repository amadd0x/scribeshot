package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/textract"
)

var (
	sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc = textract.New(sess)
)

func main() {
	filename := os.Args[1]

	file, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	resp, err := svc.DetectDocumentText(&textract.DetectDocumentTextInput{
		Document: &textract.Document{
			Bytes: file,
		},
	})
	if err != nil {
		panic(err)
	}

	for _, block := range resp.Blocks {
		if *block.BlockType == "LINE" {
			fmt.Println(*block.Text)
		}
	}
}

func init() {
	svc = textract.New(sess)
}
