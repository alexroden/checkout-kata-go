package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/alexroden/checkout-kata-go/internal/dynamo-init/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBAPI interface {
	DescribeTable(ctx context.Context, input *dynamodb.DescribeTableInput, opts ...func(*dynamodb.Options)) (*dynamodb.DescribeTableOutput, error)
	CreateTable(ctx context.Context, input *dynamodb.CreateTableInput, opts ...func(*dynamodb.Options)) (*dynamodb.CreateTableOutput, error)
	BatchWriteItem(ctx context.Context, params *dynamodb.BatchWriteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.BatchWriteItemOutput, error)
}

func main() {
	ctx := context.Background()

	endpoint := os.Getenv("DYNAMO_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8000"
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	fmt.Println("Connected to DynamoDB: 🥞", endpoint)

	tables, err := loadTables("./internal/dynamo-init/tables")
	if err != nil {
		log.Fatal(err)
	}

	for _, table := range tables {
		fmt.Printf("Create table: %s\n", table.TableName)
		if err := createTable(ctx, client, table); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("DynamoDB initialization complete ✅")

	fmt.Println("Seeding data... 🌱")
	files, _ := os.ReadDir("./internal/dynamo-init/seed")
	for _, file := range files {
		tableName := strings.TrimSuffix(file.Name(), ".json")

		runes := []rune(tableName)
		runes[0] = unicode.ToUpper(runes[0])

		tableName = fmt.Sprintf("%sTable", string(runes))

		fmt.Printf("Seeding table: %s\n", tableName)
		err := seedTable(ctx, client, tableName, "./internal/dynamo-init/seed/"+file.Name())
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Seeding complete ✅")
}

func loadTables(dir string) ([]*models.DynamoTable, error) {
	var tables []*models.DynamoTable

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if !strings.HasSuffix(f.Name(), ".json") {
			continue
		}

		path := filepath.Join(dir, f.Name())

		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		var table *models.DynamoTable
		if err := json.Unmarshal(data, &table); err != nil {
			return nil, fmt.Errorf("failed parsing %s: %w", f.Name(), err)
		}

		tables = append(tables, table)
	}

	return tables, nil
}

func createTable(
	ctx context.Context,
	client DynamoDBAPI,
	table *models.DynamoTable,
) error {
	if _, err := client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: aws.String(table.TableName),
	}); err == nil {
		return fmt.Errorf("table %s already exists", table.TableName)
	}

	input := &dynamodb.CreateTableInput{
		TableName: aws.String(table.TableName),
	}

	if table.BillingMode != "" {
		input.BillingMode = table.BillingMode
	}

	for _, attr := range table.Attributes {
		input.AttributeDefinitions = append(input.AttributeDefinitions, types.AttributeDefinition{
			AttributeName: aws.String(attr.Name),
			AttributeType: attr.Type,
		})
	}

	for _, key := range table.KeySchema {
		input.KeySchema = append(input.KeySchema, types.KeySchemaElement{
			AttributeName: aws.String(key.Name),
			KeyType:       key.KeyType,
		})
	}

	if _, err := client.CreateTable(ctx, input); err != nil {
		return fmt.Errorf("failed to create table %s: %w", table.TableName, err)
	}

	return nil
}

func seedTable(
	ctx context.Context,
	client DynamoDBAPI,
	tableName, filePath string,
) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	var items []map[string]interface{}
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}

	const batchSize = 25

	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}

		writeRequests := make([]types.WriteRequest, 0, end-i)

		for _, item := range items[i:end] {
			av, err := attributevalue.MarshalMap(item)
			if err != nil {
				return err
			}

			writeRequests = append(writeRequests, types.WriteRequest{
				PutRequest: &types.PutRequest{
					Item: av,
				},
			})
		}

		_, err := client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				tableName: writeRequests,
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}
