package auth

import (
	"errors"

	"github.com/form3tech-oss/jwt-go"
)

// KeyGetterOptions struct
type KeyGetterOptions struct {
	Aud         string
	Iss         string
	KeyProvider KeyProvider
}

// NewKeyGetter func
func NewKeyGetter(options *KeyGetterOptions) func(token *jwt.Token) (interface{}, error) {

	if options == nil || options.KeyProvider == nil || options.Aud == "" || options.Iss == "" {
		panic("invalid options")
	}

	return func(token *jwt.Token) (interface{}, error) {
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, errors.New("invalid claims")
		}

		ok = claims.VerifyAudience(options.Aud, true)
		if !ok {
			return nil, errors.New("invalid audience")
		}

		ok = claims.VerifyIssuer(options.Iss, true)
		if !ok {
			return nil, errors.New("invalid issuer")
		}

		kid, ok := token.Header["kid"]
		if !ok {
			return nil, errors.New("invalid kid")
		}

		key, err := options.KeyProvider.GetKey(kid)
		if err != nil {
			return nil, err
		}

		return key, nil
	}
}