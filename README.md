# hcler

Encode arbitrary Go types to Hashicorp's HCL format.

[![GitHub release](https://img.shields.io/github/release/creack/hcler/all.svg?maxAge=2592000)]() [![GoDoc](https://godoc.org/github.com/creack/hcler?status.svg)](https://godoc.org/github.com/creack/hcler) [![Build Status](https://travis-ci.org/creack/hcler.svg)](https://travis-ci.org/creack/hcler) [![Coverage Status](https://coveralls.io/repos/github/creack/hcler/badge.svg?branch=master)](https://coveralls.io/github/creack/hcler?branch=master)

## Usage

## Common types

`hcler.Encode()` can be used with most common types.

## Custom types

In order to support custom types, hcler provides the `hcler.Encoder` interface, similar to `json.Marshaler` & co.
The interface only requires the `EncodeHCL() (string, error)` method.

## Benchmark

```
goos: linux
goarch: amd64
pkg: github.com/creack/hcler
BenchmarkMapStringKeyEncoder/hcl_map-32                  1000000              1856 ns/op             334 B/op         14 allocs/op
BenchmarkMapStringKeyEncoder/map_str_iface-32            1000000              1883 ns/op             333 B/op         14 allocs/op
BenchmarkMapInterfaceEncoder/hcl_i_map-32                 300000              4259 ns/op            1346 B/op         20 allocs/op
BenchmarkMapInterfaceEncoder/map_iface_iface-32           300000              4561 ns/op            1346 B/op         20 allocs/op
BenchmarkDiscarded/map_str_slice_prealloc-32              500000              3654 ns/op             512 B/op         24 allocs/op
BenchmarkDiscarded/map_str_slice_noprealloc-32            500000              3759 ns/op             528 B/op         25 allocs/op
PASS
ok      github.com/creack/hcler 10.349s
```
