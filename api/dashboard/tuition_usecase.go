package dashboard

import (
	"log"
	"sort"

	"github.com/Tech-by-GL/dashboard/db"
	"github.com/Tech-by-GL/dashboard/helper"
	"github.com/jinzhu/copier"
)

type TuitionUsecaseService struct {
	Usecase *db.Querier
}

func NewTuitionUsecaseService() *TuitionUsecaseService {
	return &TuitionUsecaseService{
		Usecase: db.NewQuerier(),
	}
}

type OrderUsecase struct {
	ID             string
	UserID         int
	FullName       string
	ClassName      string
	Amount         float32
	Discount       float32
	ReasonDiscount string
}

type TuitionUsecase struct {
	TotalUser        int
	TotalOrder       int
	TotalAmount      float32
	ClassName        []string
	TotalUserInClass []int
	OthersItem       int
	OthersDetails    []OthersDetails
	Order            []OrderUsecase
}

type InvoiceUsecase struct {
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

type OthersDetails struct {
	OtherName   string
	TotalAmount int
	Price       float32
}

func (t *TuitionUsecaseService) QueryAlertStudentInTuition(date string) (res TuitionUsecase, err error) {
	var (
		orders []db.Order
		users  []int
		items  db.AllItem
	)
	orders, err = t.Usecase.QueryOrderByDate(date)
	if err != nil {
		return
	}

	err = copier.Copy(&res.Order, &orders)
	if err != nil {
		return
	}

	// Get All Items Name
	items, err = t.Usecase.GetAllItems()
	if err != nil {
		return
	}

	classCountMap := make(map[string]int, 0)
	otherDetails := make(map[string]float32, 0)
	otherDetailsCount := make(map[string]int, 0)

	res.TotalOrder = len(orders)
	for i, order := range orders {
		// Find Item Name
		itemName, typeID := items.FindItemName(order.ItemID)
		if typeID == "1" {
			if v := classCountMap[itemName]; v != 0 {
				classCountMap[itemName] += 1
			} else {
				classCountMap[itemName] = 1
			}
			res.Order[i].ClassName = itemName
		} else if typeID == "2" {
			res.OthersItem += 1
			if v := otherDetails[itemName]; v != 0 {
				otherDetails[itemName] += order.Amount - order.Discount
				otherDetailsCount[itemName] += 1
			} else {
				otherDetails[itemName] = order.Amount - order.Discount
				otherDetailsCount[itemName] = 1
			}
		} else {
			log.Printf("Cannot find itemName: %s or typeID: %s or ItemID: %s\n", itemName, typeID, order.ItemID)
		}

		res.TotalAmount += order.Amount - order.Discount
		if !helper.ContainsInArrayInt(users, order.UserID) {
			users = append(users, order.UserID)
		}
	}
	res.TotalUser = len(users)
	if err != nil {
		return
	}
	for k, v := range otherDetails {
		res.OthersDetails = append(res.OthersDetails, OthersDetails{
			OtherName:   k,
			TotalAmount: otherDetailsCount[k],
			Price:       v,
		})
	}

	// Sort the slice by values
	sort.Slice(res.OthersDetails, func(i, j int) bool {
		return res.OthersDetails[i].TotalAmount <= res.OthersDetails[j].TotalAmount
	})

	// Sort map by values
	keys := make([]string, 0, len(classCountMap))
	for key := range classCountMap {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, k := range keys {
		res.ClassName = append(res.ClassName, k)
		res.TotalUserInClass = append(res.TotalUserInClass, classCountMap[k])
	}
	return
}

func (t *TuitionUsecaseService) FindInvoiceKPIToday(date string) (res InvoiceUsecase, err error) {
	var (
		invoices             db.Invoice
		invoicePreviousDate  db.Invoice
		invoicePreviousMonth db.Invoice
	)
	invoices, err = t.Usecase.QueryInvoices(date)
	if err != nil {
		return
	}

	invoicePreviousDate, err = t.Usecase.QueryTotalAmountInvoicesByPreviousDate(date)
	if err != nil {
		log.Println("?")
		return
	}

	invoicePreviousMonth, err = t.Usecase.QueryTotalAmountInvoicesBySameDateInPreviousMonth(date)
	if err != nil {
		return
	}

	err = copier.Copy(&res, &invoices)
	res.InvoiceCompare.PreviousDateTotal = invoicePreviousDate.Total
	res.InvoiceCompare.PreviousMonthTotal = invoicePreviousMonth.Total
	return
}
