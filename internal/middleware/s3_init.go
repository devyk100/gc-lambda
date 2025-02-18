package middleware

import (
	"context"

	"gc.yashk.dev/lambda/internal/env"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/s3/actions"
)

func InitS3(ctx context.Context) (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(env.Region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(env.Key, env.Secret, ""),
		))
	s3Client := s3.NewFromConfig(cfg)
	if err != nil {
		return nil, err
	}
	return s3Client, nil
}

func InitS3PresignClient(ctx context.Context) (*s3.PresignClient, error) {
	s3Client, err := InitS3(ctx)
	presignClient := s3.NewPresignClient(s3Client)
	if err != nil {
		return nil, err
	}
	return presignClient, err
}

func InitS3Presigner(ctx context.Context) (actions.Presigner, error) {
	presignClient, err := InitS3PresignClient(ctx)
	presigner := actions.Presigner{PresignClient: presignClient}
	if err != nil {
		return presigner, err
	}
	return presigner, nil
}
