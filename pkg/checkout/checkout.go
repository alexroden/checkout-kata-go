package checkout

import "github.com/alexroden/checkout-kata-go/pkg/dynamodb"

var isTesting bool = false

type Checkout struct {
	db      dynamodb.Repository
	session string
	errors  []error
}

func New(db dynamodb.Repository) Repository {
	return &Checkout{db: db}
}

func (c *Checkout) Session() string {
	return c.session
}

func (c *Checkout) SetSession(session string) {
	c.session = session
}

func (c *Checkout) HasError() bool {
	return len(c.errors) > 0
}

func (c *Checkout) Errors(index int) error {
	if index < 0 || index >= len(c.errors) {
		return nil
	}

	return c.errors[index]
}
