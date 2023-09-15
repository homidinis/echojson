package repository

import (
	"echojson/db"
	"echojson/models"
	"fmt"
)

/*
==================================

# LOGIN

==================================
*/
func GetLogin(username string) (containerArray []models.User, err error) { //declare as array of User struct
	db := db.Conn()
	//1. grabs username and password WHERE typed-username
	//2. dumps username, password, firstname from database into container
	//3. compare password in container to typed-password
	//4. if no errors, generate token
	statement, err := db.Prepare("SELECT id, username, password, first_name, admin FROM users WHERE username=$1") //only select by Username
	if err != nil {
		fmt.Println("Prepare err in controller:", err)
	}
	var container models.User
	err = statement.QueryRow(username).Scan(&container.ID, &container.Username, &container.Password, &container.First_name, &container.Admin) //container = [user{Username,Password,Firstname}]; scan scans into each of them
	if err != nil {
		fmt.Println("Query err in controller:", err)
	}
	containerArray = append(containerArray, container)
	return
}
