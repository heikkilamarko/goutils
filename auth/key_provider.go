package auth

// KeyProvider interface
type KeyProvider interface {
	GetKey(kid string) (interface{}, error)
}
