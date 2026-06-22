package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type SpecialRow struct {
	Sku      string `dynamodbav:"sku"`
	Quantity int    `dynamodbav:"quantity"`
	Price    int    `dynamodbav:"price"`
}

func (r *SpecialRow) Unmarshal(values map[string]types.AttributeValue) error {
	return attributevalue.UnmarshalMap(values, r)
}

func (r *SpecialRow) Marshal() (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMap(r)
}
