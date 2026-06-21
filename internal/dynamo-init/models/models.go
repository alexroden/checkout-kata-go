package models

import "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

type DynamoTable struct {
	TableName   string
	BillingMode types.BillingMode
	Attributes  []*DynamoAttribute
	KeySchema   []*DynamoKey
}

type DynamoAttribute struct {
	Name string
	Type types.ScalarAttributeType
}

type DynamoKey struct {
	Name    string
	KeyType types.KeyType
}
