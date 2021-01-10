package oidc

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
	JwksURI         string
	RefreshInterval time.Duration
}

// KeyProvider struct
type KeyProvider struct {
	options  *KeyProviderOptions
	keyCache map[string]interface{}
	sync.RWMutex
}

// NewKeyProvider func
func NewKeyProvider(ctx context.Context, o KeyProviderOptions) (*KeyProvider, error) {
	if o.JwksURI == "" {
		return nil, errors.New("invalid jwks uri")
	}

	if o.RefreshInterval == 0 {
		o.RefreshInterval = time.Hour
	}

	p := &KeyProvider{
		&o,
		make(map[string]interface{}),
		sync.RWMutex{},
	}

	if err := p.start(ctx); err != nil {
		return nil, errors.New("start: metadata not found")
	}

	return p, nil
}

// GetKey method
func (p *KeyProvider) GetKey(kid string) (interface{}, error) {
	p.RLock()
	defer p.RUnlock()

	if key, ok := p.keyCache[kid]; ok {
		return key, nil
	}

	return nil, errors.New("key not found")
}

func (p *KeyProvider) start(ctx context.Context) error {
	if err := p.load(); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-time.After(p.options.RefreshInterval):
				p.load()
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (p *KeyProvider) load() error {
	resp, err := http.Get(p.options.JwksURI)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var jwks jwks

	if err = json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return err
	}

	var keyCache = make(map[string]interface{})

	for _, key := range jwks.Keys {
		if k, err := getKey(key.X5c); err == nil {
			keyCache[key.Kid] = k
		}
	}

	p.Lock()
	defer p.Unlock()

	p.keyCache = keyCache

	return nil
}

func getKey(x5c []string) (interface{}, error) {
	if len(x5c) == 0 {
		return nil, errors.New("empty x5c")
	}

	pem := fmt.Sprintf("%s\n%s\n%s", certBeg, x5c[0], certEnd)

	key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
	if err != nil {
		return nil, err
	}

	return key, nil
}
