package db

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
)

type Order struct {
	ID             string
	UserID         int
	FullName       string
	ItemID         string
	Amount         float32
	Discount       float32
	ReasonDiscount string
}

func (q *Querier) QueryOrderByDate(date string) (res []Order, err error) {
	var users AllUser
	users, err = q.GetAllUser()
	if err != nil {
		return
	}
	dateSplit := strings.Split(date, "-")
	year, err := strconv.Atoi(dateSplit[0])
	if err != nil {
		return res, err
	}
	month, err := strconv.Atoi(dateSplit[1])
	if err != nil {
		return
	}
	firstDateInMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	lastDateInMonth := firstDateInMonth.AddDate(0, 1, -1).String()[:10]

	var rows *sql.Rows

	rows, err = q.DB.Query(`SELECT o.id, -(o.amount), o.discount, o.reason_discount, ui.user_id, ui.item_id FROM (orders as o INNER JOIN user_item_order as uio ON o.id = uio.order_id) INNER JOIN user_item as ui ON uio.user_item_id = ui.id
							WHERE o.date_created >= ? and o.date_created <= ? and o.active = 1 and o.status = 0`, firstDateInMonth.String()[:19], lastDateInMonth+" 23:59:59")
	if err != nil {
		return
	}

	for rows.Next() {
		var order Order
		err = rows.Scan(&order.ID, &order.Amount, &order.Discount, &order.ReasonDiscount, &order.UserID, &order.ItemID)
		if err != nil {
			return
		}

		if (order.Amount - order.Discount) == 0 {
			continue
		}
		order.FullName = GetFullName(order.UserID, users.Users)
		res = append(res, order)
	}

	return res, err
}
