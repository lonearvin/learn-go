package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Person struct {
	UserId   int    `db:"user_id"`
	Username string `db:"username"`
	Sex      string `db:"sex"`
	Email    string `db:"email"`
}

type Place struct {
	Country string `db:"country"`
	City    string `db:"city"`
	TelCode int    `db:"telcode"`
}

var Db *sqlx.DB

func init() {

	database, err := sqlx.Open("mysql", "root:196618@tcp(10.16.64.250:3308)/test")
	if err != nil {
		fmt.Println("open mysql failed,", err)
		return
	}

	Db = database
}

func main() {

	var person []Person
	err := Db.Select(&person, "select * from person where user_id=?", 2)
	if err != nil {
		fmt.Println("exec failed, ", err)
		return
	}

	for _, p := range person {
		//fmt.Println(p)
		fmt.Printf("id:%v\nname:%v\nsex:%v\nemail:%v\n", p.UserId, p.Username, p.Sex, p.Email)
	}

}
