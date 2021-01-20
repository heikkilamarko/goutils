package bearer

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

const (
	defaultRefreshInterval = 12 * time.Hour
	minimumRefreshInterval = time.Minute
)

// KeyProviderOptions struct
type KeyProviderOptions struct {
	MetadataURI     string
	RefreshInterval time.Duration
	Logger          *zerolog.Logger
}

// KeyProvider struct
type KeyProvider struct {
	options *KeyProviderOptions
	keys    map[string]interface{}
	mu      sync.RWMutex
}

// NewKeyProvider func
func NewKeyProvider(ctx context.Context, options KeyProviderOptions) (*KeyProvider, error) {
	if options.MetadataURI == "" {
		return nil, errors.New("missing metadata uri")
	}

	if options.RefreshInterval < minimumRefreshInterval {
		if options.RefreshInterval == 0 {
			options.RefreshInterval = defaultRefreshInterval
		}
		options.RefreshInterval = minimumRefreshInterval
	}

	p := &KeyProvider{
		options: &options,
		keys:    make(map[string]interface{}),
	}

	if err := p.start(ctx); err != nil {
		return nil, err
	}

	return p, nil
}

// GetKey method
func (p *KeyProvider) GetKey(kid string) (interface{}, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	if key, ok := p.keys[kid]; ok {
		return key, nil
	}

	return nil, errors.New("key not found")
}

// Refresh method
func (p *KeyProvider) Refresh() error {
	p.logInfo("refreshing keys...")

	keys, err := getKeys(p.options.MetadataURI)
	if err != nil {
		p.logError(err)
		return err
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	p.keys = keys

	return nil
}

func (p *KeyProvider) start(ctx context.Context) error {
	if err := p.Refresh(); err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-time.After(p.options.RefreshInterval):
				p.Refresh()
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (p *KeyProvider) logInfo(msg string) {
	if p.options.Logger != nil {
		p.options.Logger.Info().Msg(msg)
	}
}

func (p *KeyProvider) logError(err error) {
	if p.options.Logger != nil {
		p.options.Logger.Err(err).Send()
	}
}
