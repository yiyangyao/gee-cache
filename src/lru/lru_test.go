package lru

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type String string

func (s String) Len() int {
	return len(s)
}

func setup() *Cache {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	caps := len(k1 + k2 + k3 + v1 + v2 + v3)
	cache := New(int64(caps), nil)
	cache.Add(k1, String(v1))
	cache.Add(k2, String(v2))
	return cache
}

func TestCache_Add(t *testing.T) {
	cache := setup()
	var cacheAddTests = []struct {
		inputA   string
		inputB   string
		excepted int
	}{
		{
			inputA:   "k3",
			inputB:   "v3",
			excepted: 3,
		},
		{
			inputA:   "key3",
			inputB:   "value3",
			excepted: 2,
		},
	}
	for _, tt := range cacheAddTests {
		_ = cache.Add(tt.inputA, String(tt.inputB))
		actual := cache.Len()
		assert.Equal(t, tt.excepted, actual)
	}
}

func TestCache_Get(t *testing.T) {
	cache := setup()
	var cacheGetTests = []struct {
		input    string
		excepted string
	}{
		{
			input:    "key1",
			excepted: "value1",
		},
		{
			input:    "key3",
			excepted: "",
		},
	}
	for _, tt := range cacheGetTests {
		if v, ok := cache.Get(tt.input); ok {
			assert.Equal(t, tt.excepted, string(v.(String)))
		} else {
			assert.Equal(t, tt.excepted, "")
		}
	}
}
