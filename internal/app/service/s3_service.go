package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"propmanager/internal/config"
)

type S3Service struct {
	cfg *config.Config
}

func NewS3Service(cfg *config.Config) *S3Service {
	return &S3Service{cfg: cfg}
}

// generateRandomPrefix creates a random string to be used as a filename prefix.
func generateRandomPrefix(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *S3Service) UploadImage(ctx context.Context, file []byte, fileName string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String(s.cfg.S3Endpoint),
		Region:   aws.String(s.cfg.S3Region),
	}, nil)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	// Generate a random prefix for the filename to prevent collisions.
	prefix, err := generateRandomPrefix(4) // generates a random 8 character hex string
	if err != nil {
		log.Printf("Error generating random prefix: %v", err)
		return "", err
	}

	// Append the random prefix to the original filename.
	modifiedFileName := fmt.Sprintf("%s-%s", prefix, fileName)

	_, err = s3.New(sess).PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.cfg.S3Bucket),
		Key:    aws.String(modifiedFileName),
		Body:   aws.ReadSeekCloser(strings.NewReader(string(file))),
		ACL:    aws.String("public-read"),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists:
				log.Println("Bucket name already in use!")
			default:
				log.Printf("Unknown S3 error: %v", aerr)
			}
		} else {
			log.Printf("Non-S3 error: %v", err)
		}
		return "", err
	}

	// Return the full URL including the modified file name.
	return fmt.Sprintf("%s/%s", s.cfg.S3Endpoint, modifiedFileName), nil
}

func (s *S3Service) DeleteImage(imageUrl string) error {
	sess, err := session.NewSession(&aws.Config{
		Endpoint: aws.String(s.cfg.S3Endpoint),
		Region:   aws.String(s.cfg.S3Region),
	}, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Extract the key from imageUrl
	urlParts := strings.Split(imageUrl, "/")
	key := urlParts[len(urlParts)-1]

	_, err = s3.New(sess).DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.cfg.S3Bucket),
		Key:    aws.String(key),
	})
	return err
}
