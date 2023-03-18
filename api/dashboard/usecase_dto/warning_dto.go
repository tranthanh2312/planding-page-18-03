package usecase_dto

type IncompatibleOrder struct {
	UserID   int     `json:"userId"`
	FullName string  `json:"fullName"`
	ItemName string  `json:"itemName"`
	Price    float32 `json:"price"`
	Reason   string  `json:"reason"`
	Text     string  `json:"text"`
}

type WarningOrderWithFutureDate struct {
	OrderID             string `json:"orderId"`
	Note                string `json:"note"`
	ItemName            string `json:"itemName"`
	CreatedUserFullName string `json:"createdUserFullName"`
	DateCreated         string `json:"dateCreated"`
}

type WarningDuplicatedOrder struct {
	OrderID  string  `json:"orderId"`
	UserID   int     `json:"userId"`
	FullName string  `json:"fullName"`
	ItemName string  `json:"itemName"`
	Price    float32 `json:"price"`
	Reason   string  `json:"reason"`
}
