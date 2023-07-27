package controller

import (
	"project/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetMyAlamat(c *fiber.Ctx) error {
	var user model.User
	sess, err := Store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  err.Error(),
			"data":    nil,
		})
	}
	id := sess.Get("id")
	err = DB.Preload("Alamats").First(&user, id).Error
	if err != nil && err.Error() != "record not found" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  err.Error(),
			"data":    nil,
		})
	}
	tempAlamats := []model.TempAlamat{}
	for _, tempAlamat := range user.Alamats {
		tempAlamats = append(tempAlamats, model.TempAlamat{
			ID:           tempAlamat.ID,
			JudulAlamat:  tempAlamat.DetailAlamat,
			NamaPenerima: tempAlamat.NamaPenerima,
			NoTelp:       tempAlamat.NoTelp,
			DetailAlamat: tempAlamat.DetailAlamat,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Success to GET data",
		"errors":  nil,
		"data":    tempAlamats,
	})
}

func GetAlamatById(c *fiber.Ctx) error{
	sess, err := Store.Get(c)
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
	alamat := new(model.Alamat)
	err = DB.First(&alamat, alamat_id).Error
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
}

func CreateAlamat(c *fiber.Ctx) error {
	sess, err := Store.Get(c)
	if (err != nil){
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to POST data",
			"errors": err.Error(),
			"data": nil,
		})
	}
	user_id := sess.Get("id")
	alamat := new(model.Alamat)
	if err := c.BodyParser(alamat); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to POST data",
			"errors": err.Error(),
			"data": nil,
		})
	}
	alamat.UserID = user_id.(uint)
	if err := DB.Create(&alamat).Error; err != nil {
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
}

func UpdateAlamat(c *fiber.Ctx) error {
	alamat_id, _ := strconv.Atoi(c.Params("id"))
	alamat := new(model.Alamat)
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
	err = DB.First(&alamat, alamat_id).Error
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
	err = DB.Model(&alamat).Update(&alamat).Error
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
}

func DeleteAlamat(c *fiber.Ctx) error {
	alamat_id, _ := strconv.Atoi(c.Params("id"))
	alamat := new(model.Alamat)
	sess, err := Store.Get(c)
	if (err != nil){
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to DELETE data",
			"errors": err.Error(),
			"data": nil,
		})
	}
	user_id := sess.Get("id")
	err = DB.First(&alamat, alamat_id).Error
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
	err = DB.Delete(&alamat, alamat_id).Error
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
}
		