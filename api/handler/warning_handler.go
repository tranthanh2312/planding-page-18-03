package handler

import (
	"net/http"

	"github.com/Tech-by-GL/dashboard/dashboard/usecase_dto"
	"github.com/Tech-by-GL/dashboard/handler/presenter"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

func (tudl *TuitionHandler) GetIncompatibleOrder(ctx *gin.Context) {
	var (
		err error
		dto []usecase_dto.IncompatibleOrder
		res []presenter.IncompatibleOrder
	)

	defer func() {
		tudl.SetError(ctx, err)
	}()

	// Get date from query
	date := ctx.Query("date")
	dto, err = tudl.TuitionHandler.FindWarningIncompatibleOrder(ctx, date)

	err = copier.Copy(&res, &dto)

	tudl.SetData(ctx, res)
	tudl.SetText(ctx, "Incompatible orders")
	tudl.SetMeta(ctx, presenter.MetaResponse{
		Code: http.StatusOK,
	})
}

func (tudl *TuitionHandler) GetWarningOrderWithFutureDate(ctx *gin.Context) {
	var (
		err error
		dto []usecase_dto.WarningOrderWithFutureDate
		res []presenter.WarningOrderWithFutureDate
	)

	defer func() {
		tudl.SetError(ctx, err)
	}()

	dto, err = tudl.TuitionHandler.FindWarningOrderWithFutureDate(ctx)

	err = copier.Copy(&res, &dto)

	tudl.SetData(ctx, res)
	tudl.SetMeta(ctx, presenter.MetaResponse{
		Code: http.StatusOK,
	})
}

func (tudl *TuitionHandler) GetWarningDuplicatedOrders(ctx *gin.Context) {
	var (
		err error
		dto []usecase_dto.WarningDuplicatedOrder
		res []presenter.WarningDuplicatedOrder
	)

	defer func() {
		tudl.SetError(ctx, err)
	}()

	date := ctx.Query("date")
	dto, err = tudl.TuitionHandler.FindWarningDuplicatedOrders(ctx, date)

	err = copier.Copy(&res, &dto)

	tudl.SetData(ctx, res)
	tudl.SetMeta(ctx, presenter.MetaResponse{
		Code: http.StatusOK,
	})
}
