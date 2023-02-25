package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/my_pos")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	handler := newHandler(db)

	r := gin.Default()
	r.POST("/token", handler.getTokenHandler)
	protected := r.Group("/")
	protected.Use(handler.authorizationMiddleware)
	protected.GET("/users", handler.getUsersHandler)
	protected.POST("/users", handler.createUserHandler)
	protected.POST("/updateusers", handler.updateUserHandler)
	protected.POST("/delete", handler.deleteUserHandler)

	if err := r.Run(); err != nil {
		log.Fatal(err)
	}
}
