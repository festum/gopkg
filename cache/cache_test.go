package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	assert := assert.New(t)

	msg := map[string]string{
		"key":   "foo",
		"value": "bar",
	}

	c := cache.NewCache()
	c.Save("foo", msg)
	data, ok := c.Get("foo")
	assert.True(ok)
	assert.Equal(msg, data)
	c.Delete("foo")
	_, ok = c.Get("foo")
	assert.False(ok)
}
