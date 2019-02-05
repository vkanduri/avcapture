package recorder

import (
	"context"

	"time"

	"github.com/go-kit/kit/log"

	"github.com/etherlabsio/pkg/logutil"
)

type Middleware func(Service) Service

func LoggingMiddleware(l log.Logger) Middleware {
	return func(svc Service) Service {
		return loggingMiddleware{
			logger: l,
			next:   svc,
		}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw loggingMiddleware) Start(ctx context.Context, req StartRecordingRequest) (resp StartRecordingResponse) {
	defer func(begin time.Time) {
		logutil.
			WithError(mw.logger, resp.Failed()).
			Log("method", "Start", "request", req, "span", time.Since(begin))
	}(time.Now())
	return mw.next.Start(ctx, req)
}

func (mw loggingMiddleware) Stop(ctx context.Context, req StopRecordingRequest) (resp StopRecordingResponse) {
	defer func(begin time.Time) {
		logutil.
			WithError(mw.logger, resp.Failed()).
			Log("method", "Stop", "request", req, "span", time.Since(begin))
	}(time.Now())
	return mw.next.Stop(ctx, req)
}
