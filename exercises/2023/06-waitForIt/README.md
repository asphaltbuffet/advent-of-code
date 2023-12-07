# [Day 6: Wait For It](https://adventofcode.com/2023/day/6)

<!-- These are helper text to make formatting the yearly readme consistent and easier...

[Day 6: Wait For It][rm6]
[Go][go6]
[Python][py6]

[rm6]: 06-waitForIt/README.md
[go6]: 06-waitForIt/go
[py6]: 06-waitForIt/py

-->

## Go

```text
< section intentionally left blank >
```

### Benchmark

I initially wrote this a very bad way (yeah, I knew it while it was happening but just kept moving forward). It was fairly easy to optimize, so I kept both versions and used benchstat to compare the outputs of `go.exe test -bench=. -run='^$' -tags test -count=10 -benchmem ./go`:

```text
goos: windows
goarch: amd64
pkg: github.com/asphaltbuffet/advent-of-code/exercises/2023/06-waitForIt/go
cpu: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
                  │     old.bench     │              new.bench               │
                  │      sec/op       │   sec/op     vs base                 │
PartOne/input_0-8        1396.0n ± 2%   890.2n ± 2%   -36.23% (p=0.000 n=10)
PartOne/input_1-8         3.648µ ± 4%   1.127µ ± 2%   -69.12% (p=0.000 n=10)
PartOne/input_2-8         3.611µ ± 2%   1.128µ ± 2%   -68.75% (p=0.000 n=10)
PartTwo/input_0-8      230895.5n ± 9%   507.5n ± 2%   -99.78% (p=0.000 n=10)
PartTwo/input_1-8   110945126.0n ± 7%   676.9n ± 1%  -100.00% (p=0.000 n=10)
PartTwo/input_2-8    75850688.5n ± 1%   674.8n ± 1%  -100.00% (p=0.000 n=10)
geomean                   181.5µ        800.1n        -99.56%

                  │     old.bench     │              new.bench              │
                  │       B/op        │    B/op     vs base                 │
PartOne/input_0-8         2000.0 ± 0%   736.0 ± 0%   -63.20% (p=0.000 n=10)
PartOne/input_1-8         8368.0 ± 0%   912.0 ± 0%   -89.10% (p=0.000 n=10)
PartOne/input_2-8         7920.0 ± 0%   912.0 ± 0%   -88.48% (p=0.000 n=10)
PartTwo/input_0-8      1294514.0 ± 0%   176.0 ± 0%   -99.99% (p=0.000 n=10)
PartTwo/input_1-8   1094566079.5 ± 0%   184.0 ± 0%  -100.00% (p=0.000 n=10)
PartTwo/input_2-8    751993030.0 ± 0%   184.0 ± 0%  -100.00% (p=0.000 n=10)
geomean                  704.7Ki        392.3        -99.95%

                  │ old.bench  │             new.bench              │
                  │ allocs/op  │ allocs/op   vs base                │
PartOne/input_0-8   22.00 ± 0%   16.00 ± 0%  -27.27% (p=0.000 n=10)
PartOne/input_1-8   27.00 ± 0%   19.00 ± 0%  -29.63% (p=0.000 n=10)
PartOne/input_2-8   27.00 ± 0%   19.00 ± 0%  -29.63% (p=0.000 n=10)
PartTwo/input_0-8   9.000 ± 0%   7.000 ± 0%  -22.22% (p=0.000 n=10)
PartTwo/input_1-8   9.000 ± 0%   7.000 ± 0%  -22.22% (p=0.000 n=10)
PartTwo/input_2-8   9.000 ± 0%   7.000 ± 0%  -22.22% (p=0.000 n=10)
geomean             15.07        11.21       -25.61%
```

## Python

```text
< section intentionally left blank >
```

## 2023 Run Times

![2023 exercise run-time graphs](../run-times.png)
