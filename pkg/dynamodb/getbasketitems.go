package dynamodb

import (
	"encoding/json"
	"log/slog"

	"github.com/alexroden/checkout-kata-go/pkg/errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (d *DynamoDB) GetBasketItems(basketId string) ([]*BasketItemRow, error) {
	tableName, err := d.getTable(BASKET_ITEMS_TABLE)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		KeyConditionExpression: aws.String("basketId = :basketId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":basketId": &types.AttributeValueMemberS{
				Value: basketId,
			},
		},
	}

	b, _ := json.Marshal(input)
	slog.Debug(string(b))

	resp, err := d.db.Query(d.ctx, input)
	if err != nil {
		slog.Error("get basket operation fail: " + err.Error())

		return nil, errors.InternalServerError("get basket operation fail")
	}

	if len(resp.Items) == 0 {
		return nil, errors.NotFound("basket items not found")
	}

	result := []*BasketItemRow{}
	for _, item := range resp.Items {
		row := &BasketItemRow{}
		if err := row.Unmarshal(item); err != nil {
			slog.Error("unmarhal basket item fail: " + err.Error())

			return nil, errors.InternalServerError("unmarhal basket item fail")
		}

		result = append(result, row)
	}

	return result, nil
}
