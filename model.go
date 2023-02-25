package main

import (
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID         int    `json:"userID"`
	UserName   string `json:"user_name"`
	Password   string `json:"password"`
	FirstName  string `json:"fname"`
	LastName   string `json:"lname"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	ProfileImg string `json:"img_profile"`
}

type Response struct {
	Message string `json:"message"`
	Item    []User `json:"item"`
}
type ResponseUnauthorized struct {
	Message string `json:"message"`
}
type UserQuery struct {
	UserID int `json:"userID"`
}
