package cx

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/truongnqse05461/ewallet/internal/log"
)

type contextKey string

const (
	traceKey    = contextKey("trace")
	loggerKey   = contextKey("logger")
	txKey       = contextKey("tx")
	userInfoKey = contextKey("userInfo")
)

func SetTrace(ctx context.Context, trace string) context.Context {
	return context.WithValue(ctx, traceKey, trace)
}

func GetTrace(ctx context.Context) (string, bool) {
	trace, ok := ctx.Value(traceKey).(string)
	return trace, ok
}

func SetLogger(ctx context.Context, logger log.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

func GetLogger(ctx context.Context) log.Logger {
	return ctx.Value(loggerKey).(log.Logger)
}

func SetTx(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

func GetTx(ctx context.Context) *sqlx.Tx {
	return ctx.Value(txKey).(*sqlx.Tx)
}
