package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"project/model"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *fiber.Ctx) error {
	user := new(model.User)

	type TempUser struct {
		ID           uint      `gorm:"primaryKey" json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Nama         string    `gorm:"type:varchar(255)" json:"nama"`
		KataSandi    string    `gorm:"type:varchar(255)" json:"kata_sandi"`
		NoTelp       string    `gorm:"type:varchar(255)" json:"no_telp" gorm:"unique:not null"`
		TanggalLahir string    `json:"tanggal_lahir"`
		JenisKelamin string    `gorm:"type:varchar(255);" json:"jenis_kelamin"`
		Tentang      string    `gorm:"type:text;" json:"tentang"`
		Pekerjaan    string    `gorm:"type:varchar(255)" json:"pekerjaan"`
		Email        string    `gorm:"type:varchar(255)" json:"email"`
		IdProvinsi   string    `gorm:"type:varchar(255)" json:"id_provinsi"`
		IdKota       string    `gorm:"type:varchar(255)" json:"id_kota"`
		IsAdmin      string    `json:"isAdmin"`
		Alamats      []model.Alamat  `gorm:"foreignKey:UserID" json:"alamats"`
	}

	tempUser := new(TempUser)
	if err := c.BodyParser(tempUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  err.Error(),
			"data":    nil,
		})
	}
	err := DB.Where("no_telp = ? AND email = ?", tempUser.NoTelp, tempUser.Email).First(&user).Error
	if err != nil && err.Error() != "record not found" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  err.Error(),
			"data":    nil,
		})
	}
	if user.Nama != "" {
		if tempUser.NoTelp == user.NoTelp {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to POST data",
				"errors":  "Error 1062: Duplicate entry " + user.NoTelp + " for key 'user.no_telp'",
				"data":    nil,
			})
		} else if tempUser.Email == user.Email {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to POST data",
				"errors":  "Error 1062: Duplicate entry " + user.Email + " for key 'users.email'",
				"data":    nil,
			})
		}
	}
	user.Nama = tempUser.Nama
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempUser.KataSandi), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  err.Error(),
			"data":    nil,
		})
	}
	user.KataSandi = string(hashedPassword)
	user.NoTelp = tempUser.NoTelp
	date, err := time.Parse("02/01/2006", tempUser.TanggalLahir)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  err.Error(),
			"data":    nil,
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
	toko := new(model.Toko)
	names := strings.Split(user.Nama, " ")

	toko.NamaToko = ""
	for _, name := range names {
		toko.NamaToko += name[0:3] + "-"
	}
	toko.NamaToko = toko.NamaToko[0 : len(toko.NamaToko)-1]
	toko.UrlFoto = ""
	if err := DB.Create(&toko).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  err.Error(),
			"data":    nil,
		})
	}
	if err := DB.Model(&toko).Update("user_id", toko.ID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  err.Error(),
			"data":    nil,
		})
	}
	if err := DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  err.Error(),
			"data":    nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to POST data",
		"errors":  nil,
		"data":    "Register Succeed",
	})
}

func Login(c *fiber.Ctx) error {
	users := []model.User{}
	err := DB.Find(&users).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  err.Error(),
			"data":    nil,
		})
	}
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
			provinces := []model.Province{}
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
			tempProvince := model.Province{}
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
			kotas := []model.Kota{}
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
			tempKota:= model.Kota{}
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
}