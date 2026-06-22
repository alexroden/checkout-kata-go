package dynamodb

import (
	"encoding/json"
	"log/slog"

	"github.com/alexroden/checkout-kata-go/pkg/errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (d *DynamoDB) GetProducts(skus ...string) ([]*ProductRow, error) {
	tableName, err := d.getTable(PRODUCTS_TABLE)
	if err != nil {
		return nil, err
	}

	keys := make([]map[string]types.AttributeValue, 0, len(skus))
	for _, sku := range skus {
		keys = append(keys, map[string]types.AttributeValue{
			"sku": &types.AttributeValueMemberS{
				Value: sku,
			},
		})
	}

	input := &dynamodb.BatchGetItemInput{
		RequestItems: map[string]types.KeysAndAttributes{
			tableName: {
				Keys: keys,
			},
		},
	}

	b, _ := json.Marshal(input)
	slog.Debug(string(b))

	resp, err := d.db.BatchGetItem(d.ctx, input)
	if err != nil {
		slog.Error("get products operation fail: " + err.Error())

		return nil, errors.InternalServerError("get products operation fail")
	}

	result := []*ProductRow{}
	for _, row := range resp.Responses[tableName] {
		product := &ProductRow{}
		if err := product.Unmarshal(row); err != nil {
			slog.Error("unmarhal product fail: " + err.Error())

			return nil, errors.InternalServerError("unmarshal product fail")
		}

		result = append(result, product)
	}

	return result, nil
}
