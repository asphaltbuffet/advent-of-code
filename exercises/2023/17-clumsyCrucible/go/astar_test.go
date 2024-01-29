package exercises

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_hashNode(t *testing.T) {
	n1 := &node{pather: &Block{Position: Point{1, 1}}, parent: nil}
	n2 := &node{pather: &Block{Position: Point{2, 2}}, parent: n1}
	n3 := &node{pather: &Block{Position: Point{3, 3}}, parent: n2}
	n4 := &node{pather: &Block{Position: Point{4, 4}}, parent: n3}
	n5 := &node{pather: &Block{Position: Point{5, 5}}, parent: n4}

	type args struct {
		n *node
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"nil", args{n: nil}, ""},
		{"empty", args{n: &node{}}, "<nil>"},
		{"no parent", args{n: n1}, "<1-1>"},
		{"three nodes", args{n: n3}, "<3-3><2-2><1-1>"},
		{"truncated", args{n: n5}, "<5-5><4-4><3-3><2-2>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, getNodeHash(tt.args.n))
		})
	}
}
