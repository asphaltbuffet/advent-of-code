package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RTLParse(t *testing.T) {

	tests := []struct {
		name      string
		arg       string
		want      []Problem
		assertion require.ErrorAssertionFunc
	}{
		{
			name: "all 2digit numbers",
			arg:  "321\n321\n+  ",
			want: []Problem{
				Problem{[]int{11, 22, 33}, "+"},
			},
			assertion: require.NoError,
		},
		{
			name: "right aligned",
			arg:  "3 1\n321\n+  ",
			want: []Problem{
				Problem{[]int{11, 2, 33}, "+"},
			},
			assertion: require.NoError,
		},
		{
			name: "left aligned",
			arg:  "321\n 21\n*  ",
			want: []Problem{
				Problem{[]int{11, 22, 3}, "*"},
			},
			assertion: require.NoError,
		},
		{
			name: "2 problems",
			arg:  "654 321\n654 321\n+   +  ",
			want: []Problem{
				Problem{[]int{11, 22, 33}, "+"},
				Problem{[]int{44, 55, 66}, "+"},
			},
			assertion: require.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RTLParse(tt.arg)

			tt.assertion(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
