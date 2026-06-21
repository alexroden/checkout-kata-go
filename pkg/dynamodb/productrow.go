package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type ProductRow struct {
	Sku       string `dynamodbav:"sku"`
	UnitPrice int    `dynamodbav:"unitPrice"`
}

func (r *ProductRow) Unmarshal(values map[string]types.AttributeValue) error {
	return attributevalue.UnmarshalMap(values, r)
}

func (r *ProductRow) Marshal() (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMap(r)
}
