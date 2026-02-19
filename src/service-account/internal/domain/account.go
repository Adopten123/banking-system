package domain

type Account struct {
	ID      int
	Balance float64
}

type AccountRepository interface {
	Ping() error
}

type AccountService interface {
	CheckHealth() string
}
