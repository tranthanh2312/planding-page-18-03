package middleware

import (
	"net/http"

	"github.com/Tech-by-GL/dashboard/appctx"
	"github.com/Tech-by-GL/dashboard/handler/presenter"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func JSONWriterMiddleware(ctx *gin.Context) {
	ctx.Next()

	// Check error if exists
	// Base on error/success to return meta object
	var (
		res      presenter.Response
		httpCode int
	)

	appErr := appctx.GetValue(ctx.Request.Context(), appctx.ErrorContextKey)
	if appErr != nil {
		_processAppError(&res, appErr)
		httpCode = res.Meta.Code
	}

	appText := appctx.GetValue(ctx.Request.Context(), appctx.TextContextKey)
	if appText != nil {
		res.Text = appText
	}

	// Respond the data object/array
	data := appctx.GetValue(ctx.Request.Context(), appctx.DataContextKey)
	if data != nil {
		res.Data = data
	}

	meta := appctx.GetValue(ctx.Request.Context(), appctx.MetaContextKey)
	if meta != nil {
		metaRes, ok1 := meta.(presenter.MetaResponse)
		if ok1 {
			res.Meta = metaRes
			httpCode = metaRes.Code
		}
	}

	if res.IsEmpty() {
		ctx.JSON(http.StatusNoContent, nil)
	} else {
		// struct ---(serialize)----> JSON ------> http.write -----> http.flush
		ctx.JSON(httpCode, res)
	}

}

func _processAppError(res *presenter.Response, appErr interface{}) {
	// bindingErr represents the error comes from validator.Validate (v10).
	bindingErr := _catchBindingError(appErr.(error))
	if bindingErr != nil {
		res.Errors = bindingErr.(presenter.ErrorResponses)
		res.Meta = presenter.MetaResponse{
			Code:    http.StatusBadRequest,
			Message: "error when binding the request",
		}

		return
	}

	// API-designed error.
	res.Errors = presenter.ErrorResponses{
		presenter.ErrorResponse{
			Code:   http.StatusBadRequest,
			Detail: appErr.(error).Error(),
			Source: &presenter.SourceResponse{
				Pointer:   "API Layer",
				Parameter: "Query-URL",
			},
		},
	}
}

func _catchBindingError(appErr error) error {
	var errs presenter.ErrorResponses

	// this check is only needed when your code could produce
	// an invalid value for validation such as interface with nil
	// value most including myself do not usually have code like this.
	if _, ok := appErr.(*validator.InvalidValidationError); ok {
		errs.Append(presenter.ErrorResponse{
			Code:   http.StatusInternalServerError,
			Detail: "invalid validation error",
		})
		return errs
	}

	if vldrErr, ok := appErr.(validator.ValidationErrors); ok {
		errs.FromValidationErrors(vldrErr)
		return errs
	}

	return nil
}
