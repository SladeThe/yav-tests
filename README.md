# Yet Another Validator's benchmarks and tests

This repo benchmarks and tests [SladeThe/yav](https://github.com/SladeThe/yav) against other validators.

## Benchmarks

Valid Account validation, measured on 2026-09-18 with Go 1.27.1 on Windows/amd64,
AMD Ryzen 9 9900X, `GOAMD64=v1`, and `GOMAXPROCS=24`.
Versions: YAV `bdf780b`, go-playground/validator v10.15.1, and ozzo-validation v4.4.1.

Each operation validates one account. Results are medians of five 300 ms samples collected in one process.

### Sequential

| Validator | ns/op | B/op | allocs/op |
|---|---:|---:|---:|
| YAV | 438.7 | 0 | 0 |
| Preallocated YAV | 375.7 | 0 | 0 |
| Playground v10.15.1 | 3078 | 559 | 15 |
| Ozzo v4.4.1 | 6324 | 6255 | 107 |

### Parallel

| Validator | ns/op | B/op | allocs/op |
|---|---:|---:|---:|
| YAV | 32.74 | 0 | 0 |
| Preallocated YAV | 29.44 | 0 | 0 |
| Playground v10.15.1 | 509.1 | 552 | 15 |
| Ozzo v4.4.1 | 2042 | 6211 | 107 |

Parallel ns/op reflects aggregate throughput across 24 workers.
Preallocated YAV reuses rules prepared before validation.

All 32 compatibility tests pass against Playground v10.15.1.
The Ozzo fixture omits some password/name rules and avatar key/value checks.

### Reproduce

Run with Go 1.27.1 from this repository:

```text
go test -mod=readonly -count=1 ./...
go test -mod=readonly -run '^$' -bench . -benchmem -benchtime=300ms -count=5 '-cpu=1,24'
```

The tables use the 24-CPU samples. The quoted CPU list works in PowerShell.
