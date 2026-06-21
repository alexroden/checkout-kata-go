package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type BasketRow struct {
	Id string `dynamodbav:"id"`
}

func (r *BasketRow) Unmarshal(values map[string]types.AttributeValue) error {
	return attributevalue.UnmarshalMap(values, r)
}

func (r *BasketRow) Marshal() (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMap(r)
}
