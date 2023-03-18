package db

import "database/sql"

type Class struct {
	ID        string
	ClassName string
	RoomCode  string
}

func (q *Querier) GetAllClasses() (classes []Class, err error) {
	bolt, err := Open()
	if err != nil {
		return
	}

	defer bolt.Close()

	err = bolt.Get("classes", classes)
	if err != nil {
		return
	}

	if err == nil && classes != nil {
		return
	}

	var rows *sql.Rows
	rows, err = q.DB.Query("SELECT id, class_name, room_code FROM classes WHERE active = 1")
	if err != nil {
		return
	}

	for rows.Next() {
		var class Class
		err = rows.Scan(&class.ID, &class.ClassName, &class.RoomCode)
		if err != nil {
			return
		}

		classes = append(classes, class)
	}

	err = bolt.Set("classes", classes) // set classes into cache
	return
}
