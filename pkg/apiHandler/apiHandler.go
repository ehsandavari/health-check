package apiHandler

import (
	"github.com/ehsandavari/go-context-plus"
	"github.com/ehsandavari/go-logger"
	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go/types"
	"health-check/application/common"
	"health-check/pkg/tracer"
	"net/http"
)

type SBaseController struct {
	ILogger logger.ILogger
	ITracer tracer.ITracer
}

func NewBaseController(iLogger logger.ILogger, iTracer tracer.ITracer) SBaseController {
	return SBaseController{
		ILogger: iLogger,
		ITracer: iTracer,
	}
}

type BaseController[TReq, TRes any] func(ctx *contextplus.Context, request TReq) (TRes, error)

func (r BaseController[TReq, TRes]) Handle(iLogger logger.ILogger) gin.HandlerFunc {
	return func(ctxGin *gin.Context) {
		reqCtx := ctxGin.Request.Context()
		ctx := contextplus.FromContext(reqCtx)

		var request TReq
		if _, ok := any(request).(*types.Nil); !ok {
			if bindErr := ctxGin.ShouldBind(&request); bindErr != nil {
				iLogger.WithError(bindErr).Warn(ctx, "error in Bind request")
				err := NewApiError(http.StatusBadRequest, "error in validate request")
				if validationErrors, ok := bindErr.(validator.ValidationErrors); ok {
					meta := make(map[string]string, len(validationErrors))
					for _, validationError := range validationErrors {
						meta[validationError.Field()] = validationError.Tag()
					}
					err.SetMeta(meta)
				}
				ctxGin.JSON(http.StatusBadRequest, NewBaseApiResponse[ApiError](
					false,
					err,
				))
				return
			}
		}

		result, err := r(ctx, request)
		if err != nil {
			iError := err.(common.IError)
			iLogger.WithUint("ErrorCode", iError.Code()).WithError(iError).Debug(ctx, "error in handler")
			ctxGin.JSON(http.StatusInternalServerError, NewBaseApiResponse[ApiError](
				false,
				NewApiError(iError.Code(), i18n.MustGetMessage(ctxGin, iError.Error())),
			))
			return
		}

		ctxGin.JSON(http.StatusOK, NewBaseApiResponse[TRes](
			true,
			result,
		))
	}
}

type BaseApiResponse[TD any] struct {
	IsSuccess bool `json:"isSuccess"`
	Data      TD   `json:"data"`
}

func NewBaseApiResponse[TD any](isSuccess bool, data TD) BaseApiResponse[TD] {
	return BaseApiResponse[TD]{
		IsSuccess: isSuccess,
		Data:      data,
	}
}

type ApiError struct {
	Code    uint   `json:"code" format:"uint32"`
	Message string `json:"message"`
	Meta    any    `json:"meta,omitempty" extensions:"x-nullable,x-omitempty"`
}

func NewApiError(code uint, message string) ApiError {
	return ApiError{
		Code:    code,
		Message: message,
	}
}

func (r *ApiError) SetMeta(meta any) {
	r.Meta = meta
}
