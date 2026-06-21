package dynamodb

import (
	"encoding/json"
	"log/slog"

	"github.com/alexroden/checkout-kata-go/pkg/errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (d *DynamoDB) GetProduct(sku string) (*ProductRow, error) {
	tableName, err := d.getTable(PRODUCTS_TABLE)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"sku": &types.AttributeValueMemberS{Value: sku},
		},
	}

	b, _ := json.Marshal(input)
	slog.Debug(string(b))

	resp, err := d.db.GetItem(d.ctx, input)
	if err != nil {
		slog.Error("get product operation fail: " + err.Error())

		return nil, errors.InternalServerError("get product operation fail")
	}

	if resp.Item == nil {
		return nil, errors.NotFound("product not found")
	}

	result := &ProductRow{}
	if err := result.Unmarshal(resp.Item); err != nil {
		slog.Error("unmarhal product fail: " + err.Error())

		return nil, errors.InternalServerError("unmarshal product fail")
	}

	return result, nil
}
