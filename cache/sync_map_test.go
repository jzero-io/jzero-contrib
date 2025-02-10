package cache

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestSyncMap(t *testing.T) {
	t.Run("TestSyncMap", func(t *testing.T) {
		cache := NewSyncMap(errors.New("not found"))

		err := cache.SetCtx(context.Background(), "JWT_ADMIN_AUTH:1:abc", "abc")
		assert.NoError(t, err)
		err = cache.SetCtx(context.Background(), "JWT_ADMIN_AUTH:1:def", "def")
		assert.NoError(t, err)
		err = cache.SetCtx(context.Background(), "JWT_ADMIN_AUTH:1:ghi", "ghi")
		assert.NoError(t, err)

		keys, err := cache.GetPrefixKeysCtx(context.Background(), "JWT_ADMIN_AUTH:1:")
		assert.NoError(t, err)
		assert.Equal(t, []string{"JWT_ADMIN_AUTH:1:abc", "JWT_ADMIN_AUTH:1:def", "JWT_ADMIN_AUTH:1:ghi"}, keys)
	})
}
