package auth

import (
	"errors"

	"github.com/form3tech-oss/jwt-go"
)

const (
	certBeg = "-----BEGIN CERTIFICATE-----"
	certEnd = "-----END CERTIFICATE-----"
)

type jwks struct {
	Keys []jwksKey `json:"keys"`
}

type jwksKey struct {
	Kty string   `json:"kty"`
	Use string   `json:"use"`
	Kid string   `json:"kid"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

// ValidationKeyGetterOptions struct
type ValidationKeyGetterOptions struct {
	Aud         string
	Iss         string
	KeyProvider *KeyProvider
}

// NewValidationKeyGetter func
func NewValidationKeyGetter(options *ValidationKeyGetterOptions) func(token *jwt.Token) (interface{}, error) {

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

		key, err := options.KeyProvider.GetKey(kid.(string))
		if err != nil {
			return nil, err
		}

		return key, nil
	}
}
