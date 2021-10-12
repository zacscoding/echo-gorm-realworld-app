package cache

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test2(t *testing.T) {
	cli, closeFn := NewTestCache(t)
	defer closeFn()

	assert.NoError(t, cli.Set(context.Background(), "key1", "value1", time.Second).Err())

	t.Logf(">> Get: %s", cli.Get(context.Background(), "key1").String())
}
