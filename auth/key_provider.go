package auth

// KeyProvider interface
type KeyProvider interface {
	GetKey(kid interface{}) (interface{}, error)
}
