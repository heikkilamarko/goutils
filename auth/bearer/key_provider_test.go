package bearer

import (
	"testing"

	"github.com/heikkilamarko/goutils/auth"
)

func TestNewKeyProvider(t *testing.T) {
	var _ auth.KeyProvider = &KeyProvider{}
}
