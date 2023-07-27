package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"project/model"

	"github.com/gofiber/fiber/v2"
)

func GetAllProvince(c *fiber.Ctx) error {
	resp, err := http.Get("https://www.emsifa.com/api-wilayah-indonesia/api/provinces.json")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  err.Error(),
			"data":    nil,
		})
	}
	provinces := []model.Province{}
	data := string(body)
	err = json.Unmarshal([]byte(data), &provinces)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
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
		"data":    provinces,
	})
}

func GetProvinceById(c *fiber.Ctx) error {
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
	provinces := []model.Province{}
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
	tempProvince := model.Province{}
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
}