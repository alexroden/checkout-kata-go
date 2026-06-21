package dynamodb

import (
	"encoding/json"
	"log/slog"

	"github.com/alexroden/checkout-kata-go/pkg/errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func (d *DynamoDB) PutBasketItem(row *BasketItemRow) error {
	tableName, err := d.getTable(BASKET_ITEMS_TABLE)
	if err != nil {
		return err
	}

	item, err := row.Marshal()
	if err != nil {
		slog.Error("data marshal fail: " + err.Error())

		return errors.InternalServerError("unable to marshal data")
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	}

	b, _ := json.Marshal(input)
	slog.Debug(string(b))

	if _, err = d.db.PutItem(d.ctx, input); err != nil {
		slog.Error("put basket item operation fail: " + err.Error())

		return errors.InternalServerError("put basket item operation fail")
	}

	return nil
}
