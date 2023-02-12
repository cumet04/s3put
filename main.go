package main

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/signal"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	usage := "Usage: cat file | s3put s3://your-bucket-name/path/to/dest"

	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	u, err := url.ParseRequestURI(os.Args[1])
	if err != nil || u.Scheme != "s3" {
		fmt.Fprintf(os.Stderr, "failed to parse destination: %v\n\n", err)
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(1)
	}

	bucket := u.Host
	key := strings.TrimPrefix(u.Path, "/")
	body := io.MultiReader(os.Stdin)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if err = uploadObject(ctx, bucket, key, body); err != nil {
		fmt.Fprintf(os.Stderr, "failed to put object; %v", err)
	}
}

func uploadObject(ctx context.Context, bucket, key string, body io.Reader) error {
	conf, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	// MEMO: To upload unknown size data, use Uploader instead of client.PutObject
	// https://stackoverflow.com/questions/43595911/how-to-save-data-streams-in-s3-aws-sdk-go-example-not-working
	client := s3.NewFromConfig(conf)
	uploader := manager.NewUploader(client)
	_, err = uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   body,
	})
	if err != nil {
		return err
	}

	return nil
}
