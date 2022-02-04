package context

import (
	"context"

	"github.com/sirupsen/logrus"
)

var loggerContext = "_context_logger"

// Logger returns the logger from the context
func Logger(ctx context.Context) *logrus.Entry {
	return ctx.Value(&loggerContext).(*logrus.Entry)
}

// WithLogger inserts a logger into the context
func WithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, &loggerContext, logger)
}
