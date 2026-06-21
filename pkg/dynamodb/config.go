package dynamodb

type Config struct {
	Region   string
	Endpoint string
	Tables   map[string]string
}
