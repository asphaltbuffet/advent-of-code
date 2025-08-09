package chessboard

// Declare a type named File which stores if a square is occupied by a piece - this will be a slice of bools
type File []bool

// Declare a type named Chessboard which contains a map of eight Files, accessed with keys from "A" to "H"
type Chessboard map[string]File

// CountInFile returns how many squares are occupied in the chessboard,
// within the given file.
func CountInFile(cb Chessboard, file string) int {
	var count int

	for _, sq := range cb[file] {
		if sq {
			count++
		}
	}

	return count
}

// CountInRank returns how many squares are occupied in the chessboard,
// within the given rank.
func CountInRank(cb Chessboard, rank int) int {
	if rank < 1 || rank > 8 {
		return 0
	}

	var c int

	for _, f := range cb {
		if f[rank-1] {
			c++
		}
	}

	return c
}

// CountAll should count how many squares are present in the chessboard.
func CountAll(cb Chessboard) int {
	var c int

	for _, f := range cb {
		c += len(f)
	}

	return c
}

// CountOccupied returns how many squares are occupied in the chessboard.
func CountOccupied(cb Chessboard) int {
	var c int

	for _, f := range cb {
		for _, v := range f {
			if v {
				c++
			}
		}
	}

	return c
}
