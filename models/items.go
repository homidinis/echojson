package models

type RequestItem struct {
	Request []Item `json:"request" validate:"required,dive"`
}

type Item struct {
	Product_id  int     `json:"product_id"`
	Name        string  `json:"item"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"min=0"`
	Quantity    int     `json:"quantity" validate:"min=0"`
}
