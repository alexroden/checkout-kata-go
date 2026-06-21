package main

import (
	"log/slog"
	"os"

	"github.com/alexroden/checkout-kata-go/pkg/controllers"
	"github.com/alexroden/checkout-kata-go/pkg/dynamodb"
	"github.com/alexroden/checkout-kata-go/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	logger.New(os.Getenv("LOG_LEVEL"))

	r := gin.Default()

	endpoint := os.Getenv("DYNAMO_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8000"
	}

	db, err := dynamodb.New(&dynamodb.Config{
		Region:   os.Getenv("AWS_REGION"),
		Endpoint: endpoint,
		Tables: map[string]string{
			dynamodb.BASKETS_TABLE: os.Getenv("BASKETS_TABLE_NAME"),
		},
	})
	if err != nil {
		slog.Error("dynamodb connection fail: " + err.Error())
	}

	checkoutController := controllers.NewCheckoutController(db)

	v1 := r.Group("/v1")
	{
		v1.POST("/start-session", checkoutController.StartSession)
	}

	r.Run(":8080")
}
