package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var s3Region, s3Endpoint, s3Bucket string
var s3Credentials *credentials.Credentials

func init() {
	var ok bool
	s3Region, ok = os.LookupEnv("S3_REGION")
	if !ok {
		log.Fatal("S3_REGION environment variable missing")
	}

	s3Endpoint, ok = os.LookupEnv("S3_ENDPOINT")
	if !ok {
		log.Fatal("S3_ENDPOINT environment variable missing")
	}

	s3Bucket, ok = os.LookupEnv("S3_BUCKET")
	if !ok {
		log.Fatal("S3_BUCKET environment variable missing")
	}

	s3Credentials = credentials.NewEnvCredentials()
}

func listOfPandas() (pandas []string, err error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(s3Endpoint),
		Region:      aws.String(s3Region),
		Credentials: s3Credentials,
	})
	if err != nil {
		return
	}

	svc := s3.New(sess)
	resp, err := svc.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(s3Bucket),
		Marker: aws.String("source/"),
	})
	if err != nil {
		return
	}

	pandas = make([]string, len(resp.Contents))
	for i, item := range resp.Contents {
		pandas[i] = *item.Key
	}

	return
}

func downloadPanda(key string) (data []byte, err error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint:    aws.String(s3Endpoint),
		Region:      aws.String(s3Region),
		Credentials: s3Credentials,
	})
	if err != nil {
		return
	}

	buf := &aws.WriteAtBuffer{}

	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return
	}

	data = buf.Bytes()
	return
}
