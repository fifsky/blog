package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIn(t *testing.T) {
	type args[T any] struct {
		s []T
	}
	type testCase[T any] struct {
		name  string
		args  args[T]
		want  string
		want1 []any
	}
	tests := []testCase[int]{
		{
			name: "empty",
			args: args[int]{
				s: []int{},
			},
			want:  "",
			want1: nil,
		},
		{
			name: "single",
			args: args[int]{
				s: []int{1},
			},
			want:  "?",
			want1: []any{1},
		},
		{
			name: "multiple",
			args: args[int]{
				s: []int{1, 2, 3},
			},
			want:  "?,?,?",
			want1: []any{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := In(tt.args.s)
			assert.Equalf(t, tt.want, got, "In(%v)", tt.args.s)
			assert.Equalf(t, tt.want1, got1, "In(%v)", tt.args.s)
		})
	}
}
