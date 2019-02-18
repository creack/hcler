package hcler

import (
	"fmt"
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

// EncodeHCL implements the hcl.Encoder interface.
func (m Map) EncodeHCL() (string, error) {
	if len(m) == 0 {
		return "{}", nil
	}
	var b strings.Builder

	b.WriteString("{ ")
	for k, v := range m {
		valueString, err := Encode(v)
		if err != nil {
			return "", errors.Wrapf(err, "encode %q", k)
		}
		// Quotes mandatory with special chars, so always set them.
		b.WriteString(`"`)
		b.WriteString(k)
		b.WriteString(`" = `)
		b.WriteString(valueString)
		b.WriteString(", ")
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
		s, err := toString(k)
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
func (l List) EncodeHCL() (string, error) {
	if len(l) == 0 {
		return "[]", nil
	}
	var b strings.Builder

	b.WriteString("[ ")
	for _, v := range l {
		valueString, err := Encode(v)
		if err != nil {
			return "", errors.Wrap(err, "encode list element")
		}
		b.WriteString(valueString)
		b.WriteString(", ")
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
		s, err := toString(v)
		if err != nil {
			return "", err
		}
		return `"` + s + `"`, nil
	}
}

// toString tries to convert the value to string.
// If nil, returns "".
func toString(v interface{}) (string, error) {
	if v == nil {
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
	case int8:
		out = strconv.FormatInt(int64(v), 10)
	case int16:
		out = strconv.FormatInt(int64(v), 10)
	case int64:
		out = strconv.FormatInt(v, 10)
	case uint:
		out = strconv.FormatUint(uint64(v), 10)
	case uint16:
		out = strconv.FormatUint(uint64(v), 10)
	case uint32:
		out = strconv.FormatUint(uint64(v), 10)
	case uint64:
		out = strconv.FormatUint(v, 10)
	case uintptr:
		out = strconv.FormatUint(uint64(v), 10)
	case float32:
		out = strconv.FormatFloat(float64(v), 'f', 2, 64)
	case float64:
		out = strconv.FormatFloat(v, 'f', 2, 64)
	default:
		return "", errors.Errorf("unsupported type %T", v)
	}
	return out, nil
}
