# [Day 19: Not Enough Minerals](https://adventofcode.com/2022/day/19)

<!-- [Day 19: Not Enough Minerals](19-notEnoughMinerals) -->

aka, I think I found the slowest way to do this...

## Go

```text
2022-19 Not Enough Minerals (Golang)

Running...

Test 1.0: pass in 1.2 s
Test 2.0: pass in 241.3 s
Part 1: 960 in 6.2 s
Part 2: 2040 in 42.6 s
```

### Benchmark Notes

I wondered if using a regex for parsing the input would be better; the output of benchstat follows:

```text
goos: linux
goarch: amd64
pkg: github.com/asphaltbuffet/advent-of-code/exercises/2022/19-notEnoughMinerals/go
cpu: Intel(R) Core(TM) i7-9700K CPU @ 3.60GHz
        │   sscanf    │                regex                  │
        │   sec/op    │    sec/op     vs base                 │
Parse-8   5.455µ ± 3%   19.760µ ± 2%  +262.24% (p=0.000 n=15)

        │   sscanf   │                regex                    │
        │    B/op    │     B/op      vs base                   │
Parse-8   158.0 ± 0%   29114.0 ± 0%  +18326.58% (p=0.000 n=15)

        │   sscanf   │                regex                   │
        │ allocs/op  │  allocs/op    vs base                  │
Parse-8   7.000 ± 0%   101.000 ± 0%  +1342.86% (p=0.000 n=15)
```

It's slower, uses more memory, and allocates that memory more often. Regex is **not** the better solution here.

## Python

```text
    < section intentionally left blank >
```

## 2022 Run Times

![2022 exercise run-time graphs](../run-times.png)
