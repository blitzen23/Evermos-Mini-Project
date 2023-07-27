package main

import (
	"project/routing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)


func main() {
	routing.Start()
}