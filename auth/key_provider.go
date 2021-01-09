package auth

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

type cache struct {
	jwks   *jwks
	keyMap sync.Map
}

func newCache(jwks *jwks) *cache {
	return &cache{jwks, sync.Map{}}
}

// KeyProvider struct
type KeyProvider struct {
	JwksURI         string
	RefreshInterval time.Duration
	cache           *cache
}

// Start method
func (p *KeyProvider) Start(ctx context.Context) {
	go func() {
		p.Load()
		for {
			select {
			case <-time.After(p.RefreshInterval):
				p.Load()
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Load method
func (p *KeyProvider) Load() {
	resp, err := http.Get(p.JwksURI)
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

	p.cache = newCache(jwks)
}

// GetKey method
func (p *KeyProvider) GetKey(kid string) (interface{}, error) {

	cache := p.cache

	if cache == nil {
		return nil, errors.New("internal error")
	}

	key, ok := cache.keyMap.Load(kid)

	if ok {
		return key, nil
	}

	var cert []byte

	for _, key := range cache.jwks.Keys {
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

	cache.keyMap.Store(kid, key)

	return key, nil
}

// NewKeyProvider func
func NewKeyProvider(jwksURI string, refreshInterval time.Duration) *KeyProvider {
	return &KeyProvider{jwksURI, refreshInterval, nil}
}
