package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	mocks "github.com/alexroden/checkout-kata-go/internal/mocks/pkg/repositories"
	"github.com/alexroden/checkout-kata-go/pkg/repositories"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/suite"
)

type DynamoInitSuite struct {
	suite.Suite
	ctx       context.Context
	tableName string
}

func (s *DynamoInitSuite) SetupTest() {
	s.ctx = context.Background()
	s.tableName = "ProductsTable"
}

func (s *DynamoInitSuite) DynamoDB(err error) repositories.DynamoDBAPI {
	result := &mocks.MockDynamoDBAPI{}

	result.On(
		"BatchWriteItem",
		s.ctx,
		&dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				s.tableName: {
					{
						PutRequest: &types.PutRequest{
							Item: map[string]types.AttributeValue{
								"sku": &types.AttributeValueMemberS{
									Value: "A",
								},
								"unitPrice": &types.AttributeValueMemberN{
									Value: "10",
								},
							},
						},
					},
					{
						PutRequest: &types.PutRequest{
							Item: map[string]types.AttributeValue{
								"sku": &types.AttributeValueMemberS{
									Value: "B",
								},
								"unitPrice": &types.AttributeValueMemberN{
									Value: "20",
								},
							},
						},
					},
				},
			},
		},
	).Return(nil, err).Once()

	return result
}

func (s *DynamoInitSuite) Test_LoadTable() {
	dir := s.T().TempDir()

	file := filepath.Join(dir, "products.json")

	data := fmt.Sprintf(`{
		"TableName": "%s",
		"BillingMode": "PAY_PER_REQUEST",
		"Attributes": [
			{"Name": "Sku", "Type": "S"}
		],
		"KeySchema": [
			{"Name": "Sku", "KeyType": "HASH"}
		]
	}`, s.tableName)

	s.NoError(os.WriteFile(file, []byte(data), 0644))

	tables, err := loadTables(dir)
	s.NoError(err)

	s.Len(tables, 1)
	s.Equal(s.tableName, tables[0].TableName)
}

func (s *DynamoInitSuite) Test_SeedTable() {
	dir := s.T().TempDir()

	file := filepath.Join(dir, "products.json")

	items := []map[string]interface{}{
		{"sku": "A", "unitPrice": 10},
		{"sku": "B", "unitPrice": 20},
	}

	b, err := json.Marshal(items)
	s.NoError(err)

	s.NoError(os.WriteFile(file, b, 0644))

	s.NoError(seedTable(s.ctx, s.DynamoDB(nil), s.tableName, file))
}

func (s *DynamoInitSuite) Test_LoadTable_DirectoryNotFound() {
	tables, err := loadTables("/does/not/exist")

	s.Error(err)
	s.Nil(tables)
}

func (s *DynamoInitSuite) Test_LoadTable_InvalidJSON() {
	dir := s.T().TempDir()

	file := filepath.Join(dir, "products.json")

	s.NoError(os.WriteFile(file, []byte("{invalid json"), 0644))

	tables, err := loadTables(dir)

	s.Error(err)
	s.Nil(tables)
	s.Contains(err.Error(), "failed parsing")
}

func (s *DynamoInitSuite) Test_SeedTable_FileNotFound() {
	mock := &mocks.MockDynamoDBAPI{}

	err := seedTable(
		s.ctx,
		mock,
		s.tableName,
		"/does/not/exist.json",
	)

	s.Error(err)
}

func (s *DynamoInitSuite) Test_SeedTable_InvalidJSON() {
	dir := s.T().TempDir()

	file := filepath.Join(dir, "products.json")

	s.NoError(
		os.WriteFile(file, []byte("{invalid json"), 0644),
	)

	mock := &mocks.MockDynamoDBAPI{}

	err := seedTable(
		s.ctx,
		mock,
		s.tableName,
		file,
	)

	s.Error(err)
}

func (s *DynamoInitSuite) Test_SeedTable_BatchWriteFails() {
	dir := s.T().TempDir()

	file := filepath.Join(dir, "products.json")

	items := []map[string]interface{}{
		{"sku": "A", "unitPrice": 10},
		{"sku": "B", "unitPrice": 20},
	}

	b, err := json.Marshal(items)
	s.NoError(err)

	s.NoError(os.WriteFile(file, b, 0644))

	expectedErr := fmt.Errorf("dynamo write failed")
	err = seedTable(s.ctx, s.DynamoDB(expectedErr), s.tableName, file)
	s.Error(err)

	s.ErrorIs(err, expectedErr)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestDynamoInitSuite(t *testing.T) {
	suite.Run(t, new(DynamoInitSuite))
}
