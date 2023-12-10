package exercises

import (
	"fmt"
	"strconv"
	"strings"
)

type Category int

const (
	SeedCategory Category = iota + 1
	SoilCategory
	FertCategory
	WaterCategory
	LightCategory
	TempCategory
	HumidCategory
	LocCategory
	UpperBound // invalid plane
)

type SeedMap struct {
	Type    Category
	Offsets []SeedMapOffset
}

type SeedMapOffset struct {
	Source int
	Dest   int
	Range  int
}

func parseSeeds(s string) []int {
	seedTokens := strings.Fields(s)

	seeds := make([]int, 0, len(seedTokens))

	for _, token := range seedTokens[1:] {
		seed, _ := strconv.Atoi(token)
		seeds = append(seeds, seed)
	}

	return seeds
}

type SeedRange struct {
	Start int
	Range int
}

func parseSeedRange(s string) []SeedRange {
	seedTokens := strings.Fields(s)

	seeds := []SeedRange{}

	for i := 1; i < len(seedTokens); i += 2 {
		seed, _ := strconv.Atoi(seedTokens[i])
		r, _ := strconv.Atoi(seedTokens[i+1])

		// fmt.Printf("start: %d, range: %d\n", seed, r)

		sr := SeedRange{
			Start: seed,
			Range: r,
		}

		seeds = append(seeds, sr)
	}

	return seeds
}

func parseAllMaps(sections []string) map[Category]*SeedMap {
	maps := make(map[Category]*SeedMap, len(sections)-1)

	for _, section := range sections {
		sm := parseMap(section)
		maps[sm.Type] = sm
	}

	return maps
}

func parseMap(s string) *SeedMap {
	lines := strings.Split(s, "\n")

	offsets := make([]SeedMapOffset, 0, len(lines)-1)

	for _, line := range lines[1:] {
		var src, dst, rng int
		_, err := fmt.Sscanf(line, "%d %d %d", &dst, &src, &rng)
		if err != nil {
			panic(err)
		}

		offsets = append(offsets, SeedMapOffset{
			Source: src,
			Dest:   dst,
			Range:  rng,
		})

	}

	sm := &SeedMap{
		Type:    categoryFromLine(lines[0]),
		Offsets: offsets,
	}

	// fmt.Printf("parsed map: %+v\n", sm)

	return sm
}

func categoryFromLine(line string) Category {
	switch {
	case strings.HasPrefix(line, "seed-to-soil"):
		return SeedCategory
	case strings.HasPrefix(line, "soil-to-fertilizer"):
		return SoilCategory
	case strings.HasPrefix(line, "fertilizer-to-water"):
		return FertCategory
	case strings.HasPrefix(line, "water-to-light"):
		return WaterCategory
	case strings.HasPrefix(line, "light-to-temperature"):
		return LightCategory
	case strings.HasPrefix(line, "temperature-to-humidity"):
		return TempCategory
	case strings.HasPrefix(line, "humidity-to-location"):
		return HumidCategory
	default:
		return UpperBound
	}
}

func getLocations(maps map[Category]*SeedMap, seeds []int) []int {
	locations := make([]int, 0, len(seeds))

	for _, seed := range seeds {
		// fmt.Println("getLocations: ", seed)
		locations = append(locations, getLocation(maps, seed))
	}

	return locations
}

func getLocation(maps map[Category]*SeedMap, seed int) int {
	soil := translate(maps[SeedCategory], seed)
	fert := translate(maps[SoilCategory], soil)
	water := translate(maps[FertCategory], fert)
	light := translate(maps[WaterCategory], water)
	temp := translate(maps[LightCategory], light)
	humidity := translate(maps[TempCategory], temp)
	location := translate(maps[HumidCategory], humidity)

	return location
}

func translate(m *SeedMap, seed int) int {
	if m == nil {
		fmt.Printf("nil map for seed %d\n", seed)
		return seed
	}

	for _, offset := range m.Offsets {
		if seed >= offset.Source && seed < offset.Source+offset.Range {
			return offset.Dest + (seed - offset.Source)
		}
	}

	return seed
}
