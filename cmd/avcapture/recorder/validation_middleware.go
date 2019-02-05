package recorder

import (
	"context"

	"github.com/etherlabsio/errors"
)

func ValidationMiddleware(svc Service) Service {
	return validationMiddleware{next: svc}
}

type validationMiddleware struct {
	next Service
}

func (mw validationMiddleware) Start(ctx context.Context, req StartRecordingRequest) (resp StartRecordingResponse) {
	var op errors.Op = "Start"
	if 0 == len(req.FFmpeg.Params) {
		resp.Err = errors.New("missing output params for encoding(ffmpeg)", op, errors.Invalid)
		return resp
	}
	if "" == req.Chrome.URL {
		resp.Err = errors.New("missing URL for capturing", op, errors.Invalid)
		return resp
	}
	return mw.next.Start(ctx, req)
}

func (mw validationMiddleware) Stop(ctx context.Context, req StopRecordingRequest) StopRecordingResponse {
	return mw.next.Stop(ctx, req)
}
