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
	jwks     *jwks
	keyCache *sync.Map
	sync.RWMutex
}

// NewKeyProvider func
func NewKeyProvider(ctx context.Context, o KeyProviderOptions) (*KeyProvider, error) {
	if o.JwksURI == "" {
		return nil, errors.New("invalid options")
	}

	if o.RefreshInterval == 0 {
		o.RefreshInterval = time.Hour
	}

	p := &KeyProvider{&o, nil, nil, sync.RWMutex{}}

	if err := p.start(ctx); err != nil {
		return nil, errors.New("provider start failed")
	}

	return p, nil
}

// GetKey method
func (p *KeyProvider) GetKey(kid interface{}) (interface{}, error) {
	p.RLock()
	defer p.RUnlock()

	if key, ok := p.keyCache.Load(kid); ok {
		return key, nil
	}

	var cert []byte

	for _, key := range p.jwks.Keys {
		if key.Kid == kid {
			cert = []byte(fmt.Sprintf("%s\n%s\n%s", certBeg, key.X5c[0], certEnd))
			break
		}
	}

	if cert == nil {
		return nil, errors.New("key not found")
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(cert)
	if err != nil {
		return nil, errors.New("key not found")
	}

	p.keyCache.Store(kid, key)

	return key, nil
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

	var jwks = &jwks{}

	if err = json.NewDecoder(resp.Body).Decode(jwks); err != nil {
		return err
	}

	p.Lock()
	defer p.Unlock()

	p.keyCache = &sync.Map{}
	p.jwks = jwks

	return nil
}
