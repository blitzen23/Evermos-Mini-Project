package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/golang-jwt/jwt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/bcrypt"
)

type Category struct {
	ID uint `gorm:"primaryKey" json:"id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
	NamaCategory string `gorm:"type:varchar(255)" json:"nama_category"`
	Produks []Produk `gorm:"foreignKey:CategoryID" json:"produks"`
	Log_Produks []Log_Produk `gorm:"foreignKey:CategoryID" json:"log_produk"`
}

type User struct{
	ID uint `gorm:"primaryKey" json:"id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
	Nama string `gorm:"type:varchar(255)" json:"nama"`
	KataSandi string `gorm:"type:varchar(255)" json:"kata_sandi"`
	NoTelp string `gorm:"type:varchar(255)" json:"no_telp" gorm:"unique:not null"`
	TanggalLahir time.Time `json:"tanggal_lahir"`
	JenisKelamin string `gorm:"type:varchar(255);" json:"jenis_kelamin"`
	Tentang string `gorm:"type:text;" json:"tentang"`
	Pekerjaan string `gorm:"type:varchar(255)" json:"pekerjaan"`
	Email string `gorm:"type:varchar(255)" json:"email"`
	IdProvinsi string `gorm:"type:varchar(255)" json:"id_provinsi"`
	IdKota string `gorm:"type:varchar(255)" json:"id_kota"`
	IsAdmin bool `json:"isAdmin"`
	Alamats []Alamat `gorm:"foreignKey:UserID" json:"alamats"`
	Toko *Toko `gorm:"foreignKey:UserID" json:"toko"`
	Trx *Trx `gorm:"foreignKey:UserID" json:"trx"`
}

type Alamat struct {
	ID uint `gorm:"primaryKey" json:"id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
	UserID uint `json:"user_id"`
	JudulAlamat string `gorm:"type:varchar(255)" json:"judul_alamat"`
	NamaPenerima string `gorm:"type:varchar(255)" json:"nama_penerima"`
	NoTelp string `gorm:"type:varchar(255)" json:"no_telp" gorm:"unique:not null"`
	DetailAlamat string `gorm:"type:varchar(255)" json:"detail_alamat"`
	User User `json:"user"`
	Trx *Trx `gorm:"foreignKey:AlamatID" json:"trx"`
}

type Trx struct {
	ID uint `gorm:"primaryKey" json:"id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
	UserID uint `json:"user_id"`
	AlamatID uint `json:"alamat_id"`
	HargaTotal int `json:"harga_total"`
	KodeInvoice string `gorm:"type:varchar(255)" json:"kode_invoice"`
	MethodBayar string `gorm:"type:varchar(255)" json:"method_bayar"`
	User *User `json:"user"`
	Alamat *Alamat `json:"alamat"`
	Detail_Trx []Detail_Trx `json:"detail_trx"`
}

type Toko struct {
	ID uint `gorm:"primaryKey" json:"id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
	UserID uint `gorm:"autoIncrement:1;autoIncrementIncrement:1" json:"user_id"`
	NamaToko string `gorm:"type:varchar(255)" json:"nama_toko"`
	UrlFoto string `gorm:"type:varchar(255)" json:"url_foto"`
	User *User `json:"user"`
	Detail_Trx []Detail_Trx `gorm:"foreignKey:TokoID" json:"detail_trx"`
	Log_Produk []Log_Produk `gorm:"foreignKey:TokoID" json:"log_produk"`
	Produks []Produk `gorm:"foreignKey:TokoID" json:"produks"`
}

type Detail_Trx struct {
	ID uint `gorm:"primaryKey" json:"id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
	Kuantitas int `json:"kuantitas"`
	HargaTotal int `json:"harga_total"`
	TrxID uint `json:"trx_id"`
	LogProdukID uint `json:"log_produk_id"`
	TokoID uint `json:"toko_id"`
	Trx *Trx `json:"trx"`
	Toko *Toko `json:"toko"`
	Log_Produk *Log_Produk `json:"log_produk"`
}

type Produk struct {
	ID uint `gorm:"primaryKey" json:"id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
	NamaProduk string `gorm:"type:varchar(255)" json:"nama_produk"`
	Slug string `gorm:"type:varchar(255)" json:"slug"`
	HargaReseller int `json:"harga_reseller"`
	HargaKonsumen int `json:"harga_konsumen"`
	Stok int `json:"stok"`
	Deskripsi string `gorm:"type:text" json:"deskripsi"`
	TokoID uint `json:"toko_id"`
	CategoryID uint `json:"category_id"`
	Toko *Toko `json:"toko"`
	Category *Category `json:"category"`
	Foto_Produks []Foto_Produk `gorm:"foreignKey:ProdukID" json:"foto_produks"`
	Log_Produk *Log_Produk `json:"log_produk"`
}

type Foto_Produk struct {
	ID uint `gorm:"primaryKey" json:"id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
	Url string `gorm:"type:varchar(255)" json:"url"`
	ProdukID uint `json:"produk_id"`
	Produk Produk `json:"produk"`
}

type Log_Produk struct {
	ID uint `gorm:"primaryKey" json:"id"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
	NamaProduk string `gorm:"type:varchar(255)" json:"nama_produk"`
	Slug string `gorm:"type:varchar(255)" json:"slug"`
	HargaReseller int `json:"harga_reseller"`
	HargaKonsumen int `json:"harga_konsumen"`
	Deskripsi string `gorm:"type:text" json:"deskripsi"`
	TokoID uint `json:"toko_id"`
	CategoryID uint `json:"category_id"`
	ProdukID uint `json:"produk_id"`
	Toko *Toko `json:"toko"`
	Category *Category `json:"category"`
	Produk *Produk `json:"Produk"`
	Detail_Trx []Detail_Trx `gorm:"foreignKey:LogProdukID" json:"detail_trx"`
	Foto_Produks []Foto_Produk `gorm:"foreignKey:ProdukID;references:ProdukID" json:"foto_produks"`
}

var db *gorm.DB
var err error
var store *session.Store

func isAdmin(c *fiber.Ctx) error {
	tokenString := c.GetReqHeaders()["Token"]
	if tokenString == ""{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false,
			"message": "Failed to Access",
			"errors": "You are unauthorized",
			"data": nil,
		})
	} 
	token, err := jwt.Parse(tokenString, func(token *jwt.Token)(interface{}, error) {
		return []byte("secretkey"), nil
	})
	if err != nil || !token.Valid{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to Access",
			"errors": err.Error(),
			"data": nil,
		})
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to Access",
			"errors": "Invalid token claims",
			"data": nil,
		})
	}
	isAdmin, exists := claims["isAdmin"].(bool)
	if !exists || !isAdmin {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false,
			"message": "Failed to Access",
			"errors": "You are unauthorized",
			"data": nil,
		})
	}
	return c.Next()
}

func isLoggedIn(c *fiber.Ctx) error {
	tokenString := c.GetReqHeaders()["Token"]
	if tokenString == ""{
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false,
			"message": "Failed to Access",
			"errors": "You are unauthorized",
			"data": nil,
		})
	} 
	token, err := jwt.Parse(tokenString, func(token *jwt.Token)(interface{}, error) {
		return []byte("secretkey"), nil
	})
	if err != nil || !token.Valid{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to Access",
			"errors": err.Error(),
			"data": nil,
		})
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to Access",
			"errors": "Invalid token claims",
			"data": nil,
		})
	}
	sess, err := store.Get(c)
	if (err != nil){
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to Access",
			"errors": err.Error(),
			"data": nil,
		})
	}
	sess.Set("id", uint(claims["id"].(float64)))
	sess.SetExpiry(time.Hour)
	if err := sess.Save(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to Access",
			"errors": err.Error(),
			"data": nil,
		})
	}
	return c.Next()
}

type TempCategory struct {
	ID uint `gorm:"primaryKey" json:"id"`
	NamaCategory string `gorm:"type:varchar(255)" json:"nama_category"`
}

func main() {
		store = session.New()
    db, err = gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/Store?charset=utf8&parseTime=True")

    if err != nil{
    	log.Println("Connection Failed to Open")
    }else{ 
    	log.Println("Connection Established")
    }

		db.SingularTable(true)
		db.AutoMigrate(&Category{}, &User{}, &Alamat{}, &Trx{}, &Toko{}, &Detail_Trx{}, &Produk{}, &Foto_Produk{}, &Log_Produk{})
		app := fiber.New()
		app.Get("/test", func(c *fiber.Ctx) error {
			alamat := new(Alamat)
			err := db.Preload("User").First(&alamat, 1).Error
			if err != nil {
				return c.JSON("error")
			}
			return c.JSON(alamat)
			// user := new(User)
			// db.Preload("Alamats").First(&user, 1)
			// return c.JSON(user)
		})
    app.Get("/category", func(c *fiber.Ctx) error {
        categories := []Category{}
				err := db.Select("id, nama_category").Find(&categories).Error
				if (err != nil && err.Error() == "record not found"){
					return c.Status(fiber.StatusOK).JSON(fiber.Map{
						"status": true,
						"message": "Succeed to GET data",
						"errors": nil,
						"data": categories,
					})
				} else if (err != nil){
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"status": false,
						"message": "Failed to GET data",
						"errors": "There is an internal server error",
						"data": nil,
					})
				} else {
					tempCategories := []TempCategory{}
					for _, category := range categories{
						tempCategories = append(tempCategories, TempCategory{
							ID: category.ID,
							NamaCategory: category.NamaCategory,
						})
					}
					return c.Status(fiber.StatusOK).JSON(fiber.Map{
						"status": true,
						"message": "Succeed to GET data",
						"errors": nil,
						"data": tempCategories,
					})
				}
    })

		app.Get("/category/:id", isAdmin, func(c *fiber.Ctx) error {
				id, _ := strconv.Atoi(c.Params("id"))
				category := new(Category)
				err := db.Select("id, nama_category").First(&category, id).Error
				if (err != nil && err.Error() == "record not found"){
					return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
						"status": false,
						"message": "Failed to GET data",
						"errors": "No Data Category",
						"data": nil,
					})
				} else if (err != nil){
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"status": false,
						"message": "Failed to GET data",
						"errors": "There is an internal server error",
						"data": nil,
					})
				}
				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"status": true,
					"message": "Succeed to GET data",
					"errors": nil,
					"data": fiber.Map{
						"id": category.ID,
						"nama_category": category.NamaCategory,
					},
				})
		})

		app.Post("/category", isAdmin, func(c *fiber.Ctx) error {
			category := new(Category)

			if err := c.BodyParser(category); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": false,
					"message": "Failed to POST data",
					"errors": "Something bad behind the scene happens",
					"data": nil,
				})
			}
			if err := db.Create(&category).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"status": false,
						"message": "Failed to POST data",
						"errors": "There is an internal server error",
						"data": nil,
					})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"message": "Succeed to POST data",
				"errors": nil,
				"data": 1,
			})
		})
		
		app.Put("/category/:id", isAdmin, func(c *fiber.Ctx) error {
			id, _ := strconv.Atoi(c.Params("id"))
			category := new(Category)
			if err := db.Select("id, nama_category").First(&category, id).Error; err != nil{
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status": false,
					"message": "Failed to PUT data",
					"errors": "No Data Category",
					"data": nil,
				})
			}
			if err := c.BodyParser(category); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": false,
					"message": "Failed to PUT data",
					"errors": "Something bad behind the scene happens",
					"data": nil,
				})
			}
			db.Model(&category).Update("nama_category", category.NamaCategory)
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"message": "Succeed to PUT data",
				"errors": nil,
				"data": "",
			})
		})

		app.Delete("/category/:id", isAdmin, func(c *fiber.Ctx) error {
			id, _ := strconv.Atoi(c.Params("id"))
			category := new(Category)
			if err := db.First(&category, id).Error; err != nil{
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": false,
					"message": "Failed to DELETE data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			if err := db.Delete(&category, id).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to DELETE data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"message": "Succeed to DELETE data",
				"errors": nil,
				"data": "",
			})
		})

		app.Post("/auth/register", func(c *fiber.Ctx) error {
			user := new(User)
			
			type TempUser struct {
				ID uint `gorm:"primaryKey" json:"id"`
				CreatedAt time.Time `json:"created_at"`
				UpdatedAt time.Time `json:"updated_at"`
				Nama string `gorm:"type:varchar(255)" json:"nama"`
				KataSandi string `gorm:"type:varchar(255)" json:"kata_sandi"`
				NoTelp string `gorm:"type:varchar(255)" json:"no_telp" gorm:"unique:not null"`
				TanggalLahir string `json:"tanggal_lahir"`
				JenisKelamin string `gorm:"type:varchar(255);" json:"jenis_kelamin"`
				Tentang string `gorm:"type:text;" json:"tentang"`
				Pekerjaan string `gorm:"type:varchar(255)" json:"pekerjaan"`
				Email string `gorm:"type:varchar(255)" json:"email"`
				IdProvinsi string `gorm:"type:varchar(255)" json:"id_provinsi"`
				IdKota string `gorm:"type:varchar(255)" json:"id_kota"`
				IsAdmin string `json:"isAdmin"`
				Alamats []Alamat `gorm:"foreignKey:UserID" json:"alamats"`
			}

			tempUser := new(TempUser)
			if err := c.BodyParser(tempUser); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": false,
					"message": "Failed to POST data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			err := db.Where("no_telp = ? AND email = ?", tempUser.NoTelp, tempUser.Email).First(&user).Error
			if err != nil && err.Error() != "record not found" {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to POST data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			if user.Nama != "" {
				if tempUser.NoTelp == user.NoTelp{
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"status": false,
						"message": "Failed to POST data",
						"errors": "Error 1062: Duplicate entry " + user.NoTelp + " for key 'user.no_telp'",
						"data": nil,
					})
				} else if tempUser.Email == user.Email{
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"status": false,
						"message": "Failed to POST data",
						"errors": "Error 1062: Duplicate entry " + user.Email + " for key 'users.email'",
						"data": nil,
					})
				}
			}
			user.Nama = tempUser.Nama
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempUser.KataSandi), bcrypt.DefaultCost)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": false,
					"message": "Failed to POST data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			user.KataSandi = string(hashedPassword)
			user.NoTelp = tempUser.NoTelp
			date, err := time.Parse("02/01/2006", tempUser.TanggalLahir)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to POST data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			user.TanggalLahir = date
			user.Pekerjaan = tempUser.Pekerjaan
			user.Email = tempUser.Email
			user.IdProvinsi = tempUser.IdProvinsi
			user.IdKota = tempUser.IdKota
			user.JenisKelamin = tempUser.JenisKelamin
			user.Tentang = tempUser.Tentang
			user.IsAdmin = true
			toko := new(Toko)
			names := strings.Split(user.Nama, " ")

			toko.NamaToko = ""
			for _, name := range names {
				toko.NamaToko  += name[0:3] + "-"
			}
			toko.NamaToko = toko.NamaToko[0:len(toko.NamaToko) - 1]
			toko.UrlFoto = ""
			if err := db.Create(&toko).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to POST data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			if err := db.Model(&toko).Update("user_id", toko.ID).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to POST data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			if err := db.Create(&user).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to POST data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"message": "Succeed to POST data",
				"errors": nil,
				"data": "Register Succeed",
			})
		})

		app.Post("/auth/login", func(c *fiber.Ctx) error {
			users := []User{}
			db.Find(&users);
			type TempUser struct {
				NoTelp string `gorm:"type:varchar(255)" json:"no_telp" gorm:"unique:not null"`
				KataSandi string `gorm:"type:varchar(255)" json:"kata_sandi"`
			}

			tempUser := new(TempUser)
			if err := c.BodyParser(tempUser); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": false,
					"message": "Failed to POST data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			for _, user := range users {
				err := bcrypt.CompareHashAndPassword([]byte(user.KataSandi), []byte(tempUser.KataSandi))
				if (err != nil) {
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
						"status": false,
						"message": "Failed to POST data",
						"errors": "No Telp atau kata sandi salah",
						"data": nil,
					})
				}
				if user.NoTelp == tempUser.NoTelp {
					resp, err := http.Get("https://www.emsifa.com/api-wilayah-indonesia/api/provinces.json")
					if err != nil {
						return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
								"status": false,
								"message": "Failed to GET Province and Kota Data",
								"errors": err.Error(),
								"data": nil,
						})
					}

					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
								"status": false,
								"message": "Failed to GET Province and Kota Data",
								"errors": err.Error(),
								"data": nil,
						})
					}
					type Province struct{
						ID string
						Name string
					}
					provinces := []Province{}
					data := string(body)
					err = json.Unmarshal([]byte(data), &provinces)
					if err != nil {
						return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
								"status": false,
								"message": "Failed to GET Province and Kota Data",
								"errors": err.Error(),
								"data": nil,
						})
					}
					tempProvince := Province{}
					for _, province := range provinces{
						if user.IdProvinsi == province.ID{
							tempProvince.ID = province.ID
							tempProvince.Name = province.Name
							break
						}
					}

					resp, err = http.Get("https://www.emsifa.com/api-wilayah-indonesia/api/regencies/" + tempProvince.ID + ".json")
					if err != nil {
						return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
								"status": false,
								"message": "Failed to GET Province and Kota Data",
								"errors": err.Error(),
								"data": nil,
						})
					}

					body, err = ioutil.ReadAll(resp.Body)
					if err != nil {
						return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
								"status": false,
								"message": "Failed to GET Province and Kota Data",
								"errors": err.Error(),
								"data": nil,
						})
					}
					type Kota struct{
						ID string
						ProvinceID string
						Name string
					}
					kotas := []Kota{}
					data = string(body)
					err = json.Unmarshal([]byte(data), &kotas)
					if err != nil {
						return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
								"status": false,
								"message": "Failed to GET Province and Kota Data",
								"errors": err.Error(),
								"data": nil,
						})
					}
					tempKota:= Kota{}
					for _, kota := range kotas{
						if user.IdKota == kota.ID{
							tempKota.ID = kota.ID
							tempKota.Name = kota.Name
							tempKota.ProvinceID = tempProvince.ID
							break
						}
					}
					claims := jwt.MapClaims{
						"id": user.ID,
						"isAdmin": user.IsAdmin,
						"exp": time.Now().Add(time.Hour * 24).Unix(),
					}

					token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("secretkey"))
					if (err != nil){
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
								"status": false,
								"message": "Failed to Generate Token",
								"errors": err.Error(),
								"data": nil,
							})
					}
					return c.Status(fiber.StatusOK).JSON(fiber.Map{
						"status": true,
						"message": "Succeed to Login",
						"errors": nil,
						"data": fiber.Map{
							"nama": user.Nama,
							"no_telp": user.NoTelp,
							"tanggal_lahir": user.TanggalLahir.Format("02/01/2006"),
							"tentang": user.Tentang,
							"pekerjaan": user.Pekerjaan,
							"email": user.Email,
							"id_provinsi": fiber.Map{
								"id": tempProvince.ID,
								"name": tempProvince.Name,
							},
							"id_kota": fiber.Map{
								"id": tempKota.ID,
								"province_id": tempKota.ProvinceID,
								"name": tempKota.Name,
							},
							"token": token,
						},
					})
				} 
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": false,
				"message": "Failed to login",
				"errors": "You are unauthorized",
				"data": nil,
			})
		})

		app.Get("/user/alamat", isLoggedIn, func (c *fiber.Ctx) error{
			var user User
			sess, err := store.Get(c)
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			id := sess.Get("id")
			err = db.Preload("Alamats").First(&user, id).Error
			if err != nil && err.Error() != "record not found"{
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			type TempAlamat struct {
				ID uint `gorm:"primaryKey" json:"id"`
				JudulAlamat string `gorm:"type:varchar(255)" json:"judul_alamat"`
				NamaPenerima string `gorm:"type:varchar(255)" json:"nama_penerima"`
				NoTelp string `gorm:"type:varchar(255)" json:"no_telp" gorm:"unique:not null"`
				DetailAlamat string `gorm:"type:varchar(255)" json:"detail_alamat"`
			}
			tempAlamats := []TempAlamat{}
			for _, tempAlamat := range user.Alamats{
				tempAlamats = append(tempAlamats, TempAlamat{
					ID: tempAlamat.ID,
					JudulAlamat: tempAlamat.DetailAlamat,
					NamaPenerima: tempAlamat.NamaPenerima,
					NoTelp: tempAlamat.NoTelp,
					DetailAlamat: tempAlamat.DetailAlamat,
				})
			} 
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"message": "Success to GET data",
				"errors": nil,
				"data": tempAlamats,
			})
		})

		app.Get("/user/alamat/:id", isLoggedIn, func (c *fiber.Ctx) error{
			sess, err := store.Get(c)
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			user_id := sess.Get("id")
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			alamat_id, _ := strconv.Atoi(c.Params("id"))
			alamat := new(Alamat)
			err = db.First(&alamat, alamat_id).Error
			if (err != nil && err.Error() != "record not found"){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": err.Error(),
					"data": nil,
				})
			} else if (err != nil && err.Error() == "record not found"){
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": "No Data Alamat",
					"data": nil,
				})
			}
			if alamat.UserID == user_id {
				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"status": true,
					"message": "Success to GET data",
					"errors": nil,
					"data": fiber.Map{
						"id": alamat.ID,
						"judul_alamat": alamat.JudulAlamat,
						"nama_penerima": alamat.NamaPenerima,
						"no_telp": alamat.NoTelp,
						"detail_alamat": alamat.DetailAlamat,
					},
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": "You are unauthorized",
					"data": nil,
				})
		})

		app.Post("/user/alamat", isLoggedIn, func (c *fiber.Ctx) error {
			sess, err := store.Get(c)
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to POST data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			user_id := sess.Get("id")
			alamat := new(Alamat)
			if err := c.BodyParser(alamat); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": false,
					"message": "Failed to POST data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			alamat.UserID = user_id.(uint)
			if err := db.Create(&alamat).Error; err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": false,
					"message": "Failed to POST data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"message": "Succeed to POST data",
				"errors": nil,
				"data": 1,
			})
		})

		app.Put("/user/alamat/:id", isLoggedIn, func(c *fiber.Ctx) error {
			alamat_id, _ := strconv.Atoi(c.Params("id"))
			alamat := new(Alamat)
			sess, err := store.Get(c)
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to PUT data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			user_id := sess.Get("id")
			err = db.First(&alamat, alamat_id).Error
			if (err != nil){
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": false,
					"message": "Failed to PUT data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			if alamat.UserID != user_id {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"status": false,
					"message": "Failed to PUT data",
					"errors": "You are unauthorized",
					"data": nil,
				})
			}
			if err := c.BodyParser(alamat); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": false,
					"message": "Failed to PUT data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			err = db.Model(&alamat).Update(alamat).Error
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to PUT data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"message": "Succeed to PUT data",
				"errors": nil,
				"data": "",
			})
		})

		app.Delete("/user/alamat/:id", isLoggedIn, func(c *fiber.Ctx) error {
			alamat_id, _ := strconv.Atoi(c.Params("id"))
			alamat := new(Alamat)
			sess, err := store.Get(c)
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to DELETE data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			user_id := sess.Get("id")
			err = db.First(&alamat, alamat_id).Error
			if (err != nil){
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": false,
					"message": "Failed to DELETE data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			if alamat.UserID != user_id {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"status": false,
					"message": "Failed to DELETE data",
					"errors": "You are unauthorized",
					"data": nil,
				})
			}
			err = db.Delete(&alamat, alamat_id).Error
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to DELETE data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"message": "Succeed to DELETE data",
				"errors": nil,
				"data": "",
			})
		})

		app.Get("/user", isLoggedIn, func (c *fiber.Ctx) error {
			user := new(User)
			sess, err := store.Get(c)
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			user_id := sess.Get("id")
			err = db.First(&user, user_id).Error
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"status": true,
					"message": "Succeed to GET data",
					"errors": nil,
					"data": fiber.Map{
						"nama": user.Nama,
						"kata_sandi": user.KataSandi,
						"no_telp": user.NoTelp,
						"tanggal_lahir": user.TanggalLahir.Format("02/01/2006"),
						"pekerjaan": user.Pekerjaan,
						"email": user.Email,
						"id_provinsi": user.IdProvinsi,
						"id_kota": user.IdKota,
					},
				})
			
		})

		app.Put("/user", isLoggedIn, func (c *fiber.Ctx) error {
			user := new(User)
			sess, err := store.Get(c)
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to PUT data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			user_id := sess.Get("id")
			err = db.First(&user, user_id).Error
			if (err != nil && err.Error() == "record not found"){
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": false,
					"message": "Failed to PUT data",
					"errors": err.Error(),
					"data": nil,
				})
			} else if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to PUT data",
					"errors": err.Error(),
					"data": nil,
				})
			}

			type TempUser struct {
				ID uint `gorm:"primaryKey" json:"id"`
				CreatedAt time.Time `json:"created_at"`
				UpdatedAt time.Time `json:"updated_at"`
				Nama string `gorm:"type:varchar(255)" json:"nama"`
				KataSandi string `gorm:"type:varchar(255)" json:"kata_sandi"`
				NoTelp string `gorm:"type:varchar(255)" json:"no_telp" gorm:"unique:not null"`
				TanggalLahir string `json:"tanggal_lahir"`
				JenisKelamin string `gorm:"type:varchar(255);" json:"jenis_kelamin"`
				Tentang string `gorm:"type:text;" json:"tentang"`
				Pekerjaan string `gorm:"type:varchar(255)" json:"pekerjaan"`
				Email string `gorm:"type:varchar(255)" json:"email"`
				IdProvinsi string `gorm:"type:varchar(255)" json:"id_provinsi"`
				IdKota string `gorm:"type:varchar(255)" json:"id_kota"`
				IsAdmin bool `json:"isAdmin"`
				Alamats []Alamat `gorm:"foreignKey:UserID" json:"alamats"`
			}

			tempUser := new(TempUser)
			tempUser.Nama = user.Nama
			tempUser.NoTelp = user.NoTelp
			tempUser.TanggalLahir = user.TanggalLahir.Format("02/01/2006")
			tempUser.Pekerjaan = user.Pekerjaan
			tempUser.Email = user.Email
			tempUser.IdProvinsi = user.IdProvinsi
			tempUser.IdKota = user.IdKota
			tempUser.JenisKelamin = user.JenisKelamin
			tempUser.IsAdmin = user.IsAdmin
			tempUser.Alamats = user.Alamats
			tempUser.Tentang = user.Tentang
			users := []User{}

			if err := c.BodyParser(tempUser); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": false,
					"message": "Failed to PUT data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			err = db.Find(&users).Error
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to PUT data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			for _, loopUser:= range users{
				if loopUser.ID != user_id{
					if loopUser.NoTelp == tempUser.NoTelp{
						return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
							"status": false,
							"message": "Failed to POST data",
							"errors": "Error 1062: Duplicate entry " + user.NoTelp + " for key 'user.no_telp'",
							"data": nil,
						})
					} else if loopUser.Email == tempUser.Email{
						return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
							"status": false,
							"message": "Failed to POST data",
							"errors": "Error 1062: Duplicate entry " + user.Email + " for key 'users.email'",
							"data": nil,
						})
					} 
				}
			}
			user.Nama = tempUser.Nama
			if tempUser.KataSandi == ""{
				tempUser.KataSandi = user.KataSandi
				user.KataSandi = tempUser.KataSandi
			} else {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempUser.KataSandi), bcrypt.DefaultCost)
				if err != nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"status": false,
						"message": "Failed to PUT data",
						"errors": err.Error(),
						"data": nil,
					})
				}
				user.KataSandi = string(hashedPassword)
			}
			user.NoTelp = tempUser.NoTelp
			date, err := time.Parse("02/01/2006", tempUser.TanggalLahir)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": false,
					"message": "Failed to PUT data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			user.TanggalLahir = date
			user.Pekerjaan = tempUser.Pekerjaan
			user.Email = tempUser.Email
			user.IdProvinsi = tempUser.IdProvinsi
			user.IdKota = tempUser.IdKota
			user.Alamats = tempUser.Alamats
			user.Tentang = tempUser.Tentang
			user.JenisKelamin = tempUser.JenisKelamin
			err = db.Model(&user).Update(&user).Error
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to PUT data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"message": "Succeed to PUT data",
				"errors": nil,
				"data": fiber.Map{
						"nama": user.Nama,
						"kata_sandi": user.KataSandi,
						"no_telp": user.NoTelp,
						"tanggal_lahir": user.TanggalLahir.Format("02/01/2006"),
						"pekerjaan": user.Pekerjaan,
						"email": user.Email,
						"id_provinsi": user.IdProvinsi,
						"id_kota": user.IdKota,
					},
			})
		})

		app.Get("/toko/my", isLoggedIn, func(c *fiber.Ctx) error {
			toko := new(Toko)
			sess, err := store.Get(c)
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			user_id := sess.Get("id")
			err = db.Where("user_id = ?", user_id).Find(&toko).Error
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"message": "Succeed to GET data",
				"errors": nil,
				"data": fiber.Map{
					"id": toko.ID,
					"nama_toko": toko.NamaToko,
					"url_foto": toko.UrlFoto,
					"user_id": toko.UserID,
				},
			})
		})

		app.Put("/toko/:id_toko", isLoggedIn, func(c *fiber.Ctx) error {
			toko := new(Toko)
			sess, err := store.Get(c)
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to PUT data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			user_id := sess.Get("id")
			toko_id, err := strconv.Atoi(c.Params("id_toko"))
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to PUT data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			err = db.First(&toko, toko_id).Error
			if (err != nil && err.Error() == "record not found"){
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status": false,
					"message": "Failed to PUT data",
					"errors": "Toko tidak ditemukan",
					"data": nil,
				})
			} else if (err != nil && err.Error() != "record not found"){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to PUT data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			if (user_id == toko.UserID){
				if form, err := c.MultipartForm(); err == nil {
					files := form.File["photo"]
					for _, file := range files {
						time := strconv.FormatInt(time.Now().Unix(), 10)
						if err := c.SaveFile(file, fmt.Sprintf("./images/%s-%s", time, file.Filename)); err != nil {
							return err
						}
						toko.UrlFoto = time + "-" + file.Filename
					}
					if nama_toko := form.Value["nama_toko"]; len(nama_toko) > 0{
						toko.NamaToko = nama_toko[0]
					}
				}
				err = db.Model(&toko).Update(&toko).Error
				if (err == nil){
					return c.Status(fiber.StatusOK).JSON(fiber.Map{
						"status": true,
						"message": "Succeed to PUT data",
						"errors": nil,
						"data": "Update toko succeed",
					})
				}
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to PUT data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": false,
				"message": "Failed to PUT data",
				"errors": "You are unauthorized",
				"data": nil,
			})
		})

		app.Get("/toko/:id_toko", isLoggedIn, func(c *fiber.Ctx) error {
			toko := new(Toko)
			toko_id, err := strconv.Atoi(c.Params("id_toko"))
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			err = db.First(&toko, toko_id).Error
			if (err != nil && err.Error() == "record not found"){
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": "Toko tidak ditemukan",
					"data": nil,
				})
			} else if (err != nil && err.Error() != "record not found"){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"message": "Succeed to GET data",
				"errors": nil,
				"data": fiber.Map{
					"id": toko.ID,
					"nama_toko": toko.NamaToko,
					"url_foto": toko.UrlFoto,
				},
			})
		})

		app.Get("/toko", isLoggedIn, func(c *fiber.Ctx) error {
			limit := c.QueryInt("limit", 10)
			page := c.QueryInt("page", 1)

			offset := (page - 1) * limit

			tokos := []Toko{}
			err := db.Limit(limit).Offset(offset).Find(&tokos).Error
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			type TempToko struct{
				ID uint `gorm:"primaryKey" json:"id"`
				NamaToko string `gorm:"type:varchar(255)" json:"nama_toko"`
				UrlFoto string `gorm:"type:varchar(255)" json:"url_foto"`
			}
			tempTokos := []TempToko{}
			for _, toko := range tokos{
				tempTokos = append(tempTokos, TempToko{
					ID: toko.ID,
					NamaToko: toko.NamaToko,
					UrlFoto: toko.UrlFoto,
				})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"message": "Succeed to GET data",
				"errors": nil,
				"data": fiber.Map{
					"page": page,
					"limit": limit,
					"data": tempTokos,
				},
			})
		})

		app.Get("/product", func(c *fiber.Ctx) error {
			limit := c.QueryInt("limit", 10)
			page := c.QueryInt("page", 1)

			offset := (page - 1) * limit

			produks := []Produk{}
			db.Limit(limit).Offset(offset).Preload("Toko").Preload("Category").Preload("Foto_Produks").Find(&produks)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			type TempToko struct{
				ID uint `gorm:"primaryKey" json:"id"`
				NamaToko string `gorm:"type:varchar(255)" json:"nama_toko"`
				UrlFoto string `gorm:"type:varchar(255)" json:"url_foto"`
			}

			type Temp_Foto_Produk struct{
				ID uint `gorm:"primaryKey" json:"id"`
				Url string `gorm:"type:varchar(255)" json:"url"`
				ProdukID uint `json:"produk_id"`
			}

			type TempProduk struct{
				ID uint `gorm:"primaryKey" json:"id"`
				NamaProduk string `gorm:"type:varchar(255)" json:"nama_produk"`
				Slug string `gorm:"type:varchar(255)" json:"slug"`
				HargaReseller int `json:"harga_reseller"`
				HargaKonsumen int `json:"harga_konsumen"`
				Stok int `json:"stok"`
				Deskripsi string `gorm:"type:text" json:"deskripsi"`
				Toko TempToko `json:"toko"`
				Category TempCategory `json:"category"`
				Photos []Temp_Foto_Produk `gorm:"foreignKey:ProdukID" json:"foto_produks"`
			}
			tempProduks := []TempProduk{}
			for _, produk := range produks{
				tempFotoProduks := []Temp_Foto_Produk{}
				for _, foto_produk := range produk.Foto_Produks{
					tempFotoProduks = append(tempFotoProduks, Temp_Foto_Produk{
						ID: foto_produk.ID,
						Url: foto_produk.Url,
						ProdukID: foto_produk.ProdukID,
					})
				}
				tempProduks = append(tempProduks, TempProduk{
					ID: produk.ID,
					NamaProduk: produk.NamaProduk,
					Slug: produk.Slug,
					HargaReseller: produk.HargaReseller,
					HargaKonsumen: produk.HargaKonsumen,
					Stok: produk.Stok,
					Deskripsi: produk.Deskripsi,
					Toko: TempToko{
						ID: produk.TokoID,
						NamaToko: produk.Toko.NamaToko,
						UrlFoto: produk.Toko.UrlFoto,
					},
					Photos: tempFotoProduks,
					Category: TempCategory{
						ID: produk.CategoryID,
						NamaCategory: produk.Category.NamaCategory,
					},
				})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"message": "Succeed to GET data",
				"errors": nil,
				"data": tempProduks,
			})
		})

		app.Get("/product/:id", func(c *fiber.Ctx) error {
			produk := Produk{}
			id, _ := strconv.Atoi(c.Params("id"))
			if err := db.Preload("Toko").Preload("Category").Preload("Foto_Produks").First(&produk, id).Error; err != nil {
				if err.Error() == "record not found"{
					return c.Status(fiber.StatusOK).JSON(fiber.Map{
						"status": false,
						"message": "Failed to GET data",
						"errors": "No Data Product",
						"data": nil,
					})
				}
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			type Temp_Foto_Produk struct{
				ID uint `gorm:"primaryKey" json:"id"`
				Url string `gorm:"type:varchar(255)" json:"url"`
				ProdukID uint `json:"produk_id"`
			}
			tempFotoProduks := []Temp_Foto_Produk{}
			for _, foto_produk := range produk.Foto_Produks{
				tempFotoProduks = append(tempFotoProduks, Temp_Foto_Produk{
					ID: foto_produk.ID,
					Url: foto_produk.Url,
					ProdukID: foto_produk.ProdukID,
				})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"message": "Succeed to GET data",
				"errors": nil,
				"data": fiber.Map{
					"id": produk.ID,
					"nama_produk": produk.NamaProduk,
					"slug": produk.Slug,
					"harga_reseller": produk.HargaReseller,
					"harga_konsumen": produk.HargaKonsumen,
					"stok": produk.Stok,
					"deskripsi": produk.Deskripsi,
					"toko": fiber.Map{
						"id": produk.TokoID,
						"nama_toko": produk.Toko.NamaToko,
						"url_foto": produk.Toko.UrlFoto,
					},
					"category": fiber.Map{
						"id": produk.CategoryID,
						"nama_category": produk.Category.NamaCategory,
					},
					"photos": tempFotoProduks,
				},
			})
		})

		app.Post("/product", isLoggedIn, func(c *fiber.Ctx) error {
			sess, err := store.Get(c)
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to POST data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			user_id := sess.Get("id")
			toko := Toko{}
			err = db.Preload("User").Where("user_id = ?", user_id).Find(&toko).Error
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to POST data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			if (user_id == toko.UserID){
				produk := Produk{}
				produk.TokoID = toko.ID
				if form, err := c.MultipartForm(); err == nil {
					if nama_produk := form.Value["nama_produk"]; len(nama_produk) > 0{
						produk.NamaProduk = nama_produk[0]
					}
					if category_id := form.Value["category_id"]; len(category_id) > 0{
						temp, _ := strconv.ParseUint(category_id[0], 10, 64)
						produk.CategoryID = uint(temp)
					}
					if harga_reseller := form.Value["harga_reseller"]; len(harga_reseller) > 0{
						temp, _ := strconv.Atoi(harga_reseller[0])
						produk.HargaReseller = temp
					}
					if harga_konsumen := form.Value["harga_konsumen"]; len(harga_konsumen) > 0{
						temp, _ := strconv.Atoi(harga_konsumen[0])
						produk.HargaKonsumen = temp
					}
					if stok := form.Value["stok"]; len(stok) > 0{
						temp, _ := strconv.Atoi(stok[0])
						produk.Stok = temp
					}
					if deskripsi := form.Value["deskripsi"]; len(deskripsi) > 0{
						produk.Deskripsi = deskripsi[0]
					}
					if slug := form.Value["slug"]; len(slug) > 0{
						produk.Slug = slug[0]
					}
					category := Category{}
					err = db.First(&category, produk.CategoryID).Error
					if err != nil && err.Error() == "record not found"{
						return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
							"status": false,
							"message": "Failed to POST data",
							"errors": "No Data Category",
							"data": nil,
						})
					} else if err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
							"status": false,
							"message": "Failed to POST data",
							"errors": err.Error(),
							"data": nil,
						})
					}
					err = db.Create(&produk).Error
					if err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
							"status": false,
							"message": "Failed to POST data",
							"errors": err.Error(),
							"data": nil,
						})
					}
					files := form.File["photos"]
					if len(files) != 0{
						for _, file := range files {
							fotoProduk := Foto_Produk{}
							time := strconv.FormatInt(time.Now().Unix(), 10)
							if err := c.SaveFile(file, fmt.Sprintf("./images/%s-%s", time, file.Filename)); err != nil {
								return err
							}
							fotoProduk.ProdukID = produk.ID
							fotoProduk.Url = "/images/" + time + "-" + file.Filename
							if err := db.Create(&fotoProduk).Error; err != nil {
								return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
									"status": false,
									"message": "Failed to POST data",
									"errors": err.Error(),
									"data": nil,
								})
							}
						}
					}
					return c.Status(fiber.StatusOK).JSON(fiber.Map{
						"status": true,
						"message": "Succeed to POST data",
						"errors": nil,
						"data": 4,
					})
				}
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to POST data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": false,
				"message": "Failed to POST data",
				"errors": "You are unauthorized",
				"data": nil,
			})
		})

		app.Put("/product/:id", isLoggedIn, func(c *fiber.Ctx) error {
			sess, err := store.Get(c)
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to PUT data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			user_id := sess.Get("id")
			toko := Toko{}
			err = db.Preload("User").Where("user_id = ?", user_id).Find(&toko).Error
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to PUT data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			if (user_id == toko.UserID){
				produk := Produk{}
				produk_id, _ := strconv.Atoi(c.Params("id"))
				err = db.First(&produk, produk_id).Error
				if err != nil && err.Error() == "record not found" {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"status": false,
						"message": "Failed to PUT data",
						"errors": "No Data Produk",
						"data": nil,
					})
				} else if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"status": false,
						"message": "Failed to PUT data",
						"errors": err.Error(),
						"data": nil,
					})
				}

				if produk.TokoID == toko.ID{
					if form, err := c.MultipartForm(); err == nil {
						if nama_produk := form.Value["nama_produk"]; len(nama_produk) > 0{
							produk.NamaProduk = nama_produk[0]
						}
						if category_id := form.Value["category_id"]; len(category_id) > 0{
							temp, _ := strconv.ParseUint(category_id[0], 10, 64)
							produk.CategoryID = uint(temp)
						}
						if harga_reseller := form.Value["harga_reseller"]; len(harga_reseller) > 0{
							temp, _ := strconv.Atoi(harga_reseller[0])
							produk.HargaReseller = temp
						}
						if harga_konsumen := form.Value["harga_konsumen"]; len(harga_konsumen) > 0{
							temp, _ := strconv.Atoi(harga_konsumen[0])
							produk.HargaKonsumen = temp
						}
						if stok := form.Value["stok"]; len(stok) > 0{
							temp, _ := strconv.Atoi(stok[0])
							produk.Stok = temp
						}
						if deskripsi := form.Value["deskripsi"]; len(deskripsi) > 0{
							produk.Deskripsi = deskripsi[0]
						}
						if slug := form.Value["slug"]; len(slug) > 0{
							produk.Slug = slug[0]
						}
						category := Category{}
						err = db.First(&category, produk.CategoryID).Error
						if err != nil && err.Error() == "record not found"{
							return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
								"status": false,
								"message": "Failed to PUT data",
								"errors": "No Data Category",
								"data": nil,
							})
						} else if err != nil {
							return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
								"status": false,
								"message": "Failed to PUT data",
								"errors": err.Error(),
								"data": nil,
							})
						}
						err = db.Model(&produk).Update(produk).Error
						if err != nil {
							return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
								"status": false,
								"message": "Failed to PUT data",
								"errors": err.Error(),
								"data": nil,
							})
						}
						files := form.File["photos"]
						if len(files) != 0{ 
							for _, file := range files {
								fotoProduk := Foto_Produk{}
								time := strconv.FormatInt(time.Now().Unix(), 10)
								if err := c.SaveFile(file, fmt.Sprintf("./images/%s-%s", time, file.Filename)); err != nil {
									return err
								}
								fotoProduk.ProdukID = produk.ID
								fotoProduk.Url = "/images/" + time + "-" + file.Filename
								if err := db.Create(&fotoProduk).Error; err != nil {
									return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
										"status": false,
										"message": "Failed to PUT data",
										"errors": err.Error(),
										"data": nil,
									})
								}
							}
						}
						return c.Status(fiber.StatusOK).JSON(fiber.Map{
							"status": true,
							"message": "Succeed to PUT data",
							"errors": nil,
							"data": "",
						})
					}
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to PUT data",
					"errors": err.Error(),
					"data": nil,
				})
				}
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"status": false,
					"message": "Failed to PUT data",
					"errors": "You are unauthorized",
					"data": nil,
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": false,
				"message": "Failed to PUT data",
				"errors": "You are unauthorized",
				"data": nil,
			})
		})

		app.Delete("/product/:id", isLoggedIn, func(c *fiber.Ctx) error {
			sess, err := store.Get(c)
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to DELETE data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			user_id := sess.Get("id")
			toko := Toko{}
			err = db.Preload("User").Where("user_id = ?", user_id).Find(&toko).Error
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to DELETE data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			if (user_id == toko.UserID){
				produk := Produk{}
				produk_id, _ := strconv.Atoi(c.Params("id"))
				err = db.First(&produk, produk_id).Error
				if err != nil && err.Error() == "record not found" {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"status": false,
						"message": "Failed to DELETE data",
						"errors": "No Data Produk",
						"data": nil,
					})
				} else if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"status": false,
						"message": "Failed to DELETE data",
						"errors": err.Error(),
						"data": nil,
					})
				}

				if produk.TokoID == toko.ID{
					err = db.Delete(&produk, produk.ID).Error
					if err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
							"status": false,
							"message": "Failed to DELETE data",
							"errors": err.Error(),
							"data": nil,
						})
					}
					return c.Status(fiber.StatusOK).JSON(fiber.Map{
						"status": true,
						"message": "Succeed to DELETE data",
						"errors": nil,
						"data": "",
					})
				}
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"status": false,
					"message": "Failed to DELETE data",
					"errors": "You are unauthorized",
					"data": nil,
				})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": false,
				"message": "Failed to DELETE data",
				"errors": "You are unauthorized",
				"data": nil,
			})
		})

		app.Get("/trx", isLoggedIn, func(c *fiber.Ctx) error {
			sess, err := store.Get(c)
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			user_id := sess.Get("id")
			trxes := []Trx{}
			err = db.Preload("Alamat").Preload("Detail_Trx").Where("user_id = ?", user_id).Find(&trxes).Error
			if err != nil && err.Error() != "record not found"{
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			type TempPhotos struct{
					ID uint `gorm:"primaryKey" json:"id"`
					Url string `gorm:"type:varchar(255)" json:"url"`
					ProdukID uint `json:"product_id"`
			}

			type TempToko struct{
					NamaToko string `gorm:"type:varchar(255)" json:"nama_toko"`
					UrlFoto string `gorm:"type:varchar(255)" json:"url_foto"`
			}

			type TempProduk struct{
					ID uint `gorm:"primaryKey" json:"id"`
					NamaProduk string `gorm:"type:varchar(255)" json:"nama_produk"`
					Slug string `gorm:"type:varchar(255)" json:"slug"`
					HargaReseller int `json:"harga_reseller"`
					HargaKonsumen int `json:"harga_konsumen"`
					Deskripsi string `gorm:"type:text" json:"deskripsi"`
					Toko TempToko `json:"toko"`
					Category TempCategory `json:"category"`
					Photos []TempPhotos `gorm:"foreignKey:ProdukID" json:"photos"`
			}
			type TempAlamat struct{
					ID uint `gorm:"primaryKey" json:"id"`
					JudulAlamat string `gorm:"type:varchar(255)" json:"judul_alamat"`
					NamaPenerima string `gorm:"type:varchar(255)" json:"nama_penerima"`
					NoTelp string `gorm:"type:varchar(255)" json:"no_telp" gorm:"unique:not null"`
					DetailAlamat string `gorm:"type:varchar(255)" json:"detail_alamat"`
			}

			type TempDetailTrx struct {
				Product TempProduk `json:"product"`
			}

			type TempTrx struct{
					ID uint `gorm:"primaryKey" json:"id"`
					HargaTotal int `json:"harga_total"`
					KodeInvoice string `gorm:"type:varchar(255)" json:"kode_invoice"`
					MethodBayar string `gorm:"type:varchar(255)" json:"method_bayar"`
					Alamat TempAlamat `json:"alamat"`
					Detail_Trx []TempDetailTrx `json:"detail_trx"`
			}
			tempTrx := []TempTrx{}
			for _, trx := range trxes{
				tempDetailTrx := []TempDetailTrx{}
				for _, detail := range trx.Detail_Trx{
					detailTrx := Detail_Trx{}
					err = db.Preload("Log_Produk").Preload("Toko").First(&detailTrx, detail.ID).Error
					if err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
							"status": false,
							"message": "Failed to GET data",
							"errors": err.Error(),
							"data": nil,
						})
					}
					logProduk := Log_Produk{}
					err = db.Preload("Category").Preload("Foto_Produks").First(&logProduk, detailTrx.LogProdukID).Error
					if err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
							"status": false,
							"message": "Failed to GET data",
							"errors": err.Error(),
							"data": nil,
						})
					}
					tempPhotos := []TempPhotos{}
					for _, fotoProduk := range logProduk.Foto_Produks{
						tempPhotos = append(tempPhotos, TempPhotos{
							ID: fotoProduk.ID,
							Url: fotoProduk.Url,
							ProdukID: fotoProduk.ProdukID,
						})
					}
					tempDetailTrx = append(tempDetailTrx, TempDetailTrx{
						Product: TempProduk{
							ID: logProduk.ProdukID,
							NamaProduk: logProduk.NamaProduk,
							Slug: logProduk.Slug,
							HargaReseller: logProduk.HargaReseller,
							HargaKonsumen: logProduk.HargaKonsumen,
							Deskripsi: logProduk.Deskripsi,
							Toko: TempToko{
								NamaToko: detailTrx.Toko.NamaToko,
								UrlFoto: detailTrx.Toko.UrlFoto,
							},
							Category: TempCategory{
								ID: logProduk.Category.ID,
								NamaCategory: logProduk.Category.NamaCategory,
							},
							Photos: tempPhotos,
						},
					})
				}
				tempTrx = append(tempTrx, TempTrx{
					ID: trx.ID,
					HargaTotal: trx.HargaTotal,
					KodeInvoice: trx.KodeInvoice,
					MethodBayar: trx.MethodBayar,
					Alamat: TempAlamat{
						ID: trx.Alamat.ID,
						NamaPenerima: trx.Alamat.NamaPenerima,
						JudulAlamat: trx.Alamat.JudulAlamat,
						NoTelp: trx.Alamat.NoTelp,
						DetailAlamat: trx.Alamat.DetailAlamat,
					},
					Detail_Trx: tempDetailTrx,
				})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"message": "Succeed to GET data",
				"errors": nil,
				"data": fiber.Map{
					"data": tempTrx, 
				},
			})
		})

		app.Get("/trx/:id", isLoggedIn, func(c *fiber.Ctx) error {
			sess, err := store.Get(c)
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			user_id := sess.Get("id")
			trx_id, _ := strconv.Atoi(c.Params("id"))
			trx := Trx{}
			err = db.Preload("Alamat").Preload("Detail_Trx").Where("id = ?", trx_id).Find(&trx).Error
			if err != nil && err.Error() != "record not found"{
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": err.Error(),
					"data": nil,
				})
			} else if err != nil && err.Error() == "record not found" {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": "No Data Trx",
					"data": nil,
				})
			}
			if trx.UserID != user_id {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": "You are unauthorized",
					"data": nil,
				})
			}
			type TempPhotos struct{
					ID uint `gorm:"primaryKey" json:"id"`
					Url string `gorm:"type:varchar(255)" json:"url"`
					ProdukID uint `json:"product_id"`
			}

			type TempProduk struct{
					ID uint `gorm:"primaryKey" json:"id"`
					NamaProduk string `gorm:"type:varchar(255)" json:"nama_produk"`
					Slug string `gorm:"type:varchar(255)" json:"slug"`
					HargaReseller int `json:"harga_reseller"`
					HargaKonsumen int `json:"harga_konsumen"`
					Deskripsi string `gorm:"type:text" json:"deskripsi"`
					Photos []TempPhotos `gorm:"foreignKey:ProdukID" json:"photos"`
			}

			type TempDetailTrx struct {
				Product TempProduk `json:"product"`
			}

			tempDetailTrx := []TempDetailTrx{}
			for _, detailTrx := range trx.Detail_Trx{
				logProduk := Log_Produk{}
				err = db.Preload("Category").Preload("Foto_Produks").First(&logProduk, detailTrx.LogProdukID).Error
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"status": false,
						"message": "Failed to GET data",
						"errors": err.Error(),
						"data": nil,
					})
				}
				tempPhotos := []TempPhotos{}
				for _, fotoProduk := range logProduk.Foto_Produks{
					tempPhotos = append(tempPhotos, TempPhotos{
						ID: fotoProduk.ID,
						Url: fotoProduk.Url,
						ProdukID: fotoProduk.ProdukID,
					})
				}
				tempDetailTrx = append(tempDetailTrx, TempDetailTrx{
					Product: TempProduk{
						ID: logProduk.ProdukID,
						NamaProduk: logProduk.NamaProduk,
						Slug: logProduk.Slug,
						HargaReseller: logProduk.HargaReseller,
						HargaKonsumen: logProduk.HargaKonsumen,
						Deskripsi: logProduk.Deskripsi,
						Photos: tempPhotos,
					},
				})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"message": "Succeed to GET data",
				"errors": nil,
				"data": fiber.Map{
					"id": trx_id,
					"harga_total": trx.HargaTotal,
					"kode_invoice": trx.KodeInvoice,
					"method_bayar": trx.MethodBayar,
					"alamat_kirim": fiber.Map{
						"id": trx.AlamatID,
						"judul_alamat": trx.Alamat.JudulAlamat,
						"nama_penerima": trx.Alamat.NamaPenerima,
						"no_telp": trx.Alamat.NoTelp,
						"detail_alamat": trx.Alamat.DetailAlamat,
					},
					"detail_trx": tempDetailTrx,
				},
			})
		})

		app.Post("/trx", isLoggedIn, func(c *fiber.Ctx) error {
			sess, err := store.Get(c)
			if (err != nil){
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to POST data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			user_id := sess.Get("id")
			type Temp_Detail_Trx struct {
				Kuantitas int `json:"kuantitas"`
				ProductID int `json:"product_id"`
			}

			type TempTrx struct {
				AlamatKirim int `json:"alamat_kirim"`
				MethodBayar string `gorm:"type:varchar(255)" json:"method_bayar"`
				Detail_Trx []Temp_Detail_Trx `json:"detail_trx"`
			}
			tempTrx := new(TempTrx)
			if err := c.BodyParser(tempTrx); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": false,
					"message": "Failed to POST data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			alamat := Alamat{}
			err = db.Where("user_id = ?", user_id).First(&alamat, tempTrx.AlamatKirim).Error
			if err != nil && err.Error() == "record not found"{
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": false,
					"message": "Failed to POST data",
					"errors": "No Data Alamat",
					"data": nil,
				})
			} else if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to POST data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			trx := Trx{}
			trx.AlamatID = uint(tempTrx.AlamatKirim)
			trx.MethodBayar = tempTrx.MethodBayar
			trx.UserID = user_id.(uint)
			time := strconv.FormatInt(time.Now().Unix(), 10)
			trx.KodeInvoice = "INV-" + time
			trx.HargaTotal = 0
			err = db.Create(&trx).Error
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to POST data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			hargaTotal := 0
			for _, detailTrx := range tempTrx.Detail_Trx {
				produk := Produk{}
				err = db.First(&produk, detailTrx.ProductID).Error
				if err != nil && err.Error() == "record not found"{
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"status": false,
						"message": "Failed to POST data",
						"errors": "No Data Produk",
						"data": nil,
					})
				} else if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"status": false,
						"message": "Failed to POST data",
						"errors": err.Error(),
						"data": nil,
					})
				}
				logProduk := Log_Produk{}
				err = db.Where("produk_id = ?", produk.ID).Find(&logProduk).Error
				if err != nil && err.Error() == "record not found"{
					logProduk.ProdukID = produk.ID
					logProduk.NamaProduk = produk.NamaProduk
					logProduk.Slug = produk.Slug
					logProduk.HargaReseller = produk.HargaReseller
					logProduk.HargaKonsumen = produk.HargaKonsumen
					logProduk.Deskripsi = produk.Deskripsi
					logProduk.TokoID = produk.TokoID
					logProduk.CategoryID = produk.CategoryID
					err = db.Create(&logProduk).Error
					if err != nil {
						return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
							"status": false,
							"message": "Failed to POST data",
							"errors": err.Error(),
							"data": nil,
						})
					}
				} else if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"status": false,
						"message": "Failed to POST data",
						"errors": err.Error(),
						"data": nil,
					})
				}
				hargaTotal += detailTrx.Kuantitas * produk.HargaKonsumen
				detail := Detail_Trx{}
				detail.TrxID = trx.ID
				detail.LogProdukID = logProduk.ID
				detail.TokoID = produk.TokoID
				detail.Kuantitas = detailTrx.Kuantitas
				detail.HargaTotal = produk.HargaKonsumen
				err = db.Create(&detail).Error
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"status": false,
						"message": "Failed to POST data",
						"errors": err.Error(),
						"data": nil,
					})
				}
				err = db.Model(&produk).Update("stok", produk.Stok - detail.Kuantitas).Error
				if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"status": false,
						"message": "Failed to POST data",
						"errors": err.Error(),
						"data": nil,
					})
				}
			}
			err = db.Model(&trx).Update("harga_total", hargaTotal).Error
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status": false,
					"message": "Failed to POST data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"message": "Succeed to POST data",
				"errors": nil,
				"data": 6,
			})
		})

		app.Get("/provcity/listprovincies", func(c *fiber.Ctx) error {
			resp, err := http.Get("https://www.emsifa.com/api-wilayah-indonesia/api/provinces.json")
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"status": false,
						"message": "Failed to GET data",
						"errors": err.Error(),
						"data": nil,
				})
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"status": false,
						"message": "Failed to GET data",
						"errors": err.Error(),
						"data": nil,
				})
			}
			type Province struct{
				ID string `json:"id"`
				Name string `json:"name"`
			}
			provinces := []Province{}
			data := string(body)
			err = json.Unmarshal([]byte(data), &provinces)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": err.Error(),
					"data": nil,
				})
			}

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"message": "Succeed to GET data",
				"errors": nil,
				"data": provinces,
			})
		})

		app.Get("/provcity/listcities/:prov_id", func(c *fiber.Ctx) error {
			resp, err := http.Get("https://www.emsifa.com/api-wilayah-indonesia/api/regencies/" + c.Params("prov_id") + ".json")
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"status": false,
						"message": "Failed to GET Province and Kota Data",
						"errors": err.Error(),
						"data": nil,
				})
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"status": false,
						"message": "Failed to GET Province and Kota Data",
						"errors": err.Error(),
						"data": nil,
				})
			}
			type Kota struct{
				ID string `json:"id"`
				ProvinceID string `json:"province_id"`
				Name string `json:"name"`
			}
			kotas := []Kota{}
			data := string(body)
			err = json.Unmarshal([]byte(data), &kotas)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET Province and Kota Data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status": true,
				"message": "Succeed to GET data",
				"errors": nil,
				"data": kotas,
			})
		})

		app.Get("/provcity/detailprovince/:prov_id", func(c *fiber.Ctx) error {
			prov_id := c.Params("prov_id")
			resp, err := http.Get("https://www.emsifa.com/api-wilayah-indonesia/api/provinces.json")
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"status": false,
						"message": "Failed to GET data",
						"errors": err.Error(),
						"data": nil,
				})
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"status": false,
						"message": "Failed to GET data",
						"errors": err.Error(),
						"data": nil,
				})
			}
			type Province struct{
				ID string `json:"id"`
				Name string `json:"name"`
			}
			provinces := []Province{}
			data := string(body)
			err = json.Unmarshal([]byte(data), &provinces)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			tempProvince := Province{}
			for _, province := range provinces {
				if province.ID == prov_id{
					tempProvince = province
					return c.Status(fiber.StatusOK).JSON(fiber.Map{
						"status": true,
						"message": "Succeed to GET data",
						"errors": nil,
						"data": tempProvince,
					})
				}
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": false,
				"message": "Failed to GET data",
				"errors": "No Data Province",
				"data": nil,
			})
		})

		app.Get("/provcity/detailcity/:city_id", func(c *fiber.Ctx) error {
			city_id := c.Params("city_id")
			resp, err := http.Get("https://www.emsifa.com/api-wilayah-indonesia/api/provinces.json")
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"status": false,
						"message": "Failed to GET data",
						"errors": err.Error(),
						"data": nil,
				})
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"status": false,
						"message": "Failed to GET data",
						"errors": err.Error(),
						"data": nil,
				})
			}
			type Province struct{
				ID string `json:"id"`
				Name string `json:"name"`
			}
			provinces := []Province{}
			data := string(body)
			err = json.Unmarshal([]byte(data), &provinces)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status": false,
					"message": "Failed to GET data",
					"errors": err.Error(),
					"data": nil,
				})
			}
			for _, province := range provinces {
				resp, err := http.Get("https://www.emsifa.com/api-wilayah-indonesia/api/regencies/" + province.ID + ".json")
				if err != nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
							"status": false,
							"message": "Failed to GET Province and Kota Data",
							"errors": err.Error(),
							"data": nil,
					})
				}

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
							"status": false,
							"message": "Failed to GET Province and Kota Data",
							"errors": err.Error(),
							"data": nil,
					})
				}
				type Kota struct{
					ID string `json:"id"`
					ProvinceID string `json:"province_id"`
					Name string `json:"name"`
				}
				kotas := []Kota{}
				data := string(body)
				err = json.Unmarshal([]byte(data), &kotas)
				if err != nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"status": false,
						"message": "Failed to GET Province and Kota Data",
						"errors": err.Error(),
						"data": nil,
					})
				}
				tempKota := Kota{}
				for _, kota := range kotas{
					if kota.ID == city_id {
						tempKota = kota
						return c.Status(fiber.StatusOK).JSON(fiber.Map{
							"status": true,
							"message": "Succeed to GET data",
							"errors": nil,
							"data": tempKota,
						})
					}
				}
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": false,
				"message": "Failed to GET data",
				"errors": "No Data City",
				"data": nil,
			})
		})
    app.Listen(":3000")
}