package dynamodb

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/stretchr/testify/suite"
)

type DynamoDBSuite struct {
	suite.Suite
	cnf *Config
}

func (s *DynamoDBSuite) SetupTest() {
	s.cnf = &Config{
		Region:   "eu-west-1",
		Endpoint: "http://dynamodb:8000",
		Tables: map[string]string{
			PRODUCTS_TABLE: PRODUCTS_TABLE,
		},
	}
}

func (s *DynamoDBSuite) Test_New() {
	repo, err := New(s.cnf)
	s.NoError(err)

	s.IsType(&DynamoDB{}, repo)
}

func (s *DynamoDBSuite) Test_ConfigError() {
	original := loadDefaultConfig
	defer func() { loadDefaultConfig = original }()

	loadDefaultConfig = func(
		ctx context.Context,
		optFns ...func(*config.LoadOptions) error,
	) (aws.Config, error) {
		return aws.Config{}, errors.New("boom")
	}

	repo, err := New(s.cnf)
	s.Error(err)

	s.Nil(repo)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestDynamoDBSuite(t *testing.T) {
	suite.Run(t, new(DynamoDBSuite))
}
