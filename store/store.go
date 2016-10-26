package store

type Keyer interface {
	Key() string
}

type Store interface {
	LogHistory(Keyer) error
	GetHistory(string) ([]byte, error)
	ListHistory(string, ...int) ([]byte, error)
}
