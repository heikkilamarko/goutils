package bearer

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

type metadata struct {
	JwksURI string `json:"jwks_uri"`
}

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

func getKeys(metadataURI string) (map[string]interface{}, error) {
	metadata, err := getMetadata(metadataURI)
	if err != nil {
		return nil, err
	}

	jwks, err := getJwks(metadata.JwksURI)
	if err != nil {
		return nil, err
	}

	var keys = make(map[string]interface{})

	for _, jwksKey := range jwks.Keys {
		if key, err := getKey(jwksKey.X5c); err == nil {
			keys[jwksKey.Kid] = key
		}
	}

	return keys, nil
}

func getMetadata(url string) (*metadata, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get metadata from url '%s' failed: %w", url, err)
	}

	defer resp.Body.Close()

	md := &metadata{}

	if err = json.NewDecoder(resp.Body).Decode(md); err != nil {
		return nil, fmt.Errorf("parse metadata failed: %w", err)
	}

	return md, nil
}

func getJwks(url string) (*jwks, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get jwks from url '%s' failed: %w", url, err)
	}

	defer resp.Body.Close()

	jwks := &jwks{}

	if err = json.NewDecoder(resp.Body).Decode(jwks); err != nil {
		return nil, fmt.Errorf("parse jwks failed: %w", err)
	}

	return jwks, nil
}

func getKey(x5c []string) (interface{}, error) {
	if len(x5c) == 0 {
		return nil, errors.New("invalid x5c")
	}

	pem := fmt.Sprintf("%s\n%s\n%s", certBeg, x5c[0], certEnd)

	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
	if err != nil {
		return nil, fmt.Errorf("parse rsa public key from pem failed: %w", err)
	}

	return key, nil
}
