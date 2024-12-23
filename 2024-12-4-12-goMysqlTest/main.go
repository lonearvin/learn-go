package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql" // 导入 MySQL 驱动程序
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

// 模块初始化
func init() {
	database, err := sqlx.Open("mysql", "root:196618@tcp(127.0.0.1:3308)/test")
	if err != nil {
		fmt.Println("open mysql failed, ", err)
		return
	}
	Db = database
}

func main() {
	r, err := Db.Exec("insert into person(username, sex, email)values(?, ?, ?)", "stu001", "man", "stu01@qq.com")
	if err != nil {
		fmt.Println("exec error,", err)
		return
	}
	id, err := r.LastInsertId()
	if err != nil {
		fmt.Println("exec failed,", err)
		return
	}
	fmt.Println("insert successful:", id)
}

// 所学习到的知识如下：
/*
为什么使用 _ 来导入包？
当你在代码中仅仅是需要包的副作用（如注册数据库驱动），而不直接引用包中的函数或类型时，可以使用 _ 来导入该包。这种做法常见于以下几种情况：

数据库驱动注册： 在 Go 中，数据库驱动需要注册到 database/sql 包中，通过导入驱动并调用其 init() 函数来完成初始化。使用 _ 来导入驱动是非常常见的做法。
*/
