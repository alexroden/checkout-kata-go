package dynamodb

import (
	"context"
	"fmt"

	"github.com/alexroden/checkout-kata-go/pkg/errors"
	"github.com/alexroden/checkout-kata-go/pkg/repositories"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

const (
	PRODUCTS_TABLE     = "products"
	SPECIALS_TABLE     = "specials"
	BASKETS_TABLE      = "baskets"
	BASKET_ITEMS_TABLE = "basketItems"
)

var loadDefaultConfig = config.LoadDefaultConfig

type DynamoDB struct {
	db  repositories.DynamoDBAPI
	cnf *Config
	ctx context.Context
}

func New(cnf *Config) (Repository, error) {
	ctx := context.Background()

	var cnfOpts []func(*config.LoadOptions) error
	if cnf.Region != "" {
		cnfOpts = append(cnfOpts, config.WithRegion(cnf.Region))
	}

	cfg, err := loadDefaultConfig(ctx, cnfOpts...)
	if err != nil {
		return nil, err
	}

	result := &DynamoDB{
		cnf: cnf,
		ctx: ctx,
	}

	var clientOpts []func(*dynamodb.Options)
	if cnf.Endpoint != "" {
		clientOpts = append(clientOpts, func(o *dynamodb.Options) {
			o.BaseEndpoint = aws.String(cnf.Endpoint)
		})
	}

	result.SetDb(dynamodb.NewFromConfig(cfg, clientOpts...))

	return result, nil
}

func (d *DynamoDB) getTable(table string) (string, error) {
	if d.cnf != nil {
		if val, ok := d.cnf.Tables[table]; ok {
			return val, nil
		}
	}

	return "", errors.BadGateway(fmt.Sprintf("%s table not set", table))
}

func (d *DynamoDB) SetDb(db repositories.DynamoDBAPI) {
	d.db = db
}
