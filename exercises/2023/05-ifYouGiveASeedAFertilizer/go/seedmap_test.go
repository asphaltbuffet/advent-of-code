package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_categoryFromLine(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want Category
	}{
		{
			name: "water",
			args: args{
				line: "water-to-light map:",
			},
			want: WaterCategory,
		},
		{
			name: "invalid",
			args: args{
				line: "lead-to-gold map:",
			},
			want: UpperBound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, categoryFromLine(tt.args.line))
		})
	}
}

func Test_parseMap(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want *SeedMap
	}{
		{
			name: "single offset",
			args: args{
				s: "seed-to-soil map:\n50 98 2",
			},
			want: &SeedMap{
				Type:    SeedCategory,
				Offsets: []SeedMapOffset{{Source: 98, Dest: 50, Range: 2}},
			},
		},
		{
			name: "multiple offsets",
			args: args{
				s: "fertilizer-to-water map:\n49 53 8\n0 11 42\n42 0 7\n57 7 4",
			},
			want: &SeedMap{
				Type: FertCategory,
				Offsets: []SeedMapOffset{
					{Source: 53, Dest: 49, Range: 8},
					{Source: 11, Dest: 0, Range: 42},
					{Source: 0, Dest: 42, Range: 7},
					{Source: 7, Dest: 57, Range: 4},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, parseMap(tt.args.s))
		})
	}
}

func Test_translate(t *testing.T) {
	type args struct {
		m    *SeedMap
		seed int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "one offset range with match",
			args: args{
				m: &SeedMap{
					Type:    SeedCategory,
					Offsets: []SeedMapOffset{{Source: 0, Dest: 100, Range: 10}},
				},
				seed: 5,
			},
			want: 105,
		},
		{
			name: "one offset range without match",
			args: args{
				m: &SeedMap{
					Type:    SeedCategory,
					Offsets: []SeedMapOffset{{Source: 0, Dest: 100, Range: 10}},
				},
				seed: 50,
			},
			want: 50,
		},
		{
			name: "mult offsets with later match",
			args: args{
				m: &SeedMap{
					Type: SeedCategory,
					Offsets: []SeedMapOffset{
						{Source: 0, Dest: 100, Range: 10},
						{Source: 50, Dest: 0, Range: 10},
					},
				},
				seed: 51,
			},
			want: 1,
		},
		{
			name: "mult offsets with no match",
			args: args{
				m: &SeedMap{
					Type: SeedCategory,
					Offsets: []SeedMapOffset{
						{Source: 0, Dest: 100, Range: 10},
						{Source: 50, Dest: 0, Range: 10},
					},
				},
				seed: 500,
			},
			want: 500,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, translate(tt.args.m, tt.args.seed))
		})
	}
}

func Test_getLocation(t *testing.T) {
	type args struct {
		maps map[Category]*SeedMap
		seed int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "not in any ranges",
			args: args{
				maps: map[Category]*SeedMap{
					SeedCategory:  {Type: SeedCategory, Offsets: []SeedMapOffset{{Source: 0, Dest: 2, Range: 2}}},
					SoilCategory:  {Type: SoilCategory, Offsets: []SeedMapOffset{{Source: 10, Dest: 20, Range: 2}}},
					FertCategory:  {Type: FertCategory, Offsets: []SeedMapOffset{{Source: 30, Dest: 40, Range: 2}}},
					WaterCategory: {Type: WaterCategory, Offsets: []SeedMapOffset{{Source: 50, Dest: 60, Range: 2}}},
					LightCategory: {Type: LightCategory, Offsets: []SeedMapOffset{{Source: 70, Dest: 80, Range: 2}}},
					TempCategory:  {Type: TempCategory, Offsets: []SeedMapOffset{{Source: 90, Dest: 0, Range: 2}}},
					HumidCategory: {Type: HumidCategory, Offsets: []SeedMapOffset{{Source: 15, Dest: 25, Range: 2}}},
				},
				seed: 1000,
			},
			want: 1000,
		},
		{
			name: "first range only",
			args: args{
				maps: map[Category]*SeedMap{
					SeedCategory:  {Type: SeedCategory, Offsets: []SeedMapOffset{{Source: 0, Dest: 2, Range: 2}}},
					SoilCategory:  {Type: SoilCategory, Offsets: []SeedMapOffset{{Source: 10, Dest: 20, Range: 2}}},
					FertCategory:  {Type: FertCategory, Offsets: []SeedMapOffset{{Source: 30, Dest: 40, Range: 2}}},
					WaterCategory: {Type: WaterCategory, Offsets: []SeedMapOffset{{Source: 50, Dest: 60, Range: 2}}},
					LightCategory: {Type: LightCategory, Offsets: []SeedMapOffset{{Source: 70, Dest: 80, Range: 2}}},
					TempCategory:  {Type: TempCategory, Offsets: []SeedMapOffset{{Source: 90, Dest: 0, Range: 20}}},
					HumidCategory: {Type: HumidCategory, Offsets: []SeedMapOffset{{Source: 15, Dest: 25, Range: 2}}},
				},
				seed: 0,
			},
			want: 2,
		},
		{
			name: "two ranges",
			args: args{
				maps: map[Category]*SeedMap{
					SeedCategory:  {Type: SeedCategory, Offsets: []SeedMapOffset{{Source: 0, Dest: 2, Range: 2}}},
					SoilCategory:  {Type: SoilCategory, Offsets: []SeedMapOffset{{Source: 10, Dest: 20, Range: 2}}},
					FertCategory:  {Type: FertCategory, Offsets: []SeedMapOffset{{Source: 30, Dest: 40, Range: 2}}},
					WaterCategory: {Type: WaterCategory, Offsets: []SeedMapOffset{{Source: 50, Dest: 60, Range: 2}}},
					LightCategory: {Type: LightCategory, Offsets: []SeedMapOffset{{Source: 70, Dest: 80, Range: 2}}},
					TempCategory:  {Type: TempCategory, Offsets: []SeedMapOffset{{Source: 90, Dest: 0, Range: 20}}},
					HumidCategory: {Type: HumidCategory, Offsets: []SeedMapOffset{{Source: 15, Dest: 25, Range: 2}}},
				},
				seed: 105,
			},
			want: 25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getLocation(tt.args.maps, tt.args.seed))
		})
	}
}

func Test_getLocations(t *testing.T) {
	type args struct {
		maps  map[Category]*SeedMap
		seeds []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "",
			args: args{
				maps: map[Category]*SeedMap{
					SeedCategory:  {Type: SeedCategory, Offsets: []SeedMapOffset{{Source: 0, Dest: 2, Range: 2}}},
					SoilCategory:  {Type: SoilCategory, Offsets: []SeedMapOffset{{Source: 10, Dest: 20, Range: 2}}},
					FertCategory:  {Type: FertCategory, Offsets: []SeedMapOffset{{Source: 30, Dest: 40, Range: 2}}},
					WaterCategory: {Type: WaterCategory, Offsets: []SeedMapOffset{{Source: 50, Dest: 60, Range: 2}}},
					LightCategory: {Type: LightCategory, Offsets: []SeedMapOffset{{Source: 70, Dest: 80, Range: 2}}},
					TempCategory:  {Type: TempCategory, Offsets: []SeedMapOffset{{Source: 90, Dest: 0, Range: 2}}},
					HumidCategory: {Type: HumidCategory, Offsets: []SeedMapOffset{{Source: 15, Dest: 25, Range: 2}}},
				},
				seeds: []int{1000, 0},
			},
			want: []int{1000, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getLocations(tt.args.maps, tt.args.seeds))
		})
	}
}

func Test_parseSeedRange(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []SeedRange
	}{
		{
			name: "example",
			args: args{
				s: "seeds: 79 14 55 13",
			},
			want: []SeedRange{
				{Start: 79, Range: 14},
				{Start: 55, Range: 13},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.ElementsMatch(t, tt.want, parseSeedRange(tt.args.s))
		})
	}
}
