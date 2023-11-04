package lasagna

func PreparationTime(layers []string, t int) int {
	if t == 0 {
		t = 2
	}

	return len(layers) * t
}

func Quantities(layers []string) (int, float64) {
	var (
		n int
		s float64
	)

	for _, l := range layers {
		switch l {
		case "noodles":
			n += 50
		case "sauce":
			s += 0.2
		}
	}

	return n, s
}

func AddSecretIngredient(fl []string, ml []string) {
	ml = append(ml[:len(ml)-1], fl[len(fl)-1:]...)
}

func ScaleRecipe(input []float64, p int) []float64 {
	var scaled []float64

	mult := float64(p) / 2.0

	for _, v := range input {
		scaled = append(scaled, v*mult)
	}

	return scaled
}

// Your first steps could be to read through the tasks, and create
// these functions with their correct parameter lists and return types.
// The function body only needs to contain `panic("")`.
//
// This will make the tests compile, but they will fail.
// You can then implement the function logic one by one and see
// an increasing number of tests passing as you implement more
// functionality.
