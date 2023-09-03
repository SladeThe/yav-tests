# Yet Another Validator's benchmarks and tests against other validators

See [SladeThe/yav](https://github.com/SladeThe/yav).

```
goos: windows
goarch: amd64
cpu: Intel(R) Core(TM) i9-10850K CPU @ 3.60GHz
```

#### Tiny struct validation

```
BenchmarkYAV              12907930       92.15 ns/op          0 B/op        0 allocs/op
BenchmarkOzzo              1334562       890.1 ns/op       1248 B/op       20 allocs/op
BenchmarkPlayground        1324868       911.8 ns/op         40 B/op        2 allocs/op
```

#### Account struct validation

```
BenchmarkYAV                729123        1658 ns/op        123 B/op        4 allocs/op
BenchmarkOzzo*               54954       21684 ns/op      19215 B/op      317 allocs/op
BenchmarkPlayground         172633        6789 ns/op        653 B/op       23 allocs/op
```
