package elling

type Product struct {
	ID          int64
	DisplayName string
	Billing     string
	Balance     Balance
	Module      string
	Type        string
}