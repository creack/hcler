package hcler_test

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/creack/hcler"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Make sure that the hcl types implements the hcler.Encoder interface.
var (
	_ hcler.Encoder = hcler.Map(nil)
	_ hcler.Encoder = hcler.IMap(nil)
	_ hcler.Encoder = hcler.List(nil)
)

// Make sure the hcl types are compatibhle with native types.
var (
	_ hcler.Map  = map[string]interface{}{}
	_ hcler.IMap = map[interface{}]interface{}{}
	_ hcler.List = []interface{}{}
)

type mapTypes struct {
	hclMap  hcler.Map
	hclIMap hcler.IMap
	stdMap  map[string]interface{}
	stdIMap map[interface{}]interface{}
}

type listTypes struct {
	hclList hcler.List
	stdList []interface{}
}

func assertEncoding(expectOneOf []string, m interface{}) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
		got, err := hcler.Encode(m)
		require.NoError(t, err)
		assert.Subset(t, expectOneOf, []string{got})
	}
}

func assertMapTestCase(t *testing.T, expectOneOf []string, vals mapTypes) {
	t.Helper()
	for name, fct := range map[string]func(t *testing.T){
		"hcl_map":   assertEncoding(expectOneOf, vals.hclMap),
		"hcl_i_map": assertEncoding(expectOneOf, vals.hclIMap),
		"std_map":   assertEncoding(expectOneOf, vals.stdMap),
		"std_i_map": assertEncoding(expectOneOf, vals.stdIMap),
	} {
		t.Run(name, fct)
	}
}

func assertListTestCase(t *testing.T, expectOneOf []string, vals listTypes) {
	t.Helper()
	for name, fct := range map[string]func(t *testing.T){
		"hcl_list": assertEncoding(expectOneOf, vals.hclList),
		"std_list": assertEncoding(expectOneOf, vals.stdList),
	} {
		t.Run(name, fct)
	}
}

func TestEncodeNil(t *testing.T) {
	t.Run("maps", func(t *testing.T) {
		expect := []string{`{}`}
		vals := mapTypes{
			hclMap:  nil,
			hclIMap: nil,
			stdMap:  nil,
			stdIMap: nil,
		}
		assertMapTestCase(t, expect, vals)
	})

	t.Run("lists", func(t *testing.T) {
		expect := []string{`[]`}
		vals := listTypes{
			hclList: nil,
			stdList: nil,
		}
		assertListTestCase(t, expect, vals)
	})
}

func TestEncodeEmpty(t *testing.T) {
	t.Run("maps", func(t *testing.T) {
		expect := []string{`{}`}
		vals := mapTypes{
			hclMap:  hcler.Map{},
			hclIMap: hcler.IMap{},
			stdMap:  map[string]interface{}{},
			stdIMap: map[interface{}]interface{}{},
		}
		assertMapTestCase(t, expect, vals)
	})

	t.Run("lists", func(t *testing.T) {
		expect := []string{`[]`}
		vals := listTypes{
			hclList: hcler.List{},
			stdList: []interface{}{},
		}
		assertListTestCase(t, expect, vals)
	})
}

func TestMapEncodeOneDepth(t *testing.T) {
	t.Run("one_var", func(t *testing.T) {
		expect := []string{`{ foo = "bar" }`}
		vals := mapTypes{
			hclMap:  hcler.Map{"foo": "bar"},
			hclIMap: hcler.IMap{"foo": "bar"},
			stdMap:  map[string]interface{}{"foo": "bar"},
			stdIMap: map[interface{}]interface{}{"foo": "bar"},
		}
		assertMapTestCase(t, expect, vals)
	})

	t.Run("two_vars", func(t *testing.T) {
		expectOneOf := []string{
			`{ foo = "bar", hello = "world" }`,
			`{ hello = "world", foo = "bar" }`,
		}
		vals := mapTypes{
			hclMap:  hcler.Map{"foo": "bar", "hello": "world"},
			hclIMap: hcler.IMap{"foo": "bar", "hello": "world"},
			stdMap:  map[string]interface{}{"foo": "bar", "hello": "world"},
			stdIMap: map[interface{}]interface{}{"foo": "bar", "hello": "world"},
		}
		assertMapTestCase(t, expectOneOf, vals)
	})
}

func TestMapEncodeMultipleDepth(t *testing.T) {
	t.Run("one_var", func(t *testing.T) {
		expect := []string{`{ foo = { bar = { baz = { hello = "world" } } } }`}
		vals := mapTypes{
			hclMap:  hcler.Map{"foo": hcler.Map{"bar": hcler.Map{"baz": hcler.Map{"hello": "world"}}}},
			hclIMap: hcler.IMap{"foo": hcler.IMap{"bar": hcler.IMap{"baz": hcler.IMap{"hello": "world"}}}},
			stdMap:  map[string]interface{}{"foo": map[string]interface{}{"bar": map[string]interface{}{"baz": map[string]interface{}{"hello": "world"}}}},
			stdIMap: map[interface{}]interface{}{
				"foo": map[interface{}]interface{}{
					"bar": map[interface{}]interface{}{
						"baz": map[interface{}]interface{}{
							"hello": "world",
						},
					},
				},
			},
		}
		assertMapTestCase(t, expect, vals)
	})

	t.Run("two_var", func(t *testing.T) {
		expectOneOf := []string{
			`{ foo = { bar = "baz", hello = "world" }, foo2 = { bar = "baz", hello = "world" } }`,
			`{ foo = { bar = "baz", hello = "world" }, foo2 = { hello = "world", bar = "baz" } }`,
			`{ foo = { hello = "world", bar = "baz" }, foo2 = { bar = "baz", hello = "world" } }`,
			`{ foo = { hello = "world", bar = "baz" }, foo2 = { hello = "world", bar = "baz" } }`,
			`{ foo2 = { bar = "baz", hello = "world" }, foo = { bar = "baz", hello = "world" } }`,
			`{ foo2 = { bar = "baz", hello = "world" }, foo = { hello = "world", bar = "baz" } }`,
			`{ foo2 = { hello = "world", bar = "baz" }, foo = { bar = "baz", hello = "world" } }`,
			`{ foo2 = { hello = "world", bar = "baz" }, foo = { hello = "world", bar = "baz" } }`,
		}
		vals := mapTypes{
			hclMap: hcler.Map{
				"foo":  hcler.Map{"bar": "baz", "hello": "world"},
				"foo2": hcler.Map{"bar": "baz", "hello": "world"},
			},
			hclIMap: hcler.IMap{
				"foo":  hcler.IMap{"bar": "baz", "hello": "world"},
				"foo2": hcler.IMap{"bar": "baz", "hello": "world"},
			},
			stdMap: map[string]interface{}{
				"foo":  map[string]interface{}{"bar": "baz", "hello": "world"},
				"foo2": map[string]interface{}{"bar": "baz", "hello": "world"},
			},
			stdIMap: map[interface{}]interface{}{
				"foo":  map[interface{}]interface{}{"bar": "baz", "hello": "world"},
				"foo2": map[interface{}]interface{}{"bar": "baz", "hello": "world"},
			},
		}
		assertMapTestCase(t, expectOneOf, vals)
	})
}

func TestEncodeListOneDepth(t *testing.T) {
	expect := []string{`[ "foo", "bar" ]`}
	vals := listTypes{
		hclList: hcler.List{"foo", "bar"},
		stdList: []interface{}{"foo", "bar"},
	}
	assertListTestCase(t, expect, vals)
}

func TestEncodeListMultipleDepth(t *testing.T) {
	expect := []string{`[ [ "foo", [ "bar", "baz" ] ], [ "hello" ], "world" ]`}
	vals := listTypes{
		hclList: hcler.List{hcler.List{"foo", hcler.List{"bar", "baz"}}, hcler.List{"hello"}, "world"},
		stdList: []interface{}{[]interface{}{"foo", []interface{}{"bar", "baz"}}, []interface{}{"hello"}, "world"},
	}
	assertListTestCase(t, expect, vals)
}

type superkey struct{}

func (sk superkey) String() string { return "fakesuperkey" }

func TestEncodeExoticMapKey(t *testing.T) {
	t.Run("non_alphanum_type", func(t *testing.T) {
		expect := []string{`{ "hello world" = "foo" }`}
		vals := mapTypes{
			hclMap:  hcler.Map{"hello world": "foo"},
			hclIMap: hcler.IMap{"hello world": "foo"},
			stdMap:  map[string]interface{}{"hello world": "foo"},
			stdIMap: map[interface{}]interface{}{"hello world": "foo"},
		}
		assertMapTestCase(t, expect, vals)
	})
	t.Run("custom_key_type", func(t *testing.T) {
		expect := []string{`{ fakesuperkey = "foo" }`}
		vals := mapTypes{
			hclMap:  hcler.Map{"fakesuperkey": "foo"},
			hclIMap: hcler.IMap{superkey{}: "foo"},
			stdMap:  map[string]interface{}{"fakesuperkey": "foo"},
			stdIMap: map[interface{}]interface{}{superkey{}: "foo"},
		}
		assertMapTestCase(t, expect, vals)
	})
	t.Run("url_key", func(t *testing.T) {
		expect := []string{`{ fakesuperkey = "foo" }`}
		vals := mapTypes{
			hclMap:  hcler.Map{"fakesuperkey": "foo"},
			hclIMap: hcler.IMap{superkey{}: "foo"},
			stdMap:  map[string]interface{}{"fakesuperkey": "foo"},
			stdIMap: map[interface{}]interface{}{superkey{}: "foo"},
		}
		assertMapTestCase(t, expect, vals)
	})
}

func TestEncodeMixedTypes(t *testing.T) {
	t.Run("nested_maps", func(t *testing.T) {
		expect := []string{`{ foo = { bar = "baz" } }`}
		vals1 := mapTypes{
			hclMap:  hcler.Map{"foo": hcler.IMap{"bar": "baz"}},
			hclIMap: hcler.IMap{"foo": hcler.Map{"bar": "baz"}},
			stdMap:  map[string]interface{}{"foo": map[interface{}]interface{}{"bar": "baz"}},
			stdIMap: map[interface{}]interface{}{"foo": map[string]interface{}{"bar": "baz"}},
		}
		vals2 := mapTypes{
			hclMap:  hcler.Map{"foo": map[interface{}]interface{}{"bar": "baz"}},
			hclIMap: hcler.IMap{"foo": map[string]interface{}{"bar": "baz"}},
			stdMap:  map[string]interface{}{"foo": hcler.IMap{"bar": "baz"}},
			stdIMap: map[interface{}]interface{}{"foo": hcler.Map{"bar": "baz"}},
		}
		assertMapTestCase(t, expect, vals1)
		assertMapTestCase(t, expect, vals2)
	})

	t.Run("map_lists", func(t *testing.T) {
		expect := []string{`{ foo = [ "bar", "baz" ] }`}
		vals1 := mapTypes{
			hclMap:  hcler.Map{"foo": hcler.List{"bar", "baz"}},
			hclIMap: hcler.IMap{"foo": hcler.List{"bar", "baz"}},
			stdMap:  map[string]interface{}{"foo": []interface{}{"bar", "baz"}},
			stdIMap: map[interface{}]interface{}{"foo": []interface{}{"bar", "baz"}},
		}
		vals2 := mapTypes{
			hclMap:  hcler.Map{"foo": []interface{}{"bar", "baz"}},
			hclIMap: hcler.IMap{"foo": []interface{}{"bar", "baz"}},
			stdMap:  map[string]interface{}{"foo": hcler.List{"bar", "baz"}},
			stdIMap: map[interface{}]interface{}{"foo": hcler.List{"bar", "baz"}},
		}
		assertMapTestCase(t, expect, vals1)
		assertMapTestCase(t, expect, vals2)
	})

	t.Run("nested_lists", func(t *testing.T) {
		expect := []string{`[ "foo", [ "bar", "baz" ], [] ]`}
		vals1 := listTypes{
			hclList: hcler.List{"foo", hcler.List{"bar", "baz"}, hcler.List(nil)},
			stdList: []interface{}{"foo", []interface{}{"bar", "baz"}, []interface{}(nil)},
		}
		vals2 := listTypes{
			hclList: hcler.List{"foo", []interface{}{"bar", "baz"}, []interface{}(nil)},
			stdList: []interface{}{"foo", hcler.List{"bar", "baz"}, hcler.List(nil)},
		}
		assertListTestCase(t, expect, vals1)
		assertListTestCase(t, expect, vals2)
	})

	t.Run("list_maps", func(t *testing.T) {
		expect := []string{`[ "foo", { bar = "baz" }, [] ]`}
		vals1 := listTypes{
			hclList: hcler.List{"foo", hcler.Map{"bar": "baz"}, hcler.List(nil)},
			stdList: []interface{}{"foo", map[string]interface{}{"bar": "baz"}, []interface{}(nil)},
		}
		vals2 := listTypes{
			hclList: hcler.List{"foo", hcler.IMap{"bar": "baz"}, hcler.List(nil)},
			stdList: []interface{}{"foo", map[interface{}]interface{}{"bar": "baz"}, []interface{}(nil)},
		}
		_ = `
"[ "foo", { bar = "baz" }, [] ]" does not contain "[ "foo", "{ bar = baz }", "[]" ]"
`
		assertListTestCase(t, expect, vals1)
		assertListTestCase(t, expect, vals2)
	})

}

func TestEncodeError(t *testing.T) {
	t.Run("hcl_imap_non_string_keys", func(t *testing.T) {
		k := unsafe.Pointer(t)
		val := hcler.IMap{k: "foo"}
		_, err := hcler.Encode(val)
		require.Error(t, err)

		expect := fmt.Sprintf("unsupported type %T", k)
		assert.Equal(t, expect, errors.Cause(err).Error())
	})
	t.Run("std_imap_non_string_keys", func(t *testing.T) {
		k := unsafe.Pointer(t)
		val := map[interface{}]interface{}{k: "foo"}
		_, err := hcler.Encode(val)
		require.Error(t, err)

		expect := fmt.Sprintf("unsupported type %T", k)
		assert.Equal(t, expect, errors.Cause(err).Error())
	})
	t.Run("unsupposed_root_type", func(t *testing.T) {
		r := unsafe.Pointer(t)
		_, err := hcler.Encode(r)
		require.Error(t, err)

		expect := fmt.Sprintf("unsupported type %T", r)
		assert.Equal(t, expect, errors.Cause(err).Error())
	})
	t.Run("unsupposed_map_type", func(t *testing.T) {
		r := unsafe.Pointer(t)
		val := hcler.Map{"foo": r}
		_, err := hcler.Encode(val)
		require.Error(t, err)

		expect := fmt.Sprintf("unsupported type %T", r)
		assert.Equal(t, expect, errors.Cause(err).Error())
	})
	t.Run("unsupposed_list_type", func(t *testing.T) {
		r := unsafe.Pointer(t)
		val := hcler.List{r}
		_, err := hcler.Encode(val)
		require.Error(t, err)

		expect := fmt.Sprintf("unsupported type %T", r)
		assert.Equal(t, expect, errors.Cause(err).Error())
	})
}

func TestIMapConvertion(t *testing.T) {
	t.Run("nil_map", func(t *testing.T) {
		var m1 hcler.IMap
		m2, err := m1.Map()
		require.NoError(t, err)
		assert.Nil(t, m2)
	})
	t.Run("empty_map", func(t *testing.T) {
		m1 := hcler.IMap{}
		m2, err := m1.Map()
		require.NoError(t, err)
		assert.Nil(t, m2)
	})
	t.Run("nested_map", func(t *testing.T) {
		m1 := map[interface{}]interface{}{
			"foo": map[string]interface{}{
				"bar": "baz",
			},
		}
		m2, err := hcler.IMap(m1).Map()
		require.NoError(t, err)

		expect := hcler.Map{"foo": map[string]interface{}{"bar": "baz"}}
		assert.Equal(t, expect, m2)
	})
}
