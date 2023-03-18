package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Tech-by-GL/dashboard/db/entity"
	"github.com/Tech-by-GL/dashboard/helper"
)

type Item struct {
	ID       string
	ItemName string
	TypeID   string
}

type AllItem struct {
	Items []Item
}

func (a *AllItem) FindItemName(itemID string) (ItemName string, TypeID string) {
	for _, data := range a.Items {
		if data.ID == itemID {
			return data.ItemName, data.TypeID
		}
	}
	return "", "-1"
}

func (q *Querier) FindItems(itemsID []string) (res []Item, err error) {
	var rows *sql.Rows

	valueString := make([]string, 0, len(itemsID))
	valueArgs := make([]interface{}, 0, len(itemsID))
	for i := 0; i < len(itemsID); i++ {
		valueString = append(valueString, "?")
		valueArgs = append(valueArgs, itemsID[i])
	}

	stmt := fmt.Sprintf("SELECT id, item_name, type_id FROM items WHERE id in (%s)", strings.Join(valueString, ","))
	rows, err = q.DB.Query(stmt, valueArgs...) // ! Because the input is existed so that's why we dont have to provide more where conditions to find is item existed or not.

	for rows.Next() {
		var item Item
		err = rows.Scan(&item.ID, &item.ItemName, &item.TypeID)
		if err != nil {
			return
		}

		res = append(res, item)
	}
	return
}

func (q *Querier) GetAllItems() (res AllItem, err error) {
	var (
		rows *sql.Rows
		bolt *Cache
	)

	bolt, err = Open()
	if err != nil {
		return
	}

	defer bolt.Close()

	err = bolt.Get("items", &res.Items)
	if err != nil {
		return
	}
	if err == nil && res.Items != nil {
		return
	}

	rows, err = q.DB.Query("SELECT id, item_name, type_id FROM items")
	if err != nil {
		return
	}

	for rows.Next() {
		var item Item
		err = rows.Scan(&item.ID, &item.ItemName, &item.TypeID)
		if err != nil {
			return
		}

		res.Items = append(res.Items, item)
	}

	err = bolt.Set("items", res.Items)
	return
}

// flag = 1: itemId; 2 = all, 3 names
func (q *Querier) QueryItem(ctx context.Context, itemId string, itemName string, flag int) (allItems []entity.Item, err error) {
	var (
		date_updated  sql.NullString
		deletedUserId sql.NullInt32
		rows          *sql.Rows
		data          entity.Item
		date_created  string
	)

	switch flag {
	case 1:
		rows, err = q.DB.QueryContext(ctx, `SELECT * FROM items WHERE id = ? AND active = 1`, itemId)
	case 2:
		rows, err = q.DB.QueryContext(ctx, `SELECT * FROM items WHERE active = 1`)
	case 3:
		rows, err = q.DB.QueryContext(ctx, `SELECT * FROM items WHERE item_name = ? AND active = 1`, itemName)
	}

	if err == sql.ErrNoRows {
		return allItems, nil
	}

	if err != nil {
		return allItems, err
	}

	for rows.Next() {
		err = rows.Scan(&data.ID, &data.ItemName, &data.Price, &data.Description, &data.TypeId, &data.RecurringDay, &data.CreatedUserID, &deletedUserId, &data.Active, &date_created, &date_updated)
		if err != nil {
			return
		}

		data.DateCreated = helper.ConvertMysqlTimeUnixTime(date_created)
		if date_updated.Valid {
			data.DateUpdated = helper.ConvertMysqlTimeUnixTime(date_updated.String)
		}

		if deletedUserId.Valid {
			data.DeletedUserID = int(deletedUserId.Int32)
		}

		allItems = append(allItems, data)
	}

	return allItems, nil
}
