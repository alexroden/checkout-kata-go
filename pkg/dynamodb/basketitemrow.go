package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type BasketItemRow struct {
	BasketId string `dynamodbav:"basketId"`
	Sku      string `dynamodbav:"sku"`
	Quantity int    `dynamodbav:"quantity"`
}

func (r *BasketItemRow) Unmarshal(values map[string]types.AttributeValue) error {
	return attributevalue.UnmarshalMap(values, r)
}

func (r *BasketItemRow) Marshal() (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMap(r)
}
