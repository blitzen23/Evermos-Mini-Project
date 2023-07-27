package controller

import (
	"log"
	"project/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Store *session.Store
var DB *gorm.DB
var err error

func Start(){
	Store = session.New()
	DB, err = gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/Store?charset=utf8&parseTime=True")

	if err != nil{
		log.Println("Connection Failed to Open")
	}else{ 
		log.Println("Connection Established")
	}

	DB.SingularTable(true)
	DB.AutoMigrate(&model.Category{}, &model.User{}, &model.Alamat{}, &model.Trx{}, &model.Toko{}, &model.Detail_Trx{}, &model.Produk{}, &model.Foto_Produk{}, &model.Log_Produk{})
}