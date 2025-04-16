package natural_transformation

import (
	E "github.com/IBM/fp-go/either"
	P "github.com/IBM/fp-go/pair"
	WR "github.com/IBM/fp-go/writer"
)

func E2W[W, A any](onLeft func(err error) W) func(E.Either[error, WR.Writer[W, A]]) WR.Writer[W, E.Either[error, A]] {
	return func(e E.Either[error, WR.Writer[W, A]]) WR.Writer[W, E.Either[error, A]] {
		wr, err := E.Unwrap(e)

		return func() P.Pair[E.Either[error, A], W] {
			if err != nil {
				return P.MakePair(E.Left[A](err), onLeft(err))
			}
			return P.MakePair(E.Right[error](WR.Evaluate(wr)), WR.Execute(wr))
		}
	}

}
