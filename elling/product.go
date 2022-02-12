package elling

type Product struct {
	ID          int64
	DisplayName string
	Billing     string
	BalanceID   int64   `json:"-"`
	Balance     Balance `json:"balance"`
	Module      string
	Type        string
	UserID      int64
	Suspended   bool
}

type ProductOrderedEvent struct {
	User    *User
	Product *Product
}

type ProductSuspendEvent struct {
	User    *User
	Product *Product
}

type ProductDeletedEvent struct {
	User    *User
	Product *Product
}
