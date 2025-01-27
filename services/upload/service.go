package upload

import (
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"golang.org/x/exp/rand"
)

type UploadService struct {
	svc *s3.S3
}

func NewUploadService() (*UploadService, error) {
	var sess, err = session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	if err != nil {
		return nil, err
	}

	return &UploadService{
		svc: s3.New(sess),
	}, nil
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateImageName(length int) string {
	rand.Seed(uint64(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (u *UploadService) GetPresignedURL() (string, error) {
	imageName := GenerateImageName(8)

	req, _ := u.svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String("kanteentest"),
		Key:    aws.String(imageName),
	})

	str, err := req.Presign(15 * time.Minute)

	if err != nil {
		log.Println("ERROR: GetPresignURL", err)
		return "", err
	}

	return str, nil
}
