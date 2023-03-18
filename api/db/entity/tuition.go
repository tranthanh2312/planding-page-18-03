package entity

type Balance struct {
	ID          string
	UserID      int
	Balance     float32
	Active      int
	DateCreated int64
	DateUpdated int64
}

type BalanceWithUser struct {
	BalanceID  string
	UserID     int
	Balance    float32
	FullName   string
	Username   string
	Phone      string
	EntityCode int
	Dob        int64 `db:"dob"`
}

type Invoice struct {
	ID                   string
	BalanceID            string
	OrderID              string
	Title                string
	Status               bool
	BalanceBeforeDeposit float32
	BalanceAfterDeposit  float32
	Note                 string
	ManualNote           string
	Payment              float32
	ClerkID              int
	ReasonDeleted        string
	DeletedUserID        int
	TransactionID        string // FK

	Active      int
	DateCreated int64
	DateUpdated int64
}

type InvoiceDeleted struct {
	ID        string
	BalanceID string
	Title     string
	ClerkName string
}

type Item struct {
	ID            string  `db:"id"`
	ItemName      string  `db:"item_name"`
	Price         float32 `db:"price"`
	Description   string  `db:"description"`
	TypeId        string  `db:"type_id"`
	RecurringDay  int     `db:"recurring_day"`
	CreatedUserID int     `db:"created_user_id"`
	DeletedUserID int     `db:"deleted_user_id"`
	Active        int     `db:"active"`
	DateCreated   int64   `db:"date_created"`
	DateUpdated   int64   `db:"date_updated"`
}

type Order struct {
	ID                     string `json:"orderId"`
	Status                 bool
	Amount                 float32
	Discount               float32
	ReasonDiscount         string
	AssignDiscountPersonID int
	Note                   string
	BackNote               string
	CreatedUserID          int
	UpdatedUserID          int
	MessageSent            int
	ReasonDeleted          string
	DeletedUserID          int
	Active                 int
	DateCreated            int64
	DateUpdated            int64
}

type Debt struct {
	Order Order
	User  User
}

type ClassDebt struct {
	UserID int
	Debt   float32
}

type UserItem struct {
	ID                        string
	ItemID                    string
	UserID                    int
	Status                    bool
	Quantity                  int
	ReasonAssigningDiscount   string
	PersonIDAssigningDiscount int
	DiscountPercent           float32
	ReasonDeleted             string
	DeletedUserID             int
	Active                    int
	DateCreated               int64
	DateUpdated               int64
}

type UserItemOrder struct {
	UserItemID  string
	OrderID     string
	DateCreated int64
	DateUpdated int64
	Active      int
}

type Transaction struct {
	ID              string
	TransactionName string
	Description     string
}

type BalanceHistory struct {
	BalanceID            string
	Note                 string
	CreatedUserSignature string
	DepositorSignature   string
	DepositorFullname    string
	DepositeAmout        int
	CreatedUserID        int
	Active               int
	DateCreated          int64
	DateUpdated          int64
}

type TransactionType struct {
	ID              string
	TransactionName string
	Desription      string
}

type UserItemCombo struct {
	Item     Item
	User     User
	UserItem UserItem
}

type ItemPriceHistory struct {
	ID            int
	ItemID        string
	Price         float32
	DateEffective int64
	DateDefective int64
	CreatedUserID int
}

type TuitionActivity struct {
	ID              int
	CreatedEntities int
	DateCreated     int64
}
