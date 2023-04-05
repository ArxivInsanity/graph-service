package util

import "context"

type contextCloser interface {
	Close(ctx context.Context) error
}

func PanicOnClosureError(err error, ctx context.Context, closer contextCloser) {
	if err != nil {
		PanicOnErr(closer.Close(ctx))
	}
}

func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
