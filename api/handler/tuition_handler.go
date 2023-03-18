package handler

import (
	"net/http"

	"github.com/Tech-by-GL/dashboard/dashboard"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type TuitionHandler struct {
	BaseHandler
	TuitionHandler *dashboard.TuitionUsecaseService
}

func NewHandler() *TuitionHandler {
	return &TuitionHandler{
		TuitionHandler: dashboard.NewTuitionUsecaseService(),
	}
}

type TuitionPresenter struct {
	TotalUser        int                     `json:"totalUser"`
	TotalOrder       int                     `json:"totalOrder"`
	TotalAmount      float32                 `json:"totalAmount"`
	ClassName        []string                `json:"className"`
	TotalUserInClass []int                   `json:"totalUserInClass"`
	OthersItem       int                     `json:"others"`
	OtherDetails     []OtherDetailsPresenter `json:"othersDetails"`
	Orders           []OrderPresenter        `json:"orders"`
}

type OrderPresenter struct {
	ID             string  `json:"id"`
	UserID         int     `json:"userId"`
	FullName       string  `json:"fullName"`
	ClassName      string  `json:"className"`
	Amount         float32 `json:"amount"`
	Discount       float32 `json:"discount"`
	ReasonDiscount string  `json:"reasonDiscount"`
}

type OtherDetailsPresenter struct {
	OtherName   string  `json:"otherName"`
	TotalAmount int     `json:"totalAmount"`
	Price       float32 `json:"price"`
}

func (h *TuitionHandler) GetTuition(ctx *gin.Context) {
	date := ctx.Query("date")
	if date == "" {
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "date is required"})
		return
	}
	var res TuitionPresenter
	orders, err := h.TuitionHandler.QueryAlertStudentInTuition(date)
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = copier.Copy(&res, &orders)
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = copier.Copy(&res.Orders, &orders.Order)

	if err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = copier.Copy(&res.OtherDetails, &orders.OthersDetails)

	if err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(200, res)
}

type InvoicePresenter struct {
	TotalCard      float32        `json:"totalCard"`
	TotalCash      float32        `json:"totalCash"`
	Total          float32        `json:"total"`
	TotalInvoices  int            `json:"totalInvoices"`
	InvoiceCompare InvoiceCompare `json:"invoiceCompare"`
}

type InvoiceCompare struct {
	PreviousDateTotal  float32 `json:"previousDateTotal"`
	PreviousMonthTotal float32 `json:"previousMonthTotal"`
}

func (h *TuitionHandler) GetKPIInvoice(ctx *gin.Context) {
	date := ctx.Query("date")
	if date == "" {
		ctx.Writer.WriteHeader(http.StatusBadRequest)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "date is required"})
		return
	}
	var res InvoicePresenter
	invoices, err := h.TuitionHandler.FindInvoiceKPIToday(date)
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = copier.Copy(&res, &invoices)
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusInternalServerError)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.IndentedJSON(200, res)
}
