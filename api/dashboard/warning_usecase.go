package dashboard

import (
	"context"

	"github.com/Tech-by-GL/dashboard/dashboard/usecase_dto"
	"github.com/Tech-by-GL/dashboard/db"
	"github.com/Tech-by-GL/dashboard/db/entity"
	"github.com/Tech-by-GL/dashboard/helper"
)

func (ti *TuitionUsecaseService) FindWarningIncompatibleOrder(ctx context.Context, date string) (res []usecase_dto.IncompatibleOrder, err error) {
	var (
		incompatibleOrders []entity.IncompatibleOrder
	)

	incompatibleOrders, err = ti.Usecase.QueryIncompatibleOrder(ctx, date)
	if err != nil {
		return
	}

	// Get All Users
	var (
		users db.AllUser
		// classes usecase_dto.AllClassCache
	)

	users, err = ti.Usecase.GetAllUser()
	if err != nil {
		return
	}

	for _, data := range incompatibleOrders {
		var (
			incompatibleOrder usecase_dto.IncompatibleOrder
			items             []entity.Item
			classItem         []entity.ClassItem
		)

		// * This logical notion check whether the user has already bought the item or not. If the user has already bought the item, then the order is not incompatible.
		classItem, err = ti.Usecase.QueryClassItemByUser(ctx, data.UserID, data.ItemID)
		if err != nil {
			return
		}

		if len(classItem) > 0 {
			continue
		}

		incompatibleOrder.UserID = data.UserID
		incompatibleOrder.FullName = users.SearchUserFullName(data.UserID)
		items, err = ti.Usecase.QueryItem(ctx, data.ItemID, "", 1)
		if err != nil {
			return
		}

		incompatibleOrder.ItemName = items[0].ItemName
		incompatibleOrder.Price = (-data.Amount) - (data.Discount)
		incompatibleOrder.Reason = "Học sinh mua sách không đúng lớp"
		res = append(res, incompatibleOrder)
	}
	return
}

func (ti *TuitionUsecaseService) FindWarningOrderWithFutureDate(ctx context.Context) (res []usecase_dto.WarningOrderWithFutureDate, err error) {
	var (
		ent []entity.WarningOrderWithFutureDate
	)

	ent, err = ti.Usecase.QueryWarningOrderWithFutureDate(ctx)
	if err != nil {
		return
	}

	users, err := ti.Usecase.GetAllUser()
	if err != nil {
		return
	}

	for _, data := range ent {
		var (
			warning usecase_dto.WarningOrderWithFutureDate
		)

		warning.OrderID = data.OrderID

		warning.Note = data.Note
		warning.ItemName = data.BackNote
		warning.CreatedUserFullName = users.SearchUserFullName(data.CreatedUserID)
		warning.DateCreated = helper.ConvertUnixTimeMySqlTime(data.DateCreated)[:10]

		res = append(res, warning)
	}
	return
}

func (ti *TuitionUsecaseService) FindWarningDuplicatedOrders(ctx context.Context, date string) (res []usecase_dto.WarningDuplicatedOrder, err error) {
	var (
		ent []entity.UserItem
	)

	ent, err = ti.Usecase.QueryAmountOfDuplicatedOrders(ctx, date)
	if err != nil {
		return
	}

	if len(ent) == 0 {
		return
	}

	users, err := ti.Usecase.GetAllUser()
	if err != nil {
		return
	}

	for _, data := range ent {
		var (
			warning usecase_dto.WarningDuplicatedOrder
			orders  []entity.Order
		)

		orders, err = ti.Usecase.QueryDuplicatedOrders(ctx, date, data.UserID, data.ItemID)
		if err != nil {
			return
		}

		warning.OrderID = orders[0].ID
		warning.UserID = data.UserID
		warning.FullName = users.SearchUserFullName(data.UserID)
		warning.ItemName = orders[0].BackNote
		warning.Price = (-orders[0].Amount) - orders[0].Discount
		warning.Reason = "Order bị trùng lặp"
		res = append(res, warning)
	}
	return
}
