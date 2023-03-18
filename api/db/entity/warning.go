package entity

type IncompatibleOrder struct {
	UserID   int
	ItemID   string
	Amount   float32
	Discount float32
}

type WarningOrderWithFutureDate struct {
	OrderID       string
	Note          string
	BackNote      string
	CreatedUserID int
	DateCreated   int64
}
