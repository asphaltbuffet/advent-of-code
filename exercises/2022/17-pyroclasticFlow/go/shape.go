package exercises

// shapes have an offset of x+2 to put them in starting position
var shapes = [][]point{
	{
		// ┌┬┬┬┐
		// └┴┴┴┘
		point{y: 0, x: 2},
		point{y: 0, x: 3},
		point{y: 0, x: 4},
		point{y: 0, x: 5},
	}, {
		//  ┌┐
		// ┌┼┼┐
		// └┼┼┘
		//  └┘
		point{y: 0, x: 3},
		point{y: 1, x: 2},
		point{y: 1, x: 3},
		point{y: 1, x: 4},
		point{y: 2, x: 3},
	}, {
		//   ┌┐
		//   ├┤
		// ┌┬┼┤
		// └┴┴┘
		point{y: 0, x: 2},
		point{y: 0, x: 3},
		point{y: 0, x: 4},
		point{y: 1, x: 4},
		point{y: 2, x: 4},
	}, {
		// ┌┐
		// ├┤
		// ├┤
		// ├┤
		// └┘
		point{y: 0, x: 2},
		point{y: 1, x: 2},
		point{y: 2, x: 2},
		point{y: 3, x: 2},
	}, {
		// ┌┬┐
		// ├┼┤
		// └┴┘
		point{y: 0, x: 2},
		point{y: 0, x: 3},
		point{y: 1, x: 2},
		point{y: 1, x: 3},
	},
}
