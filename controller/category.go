package controller

import (
	"project/model"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetAllCategory(c *fiber.Ctx) error {
	categories := []model.Category{}
	err := DB.Select("id, nama_category").Find(&categories).Error
	if err != nil && err.Error() == "record not found" {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  true,
			"message": "Succeed to GET data",
			"errors":  nil,
			"data":    categories,
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  err.Error(),
			"data":    nil,
		})
	} else {
		tempCategories := []model.TempCategory{}
		for _, category := range categories {
			tempCategories = append(tempCategories, model.TempCategory{
				ID:           category.ID,
				NamaCategory: category.NamaCategory,
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  true,
			"message": "Succeed to GET data",
			"errors":  nil,
			"data":    tempCategories,
		})
	}
}

func GetCategoryById(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	category := new(model.Category)
	err := DB.Select("id, nama_category").First(&category, id).Error
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
			"errors": err.Error(),
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
}

func CreateCategory(c *fiber.Ctx) error {
	category := new(model.Category)

	if err := c.BodyParser(category); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to POST data",
			"errors": err.Error(),
			"data": nil,
		})
	}
	if err := DB.Create(&category).Error; err != nil {
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
		"data": 1,
	})
}

func UpdateCategory(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	category := new(model.Category)
	if err := DB.Select("id, nama_category").First(&category, id).Error; err != nil{
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
			"errors": err.Error(),
			"data": nil,
		})
	}
	err := DB.Model(&category).Update("nama_category", category.NamaCategory).Error
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
		"data": "",
	})
}

func DeleteCategory(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))
	category := new(model.Category)
	if err := DB.First(&category, id).Error; err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to DELETE data",
			"errors": err.Error(),
			"data": nil,
		})
	}
	if err := DB.Delete(&category, id).Error; err != nil {
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