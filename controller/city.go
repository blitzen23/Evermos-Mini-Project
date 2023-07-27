package controller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"project/model"

	"github.com/gofiber/fiber/v2"
)

func GetCityById(c *fiber.Ctx) error {
	city_id := c.Params("city_id")
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
	for _, province := range provinces {
		resp, err := http.Get("https://www.emsifa.com/api-wilayah-indonesia/api/regencies/" + province.ID + ".json")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to GET Province and Kota Data",
				"errors":  err.Error(),
				"data":    nil,
			})
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to GET Province and Kota Data",
				"errors":  err.Error(),
				"data":    nil,
			})
		}
		kotas := []model.Kota{}
		data := string(body)
		err = json.Unmarshal([]byte(data), &kotas)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to GET Province and Kota Data",
				"errors":  err.Error(),
				"data":    nil,
			})
		}
		tempKota := model.Kota{}
		for _, kota := range kotas {
			if kota.ID == city_id {
				tempKota = kota
				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"status":  true,
					"message": "Succeed to GET data",
					"errors":  nil,
					"data":    tempKota,
				})
			}
		}
	}
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"status":  false,
		"message": "Failed to GET data",
		"errors":  "No Data City",
		"data":    nil,
	})
}

func GetALlCitiesInAProvince(c *fiber.Ctx) error {
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
	kotas := []model.Kota{}
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
}