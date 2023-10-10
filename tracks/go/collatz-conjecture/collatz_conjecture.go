package collatzconjecture

import "fmt"

func CollatzConjecture(n int) (int, error) {
    var (
        c int
        err error
    )
    
    if n <= 0 {
        return 0, fmt.Errorf("negative number: %d", n)
    }

	switch {
        case n == 1:
    		return 0, nil
        case n % 2 == 0:
    		c, err = CollatzConjecture(n/2)
        default:
    		c, err = CollatzConjecture(3 * n + 1)
    }

    return c+1, err
}
