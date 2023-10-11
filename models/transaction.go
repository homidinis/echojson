package models

type Transaction struct {
	Transaction_id   string  `json:"transactionID" `
	Transaction_date *string `json:"transactionDate" ` //pointer allows null values (one null value will turn everything blank!!)
	ID               int     `json:"id" `
	User_id          int     `json:"userID"  `
	Payment_method   string  `json:"paymentMethod" `
}
type RequestTransaction struct { //cuma untuk bulk insert/update
	Request []Transaction `json:"request" validate:"required,dive"`
}
