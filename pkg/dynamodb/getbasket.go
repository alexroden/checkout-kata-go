package dynamodb

import (
	"encoding/json"
	"log/slog"

	"github.com/alexroden/checkout-kata-go/pkg/errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (d *DynamoDB) GetBasket(id string) (*BasketRow, error) {
	tableName, err := d.getTable(BASKETS_TABLE)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	}

	b, _ := json.Marshal(input)
	slog.Debug(string(b))

	resp, err := d.db.GetItem(d.ctx, input)
	if err != nil {
		slog.Error("get basket operation fail: " + err.Error())

		return nil, errors.InternalServerError("get basket operation fail")
	}

	if resp.Item == nil {
		return nil, errors.NotFound("basket not found")
	}

	result := &BasketRow{}
	if err := result.Unmarshal(resp.Item); err != nil {
		slog.Error("unmarhal basket fail: " + err.Error())

		return nil, errors.InternalServerError("unmarshal basket fail")
	}

	return result, nil
}
