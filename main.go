package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/private/protocol/rest"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	AWS_ACCESS_KEY_ID     = "<AWS_ACCESS_KEY_ID>"
	AWS_SECRET_ACCESS_KEY = "<AWS_SECRET_ACCESS_KEY>"
	AWS_REGION            = "<AWS_REGION>"
	AWS_BUCKET            = "<AWS_BUCKET>"
	AWS_BUCKET_FOLDER     = "<AWS_BUCKET_FOLDER>"
)

type AWSConfig struct {
	AccessKeyId     string
	AccessKeySecret string
	Region          string
	Bucket          string
	BucketFolder    string
	UploadTimeout   int
	BaseURL         string
}

func newAwsConfig() AWSConfig {
	return AWSConfig{
		AccessKeyId:     AWS_ACCESS_KEY_ID,
		AccessKeySecret: AWS_SECRET_ACCESS_KEY,
		Region:          AWS_REGION,
		Bucket:          AWS_BUCKET,
		BucketFolder:    AWS_BUCKET_FOLDER,
	}
}

func createSession(awsConfig AWSConfig) *session.Session {
	sess := session.Must(session.NewSession(
		&aws.Config{
			Region: aws.String(awsConfig.Region),
			Credentials: credentials.NewStaticCredentials(
				awsConfig.AccessKeyId,
				awsConfig.AccessKeySecret,
				"",
			),
		},
	))
	return sess
}

func createS3Session(sess *session.Session) *s3.S3 {
	s3Session := s3.New(sess)
	return s3Session
}

func UploadObject(bucket string, filePath string, sess *session.Session, awsConfig AWSConfig) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	fileName := awsConfig.BucketFolder + "/" + filepath.Base(file.Name())
	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		Body:   file,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully uploaded %q to %q\n", fileName, bucket)
	return nil
}

func GetObjectLink(bucket string, fileName string, sess *session.Session, awsConfig AWSConfig) error {
	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})
	rest.Build(req)

	fmt.Printf("Object %q link: %q \n", fileName, req.HTTPRequest.URL.String())
	return nil
}

func GetObjectSecureLink(bucket string, fileName string, sess *session.Session, awsConfig AWSConfig) error {
	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})

	url, err := req.Presign(15 * time.Minute)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Object %q secure link: %q \n", fileName, url)
	return nil
}

func main() {
	awsConfig := newAwsConfig()
	filePath := "./dummy.jpg"
	sess := createSession(awsConfig)

	_ = UploadObject(awsConfig.Bucket, filePath, sess, awsConfig)
	_ = GetObjectLink(awsConfig.Bucket, "dummies/dummies.pdf", sess, awsConfig)
	_ = GetObjectSecureLink(awsConfig.Bucket, "dummies/dummies.pdf", sess, awsConfig)
}
