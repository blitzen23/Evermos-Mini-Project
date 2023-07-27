package controller

import (
	"project/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func GetUser(c *fiber.Ctx) error {
	user := new(model.User)
	sess, err := Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  err.Error(),
			"data":    nil,
		})
	}
	user_id := sess.Get("id")
	err = DB.First(&user, user_id).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  err.Error(),
			"data":    nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to GET data",
		"errors":  nil,
		"data": fiber.Map{
			"nama":          user.Nama,
			"kata_sandi":    user.KataSandi,
			"no_telp":       user.NoTelp,
			"tanggal_lahir": user.TanggalLahir.Format("02/01/2006"),
			"pekerjaan":     user.Pekerjaan,
			"email":         user.Email,
			"id_provinsi":   user.IdProvinsi,
			"id_kota":       user.IdKota,
		},
	})
}

func UpdateUser(c *fiber.Ctx) error {
	user := new(model.User)
	sess, err := Store.Get(c)
	if (err != nil){
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to PUT data",
			"errors": err.Error(),
			"data": nil,
		})
	}
	user_id := sess.Get("id")
	err = DB.First(&user, user_id).Error
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
		Alamats []model.Alamat `gorm:"foreignKey:UserID" json:"alamats"`
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
	users := []model.User{}

	if err := c.BodyParser(tempUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to PUT data",
			"errors": err.Error(),
			"data": nil,
		})
	}
	err = DB.Find(&users).Error
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
	err = DB.Model(&user).Update(&user).Error
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
}