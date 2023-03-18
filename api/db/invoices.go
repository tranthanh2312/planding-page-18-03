package db

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
)

type Invoice struct {
	TotalCard     float32 `json:"totalCard"`
	TotalCash     float32 `json:"totalCash"`
	Total         float32 `json:"total"`
	TotalInvoices int     `json:"totalInvoices"`
}

func (q *Querier) QueryInvoices(date string) (res Invoice, err error) {
	var (
		rows *sql.Rows
	)

	rows, err = q.DB.Query(`select count(*), sum(-payment), transaction_id from defaultdb.invoices 
							where date_created >= ? and date_created <= ? and active = 1 and title != 'Nhận Tiền' group by transaction_id;`,
		date+" 00:00:00", date+" 23:59:59")
	if err != nil {
		return
	}

	for rows.Next() {
		var (
			total          sql.NullFloat64
			totalInvoices  int
			transaction_id sql.NullString
		)
		err = rows.Scan(&totalInvoices, &total, &transaction_id)
		if err != nil {
			return
		}
		if transaction_id.Valid {
			if transaction_id.String == "cash" {
				res.TotalCash += float32(total.Float64)
			} else if transaction_id.String == "card" {
				res.TotalCard += float32(total.Float64)
			}
		}

		res.TotalInvoices += totalInvoices
	}
	res.Total = res.TotalCard + res.TotalCash
	return
}

func (q *Querier) QueryTotalAmountInvoicesBySameDateInPreviousMonth(date string) (res Invoice, err error) {
	var (
		total sql.NullFloat64
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
	day, err := strconv.Atoi(splitDate[2])
	if err != nil {
		return
	}

	convertDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	previousMonth := convertDate.AddDate(0, -1, 0).String()[:10]

	err = q.DB.QueryRow(`select sum(-payment) from defaultdb.invoices 
	where date_created >= ? and date_created <= ? and active = 1 and title != 'Nhận Tiền' and transaction_id != 'auto_charge';`,
		previousMonth+" 00:00:00", previousMonth+" 23:59:59").Scan(&total)
	if err != nil {
		return
	}

	if total.Valid {
		res.Total = float32(total.Float64)
	}
	return
}

func (q *Querier) QueryTotalAmountInvoicesByPreviousDate(date string) (res Invoice, err error) {
	var (
		total sql.NullFloat64
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
	day, err := strconv.Atoi(splitDate[2])
	if err != nil {
		return
	}

	convertDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	previousDate := convertDate.AddDate(0, 0, -1).String()[:10]

	err = q.DB.QueryRow(`select sum(-payment) from defaultdb.invoices 
	where date_created >= ? and date_created <= ? and active = 1 and title != 'Nhận Tiền' and transaction_id != 'auto_charge';`,
		previousDate+" 00:00:00", previousDate+" 23:59:59").Scan(&total)
	if err != nil {
		return
	}

	if total.Valid {
		res.Total = float32(total.Float64)
	}
	return
}
