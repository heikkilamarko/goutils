package bearer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/form3tech-oss/jwt-go"
)

const (
	certBeg = "-----BEGIN CERTIFICATE-----"
	certEnd = "-----END CERTIFICATE-----"
)

const (
	defaultAutomaticRefreshInterval = 24 * time.Hour
	defaultRefreshInterval          = 30 * time.Second
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

// KeyProviderOptions struct
type KeyProviderOptions struct {
	MetadataURI              string
	AutomaticRefreshInterval time.Duration
	RefreshInterval          time.Duration
}

// KeyProvider struct
type KeyProvider struct {
	options     *KeyProviderOptions
	syncAfter   time.Time
	lastRefresh time.Time
	m           sync.Mutex
	keys        map[string]interface{}
}

// NewKeyProvider func
func NewKeyProvider(ctx context.Context, options KeyProviderOptions) (*KeyProvider, error) {
	if options.MetadataURI == "" {
		return nil, errors.New("invalid metadata uri")
	}

	if options.AutomaticRefreshInterval == 0 {
		options.AutomaticRefreshInterval = defaultAutomaticRefreshInterval
	}

	if options.RefreshInterval == 0 {
		options.RefreshInterval = defaultRefreshInterval
	}

	return &KeyProvider{options: &options}, nil
}

// GetKey method
func (p *KeyProvider) GetKey(kid string) (interface{}, error) {
	now := time.Now()

	if p.keys != nil && p.syncAfter.After(now) {
		if key, ok := p.keys[kid]; ok {
			return key, nil
		}

		p.syncAfter = now.Add(p.getSmallerInterval())

		return nil, errors.New("unable to obtain key")
	}

	p.m.Lock()
	defer p.m.Unlock()

	if p.syncAfter.Before(now) || p.syncAfter.Equal(now) {
		keys, err := getKeys(p.options.MetadataURI)
		if err != nil {
			p.syncAfter = now.Add(p.getSmallerInterval())

			if p.keys == nil {
				return nil, errors.New("unable to obtain keys")
			}
		} else {
			p.keys = keys
			p.lastRefresh = now
			p.syncAfter = now.Add(p.options.AutomaticRefreshInterval)
		}
	}

	if p.keys != nil {
		if key, ok := p.keys[kid]; ok {
			return key, nil
		}
	}

	return nil, errors.New("unable to obtain key")
}

// RequestRefresh method
func (p *KeyProvider) RequestRefresh() {
	now := time.Now()
	if now.After(p.lastRefresh.Add(p.options.RefreshInterval)) {
		p.syncAfter = now
	}
}

func (p *KeyProvider) getSmallerInterval() time.Duration {
	if p.options.AutomaticRefreshInterval < p.options.RefreshInterval {
		return p.options.AutomaticRefreshInterval
	}
	return p.options.RefreshInterval
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
