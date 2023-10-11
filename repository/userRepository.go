package repository

import (
	"database/sql"
	"echojson/config"
	"echojson/models"
	"echojson/utils"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

/*
=================================

GET USERS (SELECT)
todo: add "if id=0"
==================================
*/
func GetUser(id int) (userContainer []models.User) { //return userContainer, yang map ke response di usecase
	db := config.Conn()
	query := `SELECT id, age, first_name, last_name, email, username, groups FROM public.users`
	var data []interface{}

	if id != 0 { //if id is not presented
		query += "WHERE id=$1"  //append "where" to query
		data = append(data, id) //then append the id arg to the data interface. in the case of there being a lot of arguments for a lot of WHERE conditions
	}

	fmt.Println("query: " + query)
	rows, err := db.Query(query, data...) //append data (lots of them, potentially; ... is to pass multiple values, like an array)
	if err != nil {
		fmt.Println("err in getuser query: ")
		fmt.Println(err)
	}
	for rows.Next() {
		var users models.User
		err := rows.Scan(&users.ID, &users.Age, &users.First_name, &users.Last_name, &users.Email, &users.Username, &users.Admin)
		if err != nil {
			fmt.Println("err in getuser query: ")
			fmt.Println(err)
		}
		userContainer = append(userContainer, users)
	}
	return
}

/*========================================

ADD users

========================================*/

func AddUser(users models.RequestUser, user int, tx *sql.Tx) (vals []interface{}, err error) {

	sqlStr := `INSERT INTO public.users(age, first_name, last_name, email, username, password, "group") VALUES `

	for _, row := range users.Request { //index,name_of_
		sqlStr += " (?, ?, ?, ?, ?, ?, ?),"
		vals = append(vals, row.Age, row.First_name, row.Last_name, row.Email, row.Username, row.Password, row.Admin)
	}
	// trim the last ,
	sqlStr = strings.TrimSuffix(sqlStr, ",")
	// replacing ? with $n for postgres
	sqlStr = utils.ReplaceSQL(sqlStr, "?")

	//format all vals at once
	_, err = tx.Exec(sqlStr, vals...)
	if err != nil {
		fmt.Println("exec Error:", err)
		return
	}
	return
}

func Register(users models.RequestUser, tx *sql.Tx) (vals []interface{}, err error) {

	sqlStr := `INSERT INTO public.users(age, first_name, last_name, email, username, password, admin) VALUES `

	for _, row := range users.Request { //index,name_of_
		sqlStr += " (?, ?, ?, ?, ?, ?, ?),"
		bpass, err := bcrypt.GenerateFromPassword([]byte(row.Password), 10)
		if err != nil {
			return nil, err
		}
		vals = append(vals, row.Age, row.First_name, row.Last_name, row.Email, row.Username, bpass, false)
	}
	// trim the last ,
	sqlStr = strings.TrimSuffix(sqlStr, ",")
	// replacing ? with $n for postgres
	sqlStr = utils.ReplaceSQL(sqlStr, "?")

	//format all vals at once
	_, err = tx.Exec(sqlStr, vals...)
	if err != nil {
		fmt.Println("exec Error:", err)
		return
	}
	return
}

/*
================================

# UPDATE USERS
this part updates users
====================================
*/
func UpdateUser(userContainer models.User, user int, tx *sql.Tx) (updated_id int, err error) {

	query := `UPDATE 
									public.users 
									SET 
									age=$1, first_name=$2, last_name=$3, email=$4, username=$5, 
									password=$6, "group"=$7 
									WHERE 
									id=$8 
									RETURNING id;`
	if err != nil {
		fmt.Println("Prep Error:", err)
		return
	}
	var users models.User //for Scan() container
	err = tx.QueryRow(query,
		&userContainer.Age,
		&userContainer.First_name,
		&userContainer.Last_name,
		&userContainer.Email,
		&userContainer.Username,
		&userContainer.Password,
		&userContainer.Admin,
		&userContainer.ID).Scan(&users.ID)
	if err != nil {
		fmt.Println("Exec Error:", err)
		return
	}
	updated_id = users.ID
	return
}

/*
================================

# DELETE USERS

====================================
*/
func DeleteUser(userContainer models.User, user int, tx *sql.Tx) (users models.User, err error) {
	query := `DELETE FROM public.users WHERE id=$1 RETURNING id`
	if err != nil {
		fmt.Println("Prep Error:", err)
		return
	}
	err = tx.QueryRow(query, &userContainer.ID).Scan(&users.ID)
	if err != nil {
		fmt.Println("Exec Error:", err) //response pindah ke usecase, jadi di controller cuma get value dari query;
		return
	}
	return
}
