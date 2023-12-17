package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseLine(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name      string
		args      args
		want      *Record
		assertion require.ErrorAssertionFunc
	}{
		{
			name:      "three checksums",
			args:      args{s: "???.### 1,1,3"},
			want:      &Record{Condition: "uuu.ddd", Checksum: []int{1, 1, 3}},
			assertion: require.NoError,
		},
		{
			name:      "double-digit checksum value",
			args:      args{s: "???.### 1,21,3"},
			want:      &Record{Condition: "uuu.ddd", Checksum: []int{1, 21, 3}},
			assertion: require.NoError,
		},
		{
			name:      "single checksum",
			args:      args{s: "???.### 1"},
			want:      &Record{Condition: "uuu.ddd", Checksum: []int{1}},
			assertion: require.NoError,
		},
		{
			name:      "invalid checksum",
			args:      args{s: "???.### 1,b,3"},
			want:      nil,
			assertion: require.Error,
		},
		{
			name:      "no space between sections",
			args:      args{s: "???.###1,1,3"},
			want:      nil,
			assertion: require.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLine(tt.args.s)

			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_generateRegex(t *testing.T) {
	type args struct {
		sizes []int
	}
	tests := []struct {
		name      string
		args      args
		want      string
		assertion require.ErrorAssertionFunc
	}{
		{
			name: "three sections",
			args: args{
				sizes: []int{1, 1, 3},
			},
			want:      `^[u\.]*[ud][u\.]+[ud][u\.]+[ud]{3}[u\.]*$`,
			assertion: require.NoError,
		},
		{
			name: "four sections",
			args: args{
				sizes: []int{1, 3, 1, 6},
			},
			want:      `^[u\.]*[ud][u\.]+[ud]{3}[u\.]+[ud][u\.]+[ud]{6}[u\.]*$`,
			assertion: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateRegex(tt.args.sizes)

			tt.assertion(t, err)
			if err == nil {
				assert.Equal(t, tt.want, got.String())
			}
		})
	}
}

func TestRecord_CountCombinations(t *testing.T) {
	tests := []struct {
		name      string
		r         *Record
		want      int
		assertion require.ErrorAssertionFunc
	}{
		{
			name:      "one way to match",
			r:         &Record{Condition: "uuu.ddd", Checksum: []int{1, 1, 3}},
			want:      1,
			assertion: require.NoError,
		},
		{
			name:      "4 ways",
			r:         &Record{Condition: ".uu..uu...udd.", Checksum: []int{1, 1, 3}},
			want:      4,
			assertion: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.r.countCombinations()

			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_expandAndParseLine(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name      string
		args      args
		want      *Record
		assertion require.ErrorAssertionFunc
	}{
		{
			name:      "small",
			args:      args{s: ".# 1"},
			want:      &Record{".du.du.du.du.d", []int{1, 1, 1, 1, 1}},
			assertion: require.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expandAndParseLine(tt.args.s)
			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_replaceDotsAndJoin(t *testing.T) {
	type args struct {
		left   string
		middle string
		right  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"all d", args{"dd", "d", "dd"}, "ddddd"},
		{"all u", args{"du", "d", "ud"}, "dudud"},
		{"no replacement", args{"d.u", "d", ".ud"}, "d.ud.ud"},
		{"left replacement", args{"d..u", "d", ".ud"}, "d.ud.ud"},
		{"right replacement", args{"d.u", "d", "u..d"}, "d.udu.d"},
		{"overlap left replacement", args{"du.", ".", "ud"}, "du.ud"},
		{"overlap right replacement", args{"du", ".", ".ud"}, "du.ud"},
		{"overlap all replacement", args{"du.", ".", ".ud"}, "du.ud"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, replaceDotsAndJoin(tt.args.left, tt.args.middle, tt.args.right))
		})
	}
}
