package data

import (
	"time"
)

type Book struct {
	BookID          uint       `json:"book_id"`
	Name            string     `json:"name"`
	DateToReturn    *time.Time `json:"date_to_return"`
	UserID          uint       `json:"user_id"`
	ShippingAddress string     `json:"shipping_address"`
}
