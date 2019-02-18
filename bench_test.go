package hcler_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/creack/hcler"
	"github.com/pkg/errors"
)

func run(b *testing.B, m interface{}) {
	b.Helper()
	for i := 0; i < b.N; i++ {
		if _, err := hcler.Encode(m); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMapStringKeyEncoder(b *testing.B) {
	b.Run("hcl_map", func(b *testing.B) {
		m := hcler.Map{
			"foo": "bar",
			"hello": hcler.Map{
				"world": hcler.Map{
					"ok": "bye",
				},
			},
		}
		run(b, m)
	})

	b.Run("map_str_iface", func(b *testing.B) {
		m := map[string]interface{}{
			"foo": "bar",
			"hello": map[string]interface{}{
				"world": map[string]interface{}{
					"ok": "bye",
				},
			},
		}
		run(b, m)
	})
}

func BenchmarkMapInterfaceEncoder(b *testing.B) {
	b.Run("hcl_i_map", func(b *testing.B) {
		m := hcler.IMap{
			"foo": "bar",
			"hello": hcler.IMap{
				"world": hcler.IMap{
					"ok": "bye",
				},
			},
		}
		run(b, m)
	})

	b.Run("map_iface_iface", func(b *testing.B) {
		m := map[interface{}]interface{}{
			"foo": "bar",
			"hello": map[interface{}]interface{}{
				"world": map[interface{}]interface{}{
					"ok": "bye",
				},
			},
		}
		run(b, m)
	})
}

// Bellow, alternative map encoders considered. Discarded but kept for reference & bench.

func BenchmarkDiscarded(b *testing.B) {
	b.Run("map_str_slice_prealloc", func(b *testing.B) {
		m := mapStrSlice{
			"foo": "bar",
			"hello": mapStrSlice{
				"world": mapStrSlice{
					"ok": "bye",
				},
			},
		}
		run(b, m)
	})

	b.Run("map_str_slice_noprealloc", func(b *testing.B) {
		m := mapStrSliceNoPrealloc{
			"foo": "bar",
			"hello": mapStrSliceNoPrealloc{
				"world": mapStrSliceNoPrealloc{
					"ok": "bye",
				},
			},
		}
		run(b, m)
	})
}

type mapStrSlice map[string]interface{}

func (m mapStrSlice) EncodeHCL() (string, error) {
	elements := make([]string, 0, len(m))
	for k, v := range m {
		valueString, err := hcler.Encode(v)
		if err != nil {
			return "", errors.Wrap(err, "encode value")
		}
		elements = append(elements, fmt.Sprintf("%s = %s", k, valueString))
	}
	return fmt.Sprintf("{ %s }", strings.Join(elements, ", ")), nil
}

type mapStrSliceNoPrealloc map[string]interface{}

func (m mapStrSliceNoPrealloc) EncodeHCL() (string, error) {
	elements := []string{}
	for k, v := range m {
		valueString, err := hcler.Encode(v)
		if err != nil {
			return "", errors.Wrap(err, "could not marshal value")
		}
		elements = append(elements, fmt.Sprintf("%s = %s", k, valueString))
	}
	return fmt.Sprintf("{ %s }", strings.Join(elements, ", ")), nil
}
