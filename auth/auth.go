package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

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
	Aud     string
	Iss     string
	JwksURI string
}

// NewValidationKeyGetter func
func NewValidationKeyGetter(options *ValidationKeyGetterOptions) func(token *jwt.Token) (interface{}, error) {
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

		cert, err := getCert(token, options.JwksURI)
		if err != nil {
			return nil, errors.New("invalid token")
		}

		key, err := jwt.ParseRSAPublicKeyFromPEM(cert)
		if err != nil {
			return nil, errors.New("invalid token")
		}

		return key, nil
	}
}

func getCert(token *jwt.Token, jwksURI string) ([]byte, error) {
	resp, err := http.Get(jwksURI)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var jwks = jwks{}

	err = json.NewDecoder(resp.Body).Decode(&jwks)
	if err != nil {
		return nil, err
	}

	kid, ok := token.Header["kid"]
	if !ok {
		return nil, errors.New("kid not found")
	}

	for _, key := range jwks.Keys {
		if key.Kid == kid {
			cert := fmt.Sprintf("%s\n%s\n%s", certBeg, key.X5c[0], certEnd)
			return []byte(cert), nil
		}
	}

	return nil, errors.New("key not found")
}
