package models

type User struct {
	ID         int    `json:"id"`
	First_name string `json:"firstName" validate:"required"`
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"required"`
	Age        int    `json:"age" validate:"required,gte=0,lte=80"`
	Last_name  string `json:"lastName" validate:"required"`
	Admin      *bool  `json:"isAdmin" `
	Email      string `json:"email" validate:"required,email"`
}
type RequestUser struct {
	Request []User `json:"request" validate:"required,dive"`
}

type UserID struct {
	User_id        int    `json:"userid"`
	Payment_method string `json:"payment_method"`
}
