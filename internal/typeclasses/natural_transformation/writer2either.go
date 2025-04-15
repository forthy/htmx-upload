package natural_transformation

import (
	E "github.com/IBM/fp-go/either"
	P "github.com/IBM/fp-go/pair"
	WR "github.com/IBM/fp-go/writer"
)

func W2E[W, A any](w WR.Writer[W, E.Either[error, A]]) E.Either[error, WR.Writer[W, A]] {
	return E.Map[error](func(a A) WR.Writer[W, A] {
		return func() P.Pair[A, W] {
			return P.MakePair(a, WR.Execute(w))
		}
	})(WR.Evaluate(w))
}
