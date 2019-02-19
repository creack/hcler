package hcler

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Encoder is the interface implemented by types that
// can marshal themselves into valid HCL.
type Encoder interface {
	EncodeHCL() (string, error)
}

// Map .
type Map map[string]interface{}

var re = regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

// escapeKey look if the given key contains something
// else than alphnum _ -, and wraps it in double quote if.
func escapeKey(k string) string {
	if !re.MatchString(k) {
		return `"` + k + `"`
	}
	return k
}

// EncodeHCL implements the hcl.Encoder interface.
// nolint: gosec
func (m Map) EncodeHCL() (string, error) {
	if len(m) == 0 {
		return "{}", nil
	}
	var b strings.Builder

	// Can't fail beside out of memory error.
	_, _ = b.WriteString("{ ")
	for k, v := range m {
		valueString, err := Encode(v)
		if err != nil {
			return "", errors.Wrapf(err, "encode %q", k)
		}
		// Quotes mandatory with special chars, so always set them.
		_, _ = b.WriteString(escapeKey(k))
		_, _ = b.WriteString(" = ")
		_, _ = b.WriteString(valueString)
		_, _ = b.WriteString(", ")
	}
	return strings.TrimSuffix(b.String(), ", ") + " }", nil
}

// IMap .
type IMap map[interface{}]interface{}

// Map converts an hcl.IMap to a hcl.Map.
func (m IMap) Map() (Map, error) {
	if len(m) == 0 {
		return nil, nil
	}
	out := make(Map, len(m))
	for k, v := range m {
		s, err := toString(k, false)
		if err != nil {
			return nil, errors.Wrap(err, "toString map key")
		}
		out[s] = v
	}
	return out, nil
}

// EncodeHCL implements the hcl.Encoder interface.
// Stringifies the keys and defers to hcl.Map for encoding.
func (m IMap) EncodeHCL() (string, error) {
	if len(m) == 0 {
		return "{}", nil
	}
	out, err := m.Map()
	if err != nil {
		return "", errors.Wrap(err, "convert to hcl.Map")
	}
	return out.EncodeHCL()
}

// List .
type List []interface{}

// EncodeHCL implements the hcl.Encoder interface.
// nolint: gosec
func (l List) EncodeHCL() (string, error) {
	if len(l) == 0 {
		return "[]", nil
	}
	var b strings.Builder

	// Can't fail beside out of memory error.
	_, _ = b.WriteString("[ ")
	for _, v := range l {
		valueString, err := Encode(v)
		if err != nil {
			return "", errors.Wrap(err, "encode list element")
		}
		_, _ = b.WriteString(valueString)
		_, _ = b.WriteString(", ")
	}
	return strings.TrimSuffix(b.String(), ", ") + " ]", nil
}

// Encode implements the hcl.Encoder interface.
func Encode(v interface{}) (string, error) {
	switch v := v.(type) {
	case Encoder:
		return v.EncodeHCL()
	case map[string]interface{}:
		return Map(v).EncodeHCL()
	case map[interface{}]interface{}:
		return IMap(v).EncodeHCL()
	case []interface{}:
		return List(v).EncodeHCL()
	default:
		s, err := toString(v, true)
		if err != nil {
			return "", err
		}
		return s, nil
	}
}

// toString tries to convert the value to string.
// If nil, returns "".
// nolint: gocyclo
func toString(v interface{}, escape bool) (string, error) {
	if v == nil {
		if escape {
			return `""`, nil
		}
		return "", nil
	}
	var out string
	switch v := v.(type) {
	case string:
		out = v
	case rune: // aka int32.
		out = string(v)
	case []rune:
		out = string(v)
	case byte: // aka uint8
		out = string(v)
	case []byte:
		out = string(v)
	case fmt.Stringer:
		out = v.String()
	case error:
		out = v.Error()
	case bool:
		if v {
			out = "1"
		} else {
			out = "0"
		}
	case int:
		out = strconv.FormatInt(int64(v), 10)
		escape = false
	case int8:
		out = strconv.FormatInt(int64(v), 10)
		escape = false
	case int16:
		out = strconv.FormatInt(int64(v), 10)
		escape = false
	case int64:
		out = strconv.FormatInt(v, 10)
		escape = false
	case uint:
		out = strconv.FormatUint(uint64(v), 10)
		escape = false
	case uint16:
		out = strconv.FormatUint(uint64(v), 10)
		escape = false
	case uint32:
		out = strconv.FormatUint(uint64(v), 10)
		escape = false
	case uint64:
		out = strconv.FormatUint(v, 10)
		escape = false
	case uintptr:
		out = strconv.FormatUint(uint64(v), 10)
		escape = false
	case float32:
		if v1 := int64(v); float32(v1) == v {
			return toString(v1, escape)
		}
		out = strconv.FormatFloat(float64(v), 'f', 2, 64)
		escape = false
	case float64:
		if v1 := int64(v); float64(v1) == v {
			return toString(v1, escape)
		}
		out = strconv.FormatFloat(v, 'f', 2, 64)
		escape = false
	default:
		return "", errors.Errorf("unsupported type %T", v)
	}
	if !escape {
		return out, nil
	}
	return `"` + out + `"`, nil
}
