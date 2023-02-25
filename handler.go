package main

import (
	"database/sql"


	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Handler struct {
	db *sql.DB
}

func newHandler(db *sql.DB) *Handler {
	return &Handler{db}
}

func (h *Handler) getUsersHandler(c *gin.Context) {
	rows, err := h.db.Query("SELECT * FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.UserName, &user.Password, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.ProfileImg); err != nil {
			log.Println(err)
			continue
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := Response{
		Message: "success",
		Item:    users,
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) createUserHandler(c *gin.Context) {
	// Parse request body to get user data
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Execute the SQL query to insert a new user into the database
	result, err := h.db.Exec("INSERT INTO users (user_name, password, fname, lname, email, phone, img_profile) VALUES (?, ?, ?, ?, ?, ?, ?)",
		user.UserName, user.Password, user.FirstName, user.LastName, user.Email, user.Phone, user.ProfileImg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get the ID of the newly inserted user
	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Set the ID of the new user and return the user object as JSON
	user.ID = int(id)
	c.JSON(http.StatusCreated, user)
}

func (h *Handler) updateUserHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := h.db.Prepare("UPDATE users SET user_name=?, password=?, fname=?, lname=?, email=?, phone=?, img_profile=? WHERE userID=?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	res, err := stmt.Exec(user.UserName, user.Password, user.FirstName, user.LastName, user.Email, user.Phone, user.ProfileImg, user.ID)
	if err != nil {
		panic(err.Error())
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		panic(err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "success",
		"rowsAffected": rowsAffected,
	})
}



func (h *Handler) deleteUserHandler(c *gin.Context) {
    var query UserQuery
    if err := c.Bind(&query); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    result, err := h.db.Exec("DELETE FROM users WHERE userID = ?", query.UserID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"rowsAffected": rowsAffected})
}

