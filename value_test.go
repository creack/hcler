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

	got, err := toString(m, false)
	require.NoError(t, err)
	assert.Equal(t, expect, got)
}

func assertEscapeString(t *testing.T, expect string, m interface{}) {
	t.Helper()

	got, err := toString(m, true)
	require.NoError(t, err)
	assert.Equal(t, expect, got)
}

func TestStringConvertionRaw(t *testing.T) {
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
		assertString(t, "42", float32(42))
		assertString(t, "42", float64(42))
		assertString(t, "42.01", float32(42.01))
		assertString(t, "42.01", float64(42.01))
		assertString(t, "0", false)
		assertString(t, "1", true)
	})
}

func TestStringConvertionEscape(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		assertEscapeString(t, `""`, nil)
		assertEscapeString(t, `""`, "")
		assertEscapeString(t, `"hello"`, "hello")
		assertEscapeString(t, `"hello"`, []byte("hello"))
		assertEscapeString(t, `"hello"`, []rune("hello"))
		assertEscapeString(t, `"h"`, rune('h'))
		assertEscapeString(t, `"h"`, byte('h'))
		assertEscapeString(t, `"h"`, 'h')
		assertEscapeString(t, `"error"`, errors.New("error"))

		rawURL := "https://user:password@domain.tld/path/subpath?foo=bar#foo=bar"
		u, _ := url.Parse(rawURL)
		assertEscapeString(t, `"`+rawURL+`"`, u)
	})

	t.Run("numbers", func(t *testing.T) {
		assertEscapeString(t, "42", int(42))
		assertEscapeString(t, "42", int8(42))
		assertEscapeString(t, "42", int16(42))
		assertEscapeString(t, "42", int64(42))
		assertEscapeString(t, "42", uint(42))
		assertEscapeString(t, "42", uint16(42))
		assertEscapeString(t, "42", uint32(42))
		assertEscapeString(t, "42", uint64(42))
		assertEscapeString(t, "42", uintptr(42))
		assertEscapeString(t, "42", float32(42))
		assertEscapeString(t, "42", float64(42))
		assertEscapeString(t, "42.01", float32(42.01))
		assertEscapeString(t, "42.01", float64(42.01))

		assertEscapeString(t, `"0"`, false)
		assertEscapeString(t, `"1"`, true)
	})
}
