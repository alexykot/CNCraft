package db

import (
	"context"
	"time"
)

var dbCtx context.Context

const queryTimeout = time.Second * 30

// CtxWithCancel provides a context with timeout suitable for database queries, with the corresponding cancelFunc
func CtxWithCancel() (context.Context, context.CancelFunc) {
	return context.WithTimeout(dbCtx, queryTimeout)
}

// Ctx provides a context with timeout suitable for database queries, and without the corresponding cancelFunc.
// In most cases the consumers of this context can safely assume that the parent context already has a cancelfunc
// which is correctly handle, so this cancelFunc can be safely ignored.
func Ctx() context.Context {
	ctx, _ := CtxWithCancel()
	return ctx
}
