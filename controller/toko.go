package controller

import (
	"fmt"
	"project/model"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetMyToko(c *fiber.Ctx) error {
	toko := new(model.Toko)
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
	err = DB.Where("user_id = ?", user_id).Find(&toko).Error
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
			"id":        toko.ID,
			"nama_toko": toko.NamaToko,
			"url_foto":  toko.UrlFoto,
			"user_id":   toko.UserID,
		},
	})
}

func UpdateToko(c *fiber.Ctx) error {
	toko := new(model.Toko)
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
	toko_id, err := strconv.Atoi(c.Params("id_toko"))
	if (err != nil){
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to PUT data",
			"errors": err.Error(),
			"data": nil,
		})
	}
	err = DB.First(&toko, toko_id).Error
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
		err = DB.Model(&toko).Update(&toko).Error
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
}

func GetTokoById(c *fiber.Ctx) error {
	toko := new(model.Toko)
	toko_id, err := strconv.Atoi(c.Params("id_toko"))
	if (err != nil){
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": err.Error(),
			"data": nil,
		})
	}
	err = DB.First(&toko, toko_id).Error
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
}

func GetAllToko(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 1)

	offset := (page - 1) * limit

	tokos := []model.Toko{}
	err := DB.Limit(limit).Offset(offset).Find(&tokos).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": err.Error(),
			"data": nil,
		})
	}
	tempTokos := []model.TempToko{}
	for _, toko := range tokos{
		tempTokos = append(tempTokos, model.TempToko{
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
}

func GetAllProduct(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	page := c.QueryInt("page", 1)

	offset := (page - 1) * limit

	produks := []model.Produk{}
	err := DB.Limit(limit).Offset(offset).Preload("Toko").Preload("Category").Preload("Foto_Produks").Find(&produks).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": err.Error(),
			"data": nil,
		})
	}

	type TempProduk struct{
		ID uint `gorm:"primaryKey" json:"id"`
		NamaProduk string `gorm:"type:varchar(255)" json:"nama_produk"`
		Slug string `gorm:"type:varchar(255)" json:"slug"`
		HargaReseller int `json:"harga_reseller"`
		HargaKonsumen int `json:"harga_konsumen"`
		Stok int `json:"stok"`
		Deskripsi string `gorm:"type:text" json:"deskripsi"`
		Toko model.TempToko `json:"toko"`
		Category model.TempCategory `json:"category"`
		Photos []model.Temp_Foto_Produk `gorm:"foreignKey:ProdukID" json:"foto_produks"`
	}
	tempProduks := []TempProduk{}
	for _, produk := range produks{
		tempFotoProduks := []model.Temp_Foto_Produk{}
		for _, foto_produk := range produk.Foto_Produks{
			tempFotoProduks = append(tempFotoProduks, model.Temp_Foto_Produk{
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
			Toko: model.TempToko{
				ID: produk.TokoID,
				NamaToko: produk.Toko.NamaToko,
				UrlFoto: produk.Toko.UrlFoto,
			},
			Photos: tempFotoProduks,
			Category: model.TempCategory{
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
}

func GetProductById(c *fiber.Ctx) error {
	produk := model.Produk{}
	id, _ := strconv.Atoi(c.Params("id"))
	if err := DB.Preload("Toko").Preload("Category").Preload("Foto_Produks").First(&produk, id).Error; err != nil {
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
	tempFotoProduks := []model.Temp_Foto_Produk{}
	for _, foto_produk := range produk.Foto_Produks{
		tempFotoProduks = append(tempFotoProduks, model.Temp_Foto_Produk{
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
}

func CreateProduct(c *fiber.Ctx) error {
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
	toko := model.Toko{}
	err = DB.Preload("User").Where("user_id = ?", user_id).Find(&toko).Error
	if (err != nil){
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to POST data",
			"errors": err.Error(),
			"data": nil,
		})
	}
	if (user_id == toko.UserID){
		produk := model.Produk{}
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
			category := model.Category{}
			err = DB.First(&category, produk.CategoryID).Error
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
			err = DB.Create(&produk).Error
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
					fotoProduk := model.Foto_Produk{}
					time := strconv.FormatInt(time.Now().Unix(), 10)
					if err := c.SaveFile(file, fmt.Sprintf("./images/%s-%s", time, file.Filename)); err != nil {
						return err
					}
					fotoProduk.ProdukID = produk.ID
					fotoProduk.Url = "/images/" + time + "-" + file.Filename
					if err := DB.Create(&fotoProduk).Error; err != nil {
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
}

func UpdateProduct(c *fiber.Ctx) error {
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
	toko := model.Toko{}
	err = DB.Preload("User").Where("user_id = ?", user_id).Find(&toko).Error
	if (err != nil){
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to PUT data",
			"errors": err.Error(),
			"data": nil,
		})
	}
	if (user_id == toko.UserID){
		produk := model.Produk{}
		produk_id, _ := strconv.Atoi(c.Params("id"))
		err = DB.First(&produk, produk_id).Error
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
				category := model.Category{}
				err = DB.First(&category, produk.CategoryID).Error
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
				err = DB.Model(&produk).Update(produk).Error
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
						fotoProduk := model.Foto_Produk{}
						time := strconv.FormatInt(time.Now().Unix(), 10)
						if err := c.SaveFile(file, fmt.Sprintf("./images/%s-%s", time, file.Filename)); err != nil {
							return err
						}
						fotoProduk.ProdukID = produk.ID
						fotoProduk.Url = "/images/" + time + "-" + file.Filename
						if err := DB.Create(&fotoProduk).Error; err != nil {
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
}

func DeleteProduct(c *fiber.Ctx) error {
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
	toko := model.Toko{}
	err = DB.Preload("User").Where("user_id = ?", user_id).Find(&toko).Error
	if (err != nil){
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to DELETE data",
			"errors": err.Error(),
			"data": nil,
		})
	}
	if (user_id == toko.UserID){
		produk := model.Produk{}
		produk_id, _ := strconv.Atoi(c.Params("id"))
		err = DB.First(&produk, produk_id).Error
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
			err = DB.Delete(&produk, produk.ID).Error
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
}