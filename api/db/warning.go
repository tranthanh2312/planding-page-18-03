package db

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Tech-by-GL/dashboard/db/entity"
	"github.com/Tech-by-GL/dashboard/helper"
)

// **** Orders ****

func (q *Querier) QueryIncompatibleOrder(ctx context.Context, date string) (res []entity.IncompatibleOrder, err error) {
	var (
		rows *sql.Rows
	)

	splitDate := strings.Split(date, "-")
	year, err := strconv.Atoi(splitDate[0])
	if err != nil {
		return
	}
	month, err := strconv.Atoi(splitDate[1])
	if err != nil {
		return
	}

	firstDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	nextMonth := firstDate.AddDate(0, 1, 0)

	rows, err = q.DB.QueryContext(ctx, `select ui.user_id, ui.item_id, o.amount, o.discount from orders as o INNER JOIN user_item_order as uio
											ON o.id = uio.order_id and uio.active = 1  INNER JOIN user_item as ui
											ON uio.user_item_id = ui.id and ui.active = 1 
											where o.date_created >= ? and o.date_created < ? and o.note not like 'Học Phí%' 
											and item_id not in (select item_id from class_item where date_created >= ? and date_created < ? and buying = 0 group by item_id)`, firstDate, nextMonth, firstDate, nextMonth)

	if err != nil {
		return
	}

	for rows.Next() {
		var order entity.IncompatibleOrder
		err = rows.Scan(&order.UserID, &order.ItemID, &order.Amount, &order.Discount)
		if err != nil {
			return
		}

		res = append(res, order)
	}
	return
}

func (q *Querier) QueryClassItemByUser(ctx context.Context, userId int, itemId string) (res []entity.ClassItem, err error) {
	var (
		rows *sql.Rows
	)

	rows, err = q.DB.QueryContext(ctx, `select ci.item_id, ci.class_id from class_user as cu 
											inner join class_item as ci on cu.class_id = ci.class_id 
											where cu.user_id = ? and ci.item_id = ?`, userId, itemId)
	if err != nil {
		return
	}

	for rows.Next() {
		var item entity.ClassItem
		err = rows.Scan(&item.ItemID, &item.ClassID)
		if err != nil {
			return
		}

		res = append(res, item)
	}

	return
}

func (q *Querier) QueryWarningOrderWithFutureDate(ctx context.Context) (res []entity.WarningOrderWithFutureDate, err error) {
	var (
		rows *sql.Rows
	)

	now := time.Now().String()[:10] + " 23:59:59"
	rows, err = q.DB.QueryContext(ctx, `select id, note, back_note, created_user_id, date_created from orders where date_created > ? and active = 1 order by date_created`, now)
	if err != nil {
		return
	}

	for rows.Next() {
		var (
			order        entity.WarningOrderWithFutureDate
			date_created string
		)
		err = rows.Scan(&order.OrderID, &order.Note, &order.BackNote, &order.CreatedUserID, &date_created)
		if err != nil {
			return
		}

		order.DateCreated = helper.ConvertMysqlTimeUnixTime(date_created)
		res = append(res, order)
	}
	return
}

func (q *Querier) QueryAmountOfDuplicatedOrders(ctx context.Context, date string) (res []entity.UserItem, err error) {
	var (
		rows *sql.Rows
	)
	splitDate := strings.Split(date, "-")
	year, err := strconv.Atoi(splitDate[0])
	if err != nil {
		return
	}
	month, err := strconv.Atoi(splitDate[1])
	if err != nil {
		return
	}

	firstDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	nextMonth := firstDate.AddDate(0, 1, 0)
	rows, err = q.DB.QueryContext(ctx, `SELECT item_id, user_id, count(*) as total 
											FROM user_item as ui INNER JOIN user_item_order as uio ON ui.id = uio.user_item_id 
											AND uio.active = 1 
											WHERE uio.date_created >= ? AND uio.date_created < ?
											GROUP BY item_id, user_id 
											ORDER BY total DESC;`, firstDate, nextMonth)

	for rows.Next() {
		var (
			userItem entity.UserItem
			total    int
		)
		err = rows.Scan(&userItem.ItemID, &userItem.UserID, &total)
		if err != nil {
			return
		}

		if total > 1 {
			res = append(res, userItem)
		}
	}
	return
}

func (q *Querier) QueryDuplicatedOrders(ctx context.Context, date string, userId int, itemId string) (res []entity.Order, err error) {
	var (
		rows *sql.Rows
	)

	splitDate := strings.Split(date, "-")
	year, err := strconv.Atoi(splitDate[0])
	if err != nil {
		return
	}
	month, err := strconv.Atoi(splitDate[1])
	if err != nil {
		return
	}

	firstDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	nextMonth := firstDate.AddDate(0, 1, 0)
	rows, err = q.DB.QueryContext(ctx, `SELECT o.id, o.amount, o.discount, o.back_note, o.created_user_id, o.date_created FROM user_item as ui INNER JOIN user_item_order as uio 
											ON ui.id = uio.user_item_id AND ui.active = 1 
											INNER JOIN orders as o ON o.id = uio.order_id AND o.active = 1
											WHERE uio.date_created >= ? AND uio.date_created < ?
											AND ui.item_id = ?
											AND ui.user_id = ? ORDER BY o.date_created DESC;`, firstDate, nextMonth, itemId, userId)
	if err != nil {
		return
	}

	for rows.Next() {
		var (
			order        entity.Order
			date_created string
		)
		err = rows.Scan(&order.ID, &order.Amount, &order.Discount, &order.BackNote, &order.CreatedUserID, &date_created)
		if err != nil {
			return
		}

		order.DateCreated = helper.ConvertMysqlTimeUnixTime(date_created)
		res = append(res, order)
	}
	return
}

func (q *Querier) QueryDeletedDuplicatedOrders(ctx context.Context, date string) (res []entity.Order, err error) {
	var (
		rows *sql.Rows
	)

	splitDate := strings.Split(date, "-")
	year, err := strconv.Atoi(splitDate[0])
	if err != nil {
		return
	}
	month, err := strconv.Atoi(splitDate[1])
	if err != nil {
		return
	}

	firstDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	nextMonth := firstDate.AddDate(0, 1, 0)
	rows, err = q.DB.QueryContext(ctx, "SELECT id, amount, discount, note, back_note, created_user_id, reason_deleted, date_created, date_updated FROM orders WHERE active = 0 and date_created >= ? and date_created < ? and deleted_user_id = 1", firstDate, nextMonth)
	if err != nil {
		return
	}

	for rows.Next() {
		var (
			order                        entity.Order
			reason_deleted, date_updated sql.NullString
			date_created                 string
		)
		err = rows.Scan(&order.ID, &order.Amount, &order.Discount, &order.Note, &order.BackNote, &order.CreatedUserID, &reason_deleted, &date_created, &date_updated)
		if err != nil {
			return
		}
		if reason_deleted.Valid {
			order.ReasonDeleted = reason_deleted.String
		}
		if date_updated.Valid {
			order.DateUpdated = helper.ConvertMysqlTimeUnixTime(date_updated.String)
		}
		order.DateCreated = helper.ConvertMysqlTimeUnixTime(date_created)

		res = append(res, order)
	}

	return
}

// *** Items ***
// *** Invoices ***
func (q *Querier) QueryAmountDuplicatedInvoices(ctx context.Context, date string) (res []entity.Invoice, err error) {
	var (
		rows *sql.Rows
	)

	splitDate := strings.Split(date, "-")
	year, err := strconv.Atoi(splitDate[0])
	if err != nil {
		return
	}
	month, err := strconv.Atoi(splitDate[1])
	if err != nil {
		return
	}

	firstDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	nextMonth := firstDate.AddDate(0, 1, 0)
	rows, err = q.DB.QueryContext(ctx, `select balance_id, order_id, title, status, balance_before_deposit, balance_after_deposit, note, manual_note, payment, transaction_id, clerk_id, count(*) as total
										from invoices WHERE date_created >= ? and date_created < ? and active = 1
										group by balance_id, order_id, title, status, balance_before_deposit, balance_after_deposit, note, manual_note, payment, transaction_id, clerk_id
										ORDER BY total desc;`, firstDate, nextMonth)
	if err != nil {
		return
	}

	for rows.Next() {
		var (
			invoice               entity.Invoice
			total                 int
			order_id, manual_note sql.NullString
		)
		err = rows.Scan(&invoice.BalanceID, &order_id, &invoice.Title, &invoice.Status, &invoice.BalanceBeforeDeposit, &invoice.BalanceAfterDeposit, &invoice.Note, &manual_note, &invoice.Payment, &invoice.TransactionID, &invoice.ClerkID, &total)
		if err != nil {
			return
		}

		if order_id.Valid {
			invoice.OrderID = order_id.String
		}

		if manual_note.Valid {
			invoice.ManualNote = manual_note.String
		}
		if total > 1 {
			res = append(res, invoice)
		}
	}
	return
}

func (q *Querier) QuerySpecificDuplicatedInvoices(ctx context.Context, invoice entity.Invoice) (res []entity.Invoice, err error) {
	var (
		rows *sql.Rows
	)

	rows, err = q.DB.QueryContext(ctx, `SELECT * FROM invoices WHERE balance_id = ? AND order_id = ? AND title = ? 
							AND status = ? AND balance_before_deposit = ? AND balance_after_deposit = ? AND note = ? AND manual_note = ? AND payment = ? AND transaction_id = ? AND clerk_id = ? AND active = 1
							ORDER BY date_created DESC;`,
		invoice.BalanceID, invoice.OrderID, invoice.Title, invoice.Status, invoice.BalanceBeforeDeposit, invoice.BalanceAfterDeposit, invoice.Note, invoice.ManualNote, invoice.Payment, invoice.TransactionID, invoice.ClerkID)

	for rows.Next() {
		var (
			invoice                                             entity.Invoice
			order_id, manual_note, reason_deleted, date_updated sql.NullString
			deleted_user_id                                     sql.NullInt32
			date_created                                        string
		)
		err = rows.Scan(&invoice.ID, &invoice.BalanceID, &order_id, &invoice.Title, &invoice.Status, &invoice.BalanceBeforeDeposit,
			&invoice.BalanceAfterDeposit, &invoice.Note, &manual_note, &invoice.Payment,
			&invoice.TransactionID, &invoice.ClerkID, &reason_deleted, &deleted_user_id, &invoice.Active, &date_created, &date_updated)
		if err != nil {
			return
		}

		if order_id.Valid {
			invoice.OrderID = order_id.String
		}

		if manual_note.Valid {
			invoice.ManualNote = manual_note.String
		}

		if reason_deleted.Valid {
			invoice.ReasonDeleted = reason_deleted.String
		}

		if deleted_user_id.Valid {
			invoice.DeletedUserID = int(deleted_user_id.Int32)
		}

		if date_updated.Valid {
			invoice.DateUpdated = helper.ConvertMysqlTimeUnixTime(date_updated.String)
		}

		invoice.DateCreated = helper.ConvertMysqlTimeUnixTime(date_created)
		res = append(res, invoice)
	}
	return
}

func (q *Querier) QueryDeletedDuplicatedInvoices(ctx context.Context, date string) (res []entity.Invoice, err error) {
	var (
		rows *sql.Rows
	)

	splitDate := strings.Split(date, "-")
	year, err := strconv.Atoi(splitDate[0])
	if err != nil {
		return
	}
	month, err := strconv.Atoi(splitDate[1])
	if err != nil {
		return
	}

	firstDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	nextMonth := firstDate.AddDate(0, 1, 0)
	log.Println(nextMonth)

	// rows, err = q.DB(ctx).QueryContext(ctx, )

	for rows.Next() {
		var (
			invoice                                             entity.Invoice
			order_id, manual_note, reason_deleted, date_updated sql.NullString
			deleted_user_id                                     sql.NullInt32
			date_created                                        string
		)
		err = rows.Scan(&invoice.ID, &invoice.BalanceID, &order_id, &invoice.Title, &invoice.Status, &invoice.BalanceBeforeDeposit,
			&invoice.BalanceAfterDeposit, &invoice.Note, &manual_note, &invoice.Payment,
			&invoice.TransactionID, &invoice.ClerkID, &reason_deleted, &deleted_user_id, &invoice.Active, &date_created, &date_updated)
		if err != nil {
			return
		}

		if order_id.Valid {
			invoice.OrderID = order_id.String
		}

		if manual_note.Valid {
			invoice.ManualNote = manual_note.String
		}

		if reason_deleted.Valid {
			invoice.ReasonDeleted = reason_deleted.String
		}

		if deleted_user_id.Valid {
			invoice.DeletedUserID = int(deleted_user_id.Int32)
		}

		if date_updated.Valid {
			invoice.DateUpdated = helper.ConvertMysqlTimeUnixTime(date_updated.String)
		}

		invoice.DateCreated = helper.ConvertMysqlTimeUnixTime(date_created)
		res = append(res, invoice)
	}
	return
}
