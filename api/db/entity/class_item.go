package entity

type ClassItem struct {
	ClassItemID   string
	ItemID        string
	ClassID       string
	Buying        bool
	Active        int
	CreatedUserID int
	DeletedUserID int
	DateCreated   int64
	DateUpdated   int64
}
