package db

type User struct {
	ID       int
	FullName string
	Phone    string
}

func GetFullName(userId int, users []User) string {
	// Binary Search
	l, r := 0, len(users)-1
	for l <= r {
		m := (l + r) / 2
		if users[m].ID == userId {
			return users[m].FullName
		}
		if users[m].ID < userId {
			l = m + 1
		} else {
			r = m - 1
		}
	}
	return ""
}

// Cache user's data
func (q *Querier) GetAllUsers() (users []User, err error) {
	// If not, query from database
	rows, err := q.DB.Query("SELECT id, full_name, phone FROM users WHERE active = 1 ORDER BY id")
	if err != nil {
		return
	}

	for rows.Next() {
		var user User
		err = rows.Scan(&user.ID, &user.FullName, &user.Phone)
		if err != nil {
			return
		}
		users = append(users, user)
	}

	return users, err
}
