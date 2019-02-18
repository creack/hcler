package hcler

import (
	"errors"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertString(t *testing.T, expect string, m interface{}) {
	t.Helper()

	got, err := toString(m)
	require.NoError(t, err)
	assert.Equal(t, expect, got)
}

func TestStringConvertion(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		assertString(t, "", nil)
		assertString(t, "", "")
		assertString(t, "hello", "hello")
		assertString(t, "hello", []byte("hello"))
		assertString(t, "hello", []rune("hello"))
		assertString(t, "h", rune('h'))
		assertString(t, "h", byte('h'))
		assertString(t, "h", 'h')
		assertString(t, "error", errors.New("error"))

		rawURL := "https://user:password@domain.tld/path/subpath?foo=bar#foo=bar"
		u, _ := url.Parse(rawURL)
		assertString(t, rawURL, u)
	})

	t.Run("numbers", func(t *testing.T) {
		assertString(t, "42", int(42))
		assertString(t, "42", int8(42))
		assertString(t, "42", int16(42))
		assertString(t, "42", int64(42))
		assertString(t, "42", uint(42))
		assertString(t, "42", uint16(42))
		assertString(t, "42", uint32(42))
		assertString(t, "42", uint64(42))
		assertString(t, "42", uintptr(42))
		assertString(t, "42.00", float32(42))
		assertString(t, "42.00", float64(42))
		assertString(t, "0", false)
		assertString(t, "1", true)
	})
}
