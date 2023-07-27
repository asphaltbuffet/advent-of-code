package exercises

import "testing"

func BenchmarkParse(b *testing.B) {
	line := "Blueprint 123: Each ore robot costs 4 ore. Each clay robot costs 6 ore. Each obsidian robot costs 8 ore and 2 clay. Each geode robot costs 10 ore and 5 obsidian."
	for i := 0; i < b.N; i++ {
		_, _ = parse(line)
	}
}

// func BenchmarkParseWithRegex(b *testing.B) {
// 	line := "Blueprint 123: Each ore robot costs 4 ore. Each clay robot costs 6 ore. Each obsidian robot costs 8 ore and 2 clay. Each geode robot costs 10 ore and 5 obsidian."
// 	for i := 0; i < b.N; i++ {
// 		_, _ = parseWithRegex(line)
// 	}
// }
