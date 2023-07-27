package controller

import (
	"project/model"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetAllTrx(c *fiber.Ctx) error {
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
	trxes := []model.Trx{}
	err = DB.Preload("Alamat").Preload("Detail_Trx").Where("user_id = ?", user_id).Find(&trxes).Error
	if err != nil && err.Error() != "record not found" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to GET data",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	tempTrx := []model.TempTrx{}
	for _, trx := range trxes {
		tempDetailTrx := []model.TempDetailTrx{}
		for _, detail := range trx.Detail_Trx {
			detailTrx := model.Detail_Trx{}
			err = DB.Preload("Log_Produk").Preload("Toko").First(&detailTrx, detail.ID).Error
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  false,
					"message": "Failed to GET data",
					"errors":  err.Error(),
					"data":    nil,
				})
			}
			logProduk := model.Log_Produk{}
			err = DB.Preload("Category").Preload("Foto_Produks").First(&logProduk, detailTrx.LogProdukID).Error
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  false,
					"message": "Failed to GET data",
					"errors":  err.Error(),
					"data":    nil,
				})
			}
			tempPhotos := []model.Temp_Foto_Produk{}
			for _, fotoProduk := range logProduk.Foto_Produks {
				tempPhotos = append(tempPhotos, model.Temp_Foto_Produk{
					ID:       fotoProduk.ID,
					Url:      fotoProduk.Url,
					ProdukID: fotoProduk.ProdukID,
				})
			}
			tempDetailTrx = append(tempDetailTrx, model.TempDetailTrx{
				Product: model.TempProduk{
					ID:            logProduk.ProdukID,
					NamaProduk:    logProduk.NamaProduk,
					Slug:          logProduk.Slug,
					HargaReseller: logProduk.HargaReseller,
					HargaKonsumen: logProduk.HargaKonsumen,
					Deskripsi:     logProduk.Deskripsi,
					Toko: model.TempToko{
						NamaToko: detailTrx.Toko.NamaToko,
						UrlFoto:  detailTrx.Toko.UrlFoto,
					},
					Category: model.TempCategory{
						ID:           logProduk.Category.ID,
						NamaCategory: logProduk.Category.NamaCategory,
					},
					Photos: tempPhotos,
				},
			})
		}
		tempTrx = append(tempTrx, model.TempTrx{
			ID:          trx.ID,
			HargaTotal:  trx.HargaTotal,
			KodeInvoice: trx.KodeInvoice,
			MethodBayar: trx.MethodBayar,
			Alamat: model.TempAlamat{
				ID:           trx.Alamat.ID,
				NamaPenerima: trx.Alamat.NamaPenerima,
				JudulAlamat:  trx.Alamat.JudulAlamat,
				NoTelp:       trx.Alamat.NoTelp,
				DetailAlamat: trx.Alamat.DetailAlamat,
			},
			Detail_Trx: tempDetailTrx,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to GET data",
		"errors":  nil,
		"data": fiber.Map{
			"data": tempTrx,
		},
	})
}

func GetTrxById(c *fiber.Ctx) error {
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
	trx_id, _ := strconv.Atoi(c.Params("id"))
	trx := model.Trx{}
	err = DB.Preload("Alamat").Preload("Detail_Trx").Where("id = ?", trx_id).Find(&trx).Error
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
		logProduk := model.Log_Produk{}
		err = DB.Preload("Category").Preload("Foto_Produks").First(&logProduk, detailTrx.LogProdukID).Error
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
}

func CreateTrx(c *fiber.Ctx) error {
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
	alamat := model.Alamat{}
	err = DB.Where("user_id = ?", user_id).First(&alamat, tempTrx.AlamatKirim).Error
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
	trx := model.Trx{}
	trx.AlamatID = uint(tempTrx.AlamatKirim)
	trx.MethodBayar = tempTrx.MethodBayar
	trx.UserID = user_id.(uint)
	time := strconv.FormatInt(time.Now().Unix(), 10)
	trx.KodeInvoice = "INV-" + time
	trx.HargaTotal = 0
	err = DB.Create(&trx).Error
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
		produk := model.Produk{}
		err = DB.First(&produk, detailTrx.ProductID).Error
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
		logProduk := model.Log_Produk{}
		err = DB.Where("produk_id = ?", produk.ID).Find(&logProduk).Error
		if err != nil && err.Error() == "record not found"{
			logProduk.ProdukID = produk.ID
			logProduk.NamaProduk = produk.NamaProduk
			logProduk.Slug = produk.Slug
			logProduk.HargaReseller = produk.HargaReseller
			logProduk.HargaKonsumen = produk.HargaKonsumen
			logProduk.Deskripsi = produk.Deskripsi
			logProduk.TokoID = produk.TokoID
			logProduk.CategoryID = produk.CategoryID
			err = DB.Create(&logProduk).Error
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
		detail := model.Detail_Trx{}
		detail.TrxID = trx.ID
		detail.LogProdukID = logProduk.ID
		detail.TokoID = produk.TokoID
		detail.Kuantitas = detailTrx.Kuantitas
		detail.HargaTotal = produk.HargaKonsumen
		err = DB.Create(&detail).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": false,
				"message": "Failed to POST data",
				"errors": err.Error(),
				"data": nil,
			})
		}
		err = DB.Model(&produk).Update("stok", produk.Stok - detail.Kuantitas).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": false,
				"message": "Failed to POST data",
				"errors": err.Error(),
				"data": nil,
			})
		}
	}
	err = DB.Model(&trx).Update("harga_total", hargaTotal).Error
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
}