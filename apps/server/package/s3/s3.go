package s3

import (
	"context"
	"errors"
	"net/http"

	"game-random-api/package/config"
	"game-random-api/package/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	s3cfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/traphamxuan/gobs"
)

type S3Config struct {
	AccessKey string `env:"S3_ACCESS_KEY"`
	SecretKey string `env:"S3_SECRET_KEY"`
	Region    string `env:"S3_REGION"`
	Bucket    string `env:"S3_BUCKET_NAME"`
	URI       string `env:"S3_URI"`
}

type S3 struct {
	log    *logger.Logger
	config *S3Config
	client *s3.Client
}

var _ gobs.IService = (*S3)(nil)

func NewS3(ctx context.Context) *S3 {
	return &S3{}
}

func (s *S3) Init(ctx context.Context, sb *gobs.Component) error {
	sb.Deps = []gobs.IService{
		&logger.Logger{},
		&config.Configuration{},
	}

	onnSetup := func(ctx context.Context, dependencies []gobs.IService, _ []gobs.CustomService) error {
		s.log = dependencies[0].(*logger.Logger)
		config := dependencies[1].(*config.Configuration)
		var cfg S3Config
		if err := config.ParseConfig(&cfg); err != nil {
			return err
		}
		s.config = &cfg
		s.Setup(ctx)
		return nil
	}
	sb.OnSetup = &onnSetup

	return nil
}

func (s *S3) Setup(ctx context.Context) error {

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: s.config.URI,
		}, nil
	})

	s3Cfg, err := s3cfg.LoadDefaultConfig(ctx,
		s3cfg.WithEndpointResolverWithOptions(r2Resolver),
		s3cfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s.config.AccessKey, s.config.SecretKey, "")),
		s3cfg.WithRegion("auto"),
	)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(s3Cfg)

	s.client = client
	return nil
}

func (s *S3) IsExisted(ctx context.Context, key string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.config.Bucket),
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
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return resp.Metadata, nil
}

func (s *S3) CopyObject(c context.Context, sourceKey string, destKey string, metadata map[string]string) error {
	_, err := s.client.CopyObject(c, &s3.CopyObjectInput{
		Bucket:            aws.String(s.config.Bucket),
		CopySource:        aws.String(s.config.Bucket + "/" + sourceKey),
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
		Bucket: aws.String(s.config.Bucket),
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
