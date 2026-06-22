package dynamodb

import (
	"encoding/json"
	"log/slog"

	"github.com/alexroden/checkout-kata-go/pkg/errors"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (d *DynamoDB) GetSpecials(skus ...string) ([]*SpecialRow, error) {
	tableName, err := d.getTable(SPECIALS_TABLE)
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
		slog.Error("get specials operation fail: " + err.Error())

		return nil, errors.InternalServerError("get specials operation fail")
	}

	result := []*SpecialRow{}
	for _, row := range resp.Responses[tableName] {
		user := &SpecialRow{}
		if err := user.Unmarshal(row); err != nil {
			slog.Error("unmarhal special fail: " + err.Error())

			return nil, errors.InternalServerError("unmarshal special fail")
		}

		result = append(result, user)
	}

	return result, nil
}
