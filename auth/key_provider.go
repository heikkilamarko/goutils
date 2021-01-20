package auth

// KeyProvider interface
type KeyProvider interface {
	GetKey(string) (interface{}, error)
	Refresh() error
}
