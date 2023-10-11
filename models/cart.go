package models

type RequestCart struct {
	Request []Cart `json:"request" validate:"required,dive"`
}

type Cart struct {
	ID         int `json:"id"`
	User_id    int `json:"userID" validate:"required"`
	Product_id int `json:"productID" validate:"required"`
	Quantity   int `json:"quantity" validate:"required,gte=0"`
	Price      int `json:"price" validate:"required"`
}
type PaymentMethodCart struct {
	Cart           Cart   `json:"cart"`
	Payment_method string `json:"payment_method"`
}
