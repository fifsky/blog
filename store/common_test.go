package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIn(t *testing.T) {
	type args[T any] struct {
		s          []T
		startIndex int
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
				s:          []int{},
				startIndex: 1,
			},
			want:  "",
			want1: nil,
		},
		{
			name: "single",
			args: args[int]{
				s:          []int{1},
				startIndex: 1,
			},
			want:  "$1",
			want1: []any{1},
		},
		{
			name: "multiple",
			args: args[int]{
				s:          []int{1, 2, 3},
				startIndex: 1,
			},
			want:  "$1,$2,$3",
			want1: []any{1, 2, 3},
		},
		{
			name: "multiple with offset",
			args: args[int]{
				s:          []int{1, 2, 3},
				startIndex: 5,
			},
			want:  "$5,$6,$7",
			want1: []any{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := In(tt.args.s, tt.args.startIndex)
			assert.Equalf(t, tt.want, got, "In(%v, %d)", tt.args.s, tt.args.startIndex)
			assert.Equalf(t, tt.want1, got1, "In(%v, %d)", tt.args.s, tt.args.startIndex)
		})
	}
}
