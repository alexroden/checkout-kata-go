package checkout

type Repository interface {
	SessionManager
	Validator

	ScanItem(sku string)
}

type SessionManager interface {
	Session() string
	SetSession(session string)
	StartSession() (string, error)
}

type Validator interface {
	HasError() bool
	Errors(int) error
}
