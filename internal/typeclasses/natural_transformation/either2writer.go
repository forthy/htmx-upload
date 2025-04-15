package natural_transformation

import (
	E "github.com/IBM/fp-go/either"
	L "github.com/IBM/fp-go/lazy"
	P "github.com/IBM/fp-go/pair"
	WR "github.com/IBM/fp-go/writer"
)

func E2W[W, A any](onLeft L.Lazy[W]) func(E.Either[error, WR.Writer[W, A]]) WR.Writer[W, E.Either[error, A]] {
	return func(e E.Either[error, WR.Writer[W, A]]) WR.Writer[W, E.Either[error, A]] {
		wr, err := E.Unwrap(e)

		return func() P.Pair[E.Either[error, A], W] {
			if err != nil {
				return P.MakePair(E.Left[A](err), onLeft())
			}
			return P.MakePair(E.Right[error](WR.Evaluate(wr)), WR.Execute(wr))
		}
	}

}
