package entity

type User struct {
	ID            int     `db:"id"`
	FullName      string  `db:"fullname"`
	Username      string  `db:"username"`
	Password      string  `db:"password"`
	Gender        string  `db:"gender"`
	Address       string  `db:"address"`
	Mail          string  `db:"mail"`
	Phone         string  `db:"phone"`
	Dob           int64   `db:"dob"`
	Qualification string  `db:"qualification"`
	EntityCode    int     `db:"entity_code"`
	Elo           float32 `db:"elo"`
	Active        int     `db:"active"`

	DateCreated int64 `db:"datecreated"`
	DateUpdated int64 `db:"dateupdated"`
}
