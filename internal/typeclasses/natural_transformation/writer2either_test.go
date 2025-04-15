package natural_transformation

import (
	"errors"
	"testing"

	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	P "github.com/IBM/fp-go/pair"
	W "github.com/IBM/fp-go/writer"
	"github.com/stretchr/testify/assert"
)

func TestWriter2Either(t *testing.T) {
	tests := []struct {
		name     string
		input    W.Writer[[]string, E.Either[error, int]]
		expected E.Either[error, W.Writer[[]string, int]]
	}{
		{
			name: "Left case",
			input: func() P.Pair[E.Either[error, int], []string] {
				return P.MakePair(E.Left[int](errors.New("Left case")), A.Of("Left case"))
			},
			expected: E.Left[W.Writer[[]string, int]](errors.New("Left case")),
		},
		{
			name: "Right case",
			input: func() P.Pair[E.Either[error, int], []string] {
				return P.MakePair(E.Right[error](42), A.Of("Success"))
			},
			expected: E.Right[error, W.Writer[[]string, int]](func() P.Pair[int, []string] {
				return P.MakePair(42, A.Of("Success"))
			}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := W2E(test.input)

			E.Fold(
				func(err error) bool {
					return assert.Equal(t, test.expected, E.Left[W.Writer[[]string, int]](err))
				},
				func(w W.Writer[[]string, int]) bool {
					assert.Equal(t, 42, W.Evaluate(w))
					return assert.Equal(t, "Success", W.Execute(w)[0])
				},
			)(result)
		})
	}
}
