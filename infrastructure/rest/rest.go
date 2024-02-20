package rest

import (
	"github.com/ehsandavari/go-context-plus"
	"github.com/ehsandavari/go-logger"
	"github.com/go-resty/resty/v2"
	"health-check/application/interfaces"
	"health-check/domain/enums"
	"net/http"
)

type sRest struct {
	iLogger logger.ILogger
	client  *resty.Client
}

func NewRest(logger logger.ILogger) interfaces.IRest {
	return &sRest{
		iLogger: logger,
		client: resty.New().
			SetPreRequestHook(
				func(c *resty.Client, r *http.Request) error {
					logger.WithHttpRequest(r).Info(contextplus.FromContext(r.Context()), "request")
					return nil
				},
			).
			OnAfterResponse(
				func(c *resty.Client, r *resty.Response) error {
					logger.WithHttpResponse(r.RawResponse).Info(contextplus.FromContext(r.Request.Context()), "response")
					return nil
				},
			),
	}
}

func (r *sRest) Execute(ctx *contextplus.Context, method enums.HttpMethod, url string, headers map[string]string, body map[string]any) (int, http.Header, string, error) {
	resp, err := r.client.R().
		SetContext(ctx).
		SetHeaders(headers).
		SetBody(body).
		Execute(method.String(), url)
	if err != nil {
		r.iLogger.WithError(err).Error(contextplus.Background(), "error in Execute request")
		return 0, nil, "", err
	}

	return resp.StatusCode(), resp.Header(), resp.String(), nil
}
