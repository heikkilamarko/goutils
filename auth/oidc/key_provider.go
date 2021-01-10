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
func NewKeyProvider(ctx context.Context, o KeyProviderOptions) *KeyProvider {
	if o.JwksURI == "" {
		panic("invalid options")
	}

	if o.RefreshInterval == 0 {
		o.RefreshInterval = time.Hour
	}

	p := &KeyProvider{&o, nil, nil, sync.RWMutex{}}

	p.start(ctx)

	return p
}

// GetKey method
func (p *KeyProvider) GetKey(kid interface{}) (interface{}, error) {
	p.RLock()
	defer p.RUnlock()

	key, ok := p.keyCache.Load(kid)

	if ok {
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
	go func() {
		p.load()
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

func (p *KeyProvider) load() {
	resp, err := http.Get(p.options.JwksURI)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	defer resp.Body.Close()

	var jwks = &jwks{}

	err = json.NewDecoder(resp.Body).Decode(jwks)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	p.Lock()
	defer p.Unlock()

	p.keyCache = &sync.Map{}
	p.jwks = jwks
}
