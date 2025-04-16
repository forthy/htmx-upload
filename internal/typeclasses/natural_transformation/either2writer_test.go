package natural_transformation

import (
	"errors"
	"fmt"
	"log"
	"testing"

	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	P "github.com/IBM/fp-go/pair"
	W "github.com/IBM/fp-go/writer"

	"github.com/stretchr/testify/assert"
)

func TestEither2Writer(t *testing.T) {
	tests := []struct {
		name     string
		input    E.Either[error, W.Writer[[]string, int]]
		expected W.Writer[[]string, E.Either[error, int]]
	}{
		{
			name:  "Left case",
			input: E.Left[W.Writer[[]string, int]](errors.New("Left case")),
			expected: func() P.Pair[E.Either[error, int], []string] {
				return P.MakePair(E.Left[int](errors.New("Left case")), A.Of("Left case"))
			},
		},
		{
			name:  "Right case",
			input: E.Right[error, W.Writer[[]string, int]](func() P.Pair[int, []string] { return P.MakePair(42, A.Of("Success")) }),
			expected: func() P.Pair[E.Either[error, int], []string] {
				return P.MakePair(E.Right[error](42), A.Of("Success"))
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := E2W[[]string, int](func(err error) []string {
				return A.Of(fmt.Sprintf("Error:[%v]", err))
			})(test.input)

			// DEBUG
			log.Println(fmt.Sprintf("Result:[%s]", W.Evaluate(result)))
			log.Println(fmt.Sprintf("Expected:[%s]", W.Evaluate(test.expected)))

			assert.Equal(t, W.Evaluate(test.expected), W.Evaluate(result))
		})
	}
}
