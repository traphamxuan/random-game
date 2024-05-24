package s3

import (
	"context"
	"errors"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	s3cfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/traphamxuan/random-game/package/config"
	"github.com/traphamxuan/random-game/package/logger"
	servicemanager "github.com/traphamxuan/random-game/package/service_manager"
)

type S3 struct {
	log        *logger.Logger
	config     *config.Configuration
	client     *s3.Client
	bucketName string
}

var _ servicemanager.IService = (*S3)(nil)

type S3Config struct {
	AccessKey string `env:"S3_ACCESS_KEY"`
	SecretKey string `env:"S3_SECRET_KEY"`
	Region    string `env:"S3_REGION"`
	Bucket    string `env:"S3_BUCKET_NAME"`
	URI       string `env:"S3_URI"`
}

func NewS3(ctx context.Context, sm *servicemanager.ServiceManager) *S3 {
	return &S3{
		log:    servicemanager.GetServiceOrPanic[*logger.Logger](sm, "Logger", "S3"),
		config: servicemanager.GetServiceOrPanic[*config.Configuration](sm, "Configuration", "S3"),
	}
}

func (s *S3) Setup(ctx context.Context) error {
	var s3Config S3Config
	if err := s.config.ParseConfig(&s3Config); err != nil {
		return err
	}
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: s3Config.URI,
		}, nil
	})

	s3Cfg, err := s3cfg.LoadDefaultConfig(ctx,
		s3cfg.WithEndpointResolverWithOptions(r2Resolver),
		s3cfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s3Config.AccessKey, s3Config.SecretKey, "")),
		s3cfg.WithRegion("auto"),
	)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(s3Cfg)

	s.client = client
	s.bucketName = s3Config.Bucket
	return nil
}

func (s *S3) Start(ctx context.Context) error {
	return nil
}

func (s *S3) Stop(ctx context.Context) error {
	return nil
}

func (s *S3) IsExisted(ctx context.Context, key string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		var responseError *awshttp.ResponseError
		if errors.As(err, &responseError) && responseError.ResponseError.HTTPStatusCode() == http.StatusNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *S3) FetchMetadata(c context.Context, key string) (map[string]string, error) {
	resp, err := s.client.HeadObject(c, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return resp.Metadata, nil
}

func (s *S3) CopyObject(c context.Context, sourceKey string, destKey string, metadata map[string]string) error {
	_, err := s.client.CopyObject(c, &s3.CopyObjectInput{
		Bucket:            aws.String(s.bucketName),
		CopySource:        aws.String(s.bucketName + "/" + sourceKey),
		Key:               aws.String(destKey),
		Metadata:          metadata,
		MetadataDirective: types.MetadataDirectiveReplace,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *S3) DeleteObject(c context.Context, key string) error {
	_, err := s.client.DeleteObject(c, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *S3) CommitFileUploaded(c context.Context, key string, parent string, target string, prefix FilePrefix) (string, func(*string) error, func() error, error) {
	metadata, err := s.FetchMetadata(c, "upload/"+key)
	if err != nil {
		return "", nil, nil, err
	}
	if metadata == nil || metadata["parent"] != parent {
		return "", nil, nil, errors.New("invalid file")
	}
	// prodKey := slug.Make("product " + product.Name + " " + strconv.Itoa(int(time.Now().Unix()%100)))
	prodKey := string(prefix) + "/" + parent + "/" + target
	if err := s.CopyObject(c, "upload/"+key, "data/"+prodKey, map[string]string{
		"id":       target,
		"parent":   parent,
		"filename": metadata["filename"],
	}); err != nil {
		return "", nil, nil, err
	}

	finish := func(oldKey *string) error {
		if oldKey != nil && *oldKey != "" {
			if err := s.DeleteObject(c, *oldKey); err != nil {
				return err
			}
		}
		return s.DeleteObject(c, "upload/"+key)
	}
	rollback := func() error {
		return s.DeleteObject(c, "data/"+prodKey)
	}

	return prodKey, finish, rollback, nil
}
