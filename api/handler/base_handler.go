package handler

import (
	"github.com/Tech-by-GL/dashboard/appctx"
	"github.com/Tech-by-GL/dashboard/handler/presenter"
	"github.com/gin-gonic/gin"
)

// BaseHandler help us respond to client
type BaseHandler struct{}

// SetMeta to put meta information into context
func (h *BaseHandler) SetMeta(ctx *gin.Context, meta presenter.MetaResponse) {
	newCtx := appctx.SetValue(ctx.Request.Context(), appctx.MetaContextKey, meta)
	ctx.Request = ctx.Request.WithContext(newCtx)
}

// SetData to put data information into context
func (h *BaseHandler) SetData(ctx *gin.Context, data interface{}) {
	newCtx := appctx.SetValue(ctx.Request.Context(), appctx.DataContextKey, data)
	ctx.Request = ctx.Request.WithContext(newCtx)
}

func (h *BaseHandler) SetText(ctx *gin.Context, text interface{}) {
	newCtx := appctx.SetValue(ctx.Request.Context(), appctx.TextContextKey, text)
	ctx.Request = ctx.Request.WithContext(newCtx)
}

// SetError to put meta information into context
func (h *BaseHandler) SetError(ctx *gin.Context, err error) {
	newCtx := appctx.SetValue(ctx.Request.Context(), appctx.ErrorContextKey, err)
	ctx.Request = ctx.Request.WithContext(newCtx)

}
