package hamming

import "fmt"

func Distance(a, b string) (int, error) {
	if len(a) != len(b) {
        return 0, fmt.Errorf("unequal string lengths: a=%d, b=%d", len(a), len(b))
    }

    var d int
    for i := 0; i < len(a); i++ {
        if a[i] != b[i] {
            d++
        }
    }

    return d, nil
}
