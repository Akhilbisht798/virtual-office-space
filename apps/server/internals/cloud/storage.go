package cloud

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetPreSignedUrl(
	bucketName string, objecetKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	s3Client, err := getS3Client()
	if err != nil {
		return nil, err
	}
	client := s3.NewPresignClient(&s3Client)
	url, err := client.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objecetKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Duration(lifetimeSecs * int64(time.Minute))
	})
	if err != nil {
		log.Printf("Could'nt get a Presigned request to get %v:%v, Here why: %v\n",
			bucketName, objecetKey, err)
		return nil, err
	}
	return url, nil
}

func PutPreSignedUrl(
	bucketName string, objecetKey string, lifetimeSecs int64) (*v4.PresignedHTTPRequest, error) {
	s3Client, err := getS3Client()
	if err != nil {
		return nil, err
	}
	client := s3.NewPresignClient(&s3Client)
	url, err := client.PresignPutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objecetKey),
	}, func(po *s3.PresignOptions) {
		po.Expires = time.Duration(lifetimeSecs * int64(time.Minute))
	})
	if err != nil {
		log.Printf("Could'nt get a Presigned request to get %v:%v, Here why: %v\n",
			bucketName, objecetKey, err)
		return nil, err
	}
	return url, nil
}

func getS3Client() (s3.Client, error) {
	var client s3.Client
	region := os.Getenv("REGION")
	if os.Getenv("APP_ENV") == "production" {
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithDefaultRegion(region),
		)
		if err != nil {
			return s3.Client{}, err
		}
		client = *s3.NewFromConfig(cfg)
		return client, nil
	}

	// Devlopment using MinIO
	//log.Println("using minio for devlopment side of thing")
	staticResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			URL:               "http://localhost:9000",
			SigningRegion:     region,
			HostnameImmutable: true,
		}, nil
	})

	cfg := aws.Config{
		Region:           region,
		Credentials:      credentials.NewStaticCredentialsProvider("ROOT", "password", ""),
		EndpointResolver: staticResolver,
	}

	client = *s3.NewFromConfig(cfg)
	return client, nil
}
