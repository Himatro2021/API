package db

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestCache_Get(t *testing.T) {
	mr, err := miniredis.Run()
	assert.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	cacher := NewCacher(client)
	ctx := context.TODO()

	t.Run("ok", func(t *testing.T) {
		err := mr.Set("test", "ok")
		assert.NoError(t, err)

		res, err := cacher.Get(ctx, "test")

		assert.NoError(t, err)
		assert.Equal(t, res, "ok")
	})

	t.Run("ok - nil response", func(t *testing.T) {
		_, err := cacher.Get(ctx, "keynotexist")

		assert.Error(t, err)
		assert.Equal(t, err, redis.Nil)
	})

	// close miniredis to simulate server eror
	mr.Close()

	t.Run("err - server closed", func(t *testing.T) {
		_, err := cacher.Get(ctx, "oops-shutdown")

		assert.Error(t, err)
	})

}

func TestCache_Set(t *testing.T) {
	mr, err := miniredis.Run()
	assert.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	cacher := NewCacher(client)
	ctx := context.TODO()

	t.Run("ok", func(t *testing.T) {
		err = cacher.Set(ctx, "test", "ok", time.Second*100)
		assert.NoError(t, err)
	})

	// close redis and simulate error condition
	mr.Close()

	t.Run("err - server unavailable", func(t *testing.T) {
		err := cacher.Set(ctx, "test", "ting", time.Second*100)

		assert.Error(t, err)
	})
}
