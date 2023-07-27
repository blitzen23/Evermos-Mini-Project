package routing

import (
	"project/controller"

	"github.com/gofiber/fiber/v2"
)

var App *fiber.App

func Start() {
	controller.Start()
	App = fiber.New()
	App.Post("/auth/register", controller.CreateUser)
	App.Post("/auth/login", controller.Login)

	App.Get("/category", controller.GetAllCategory)
	App.Get("/category/:id", controller.IsAdmin, controller.GetCategoryById)
	App.Post("/category", controller.IsAdmin, controller.CreateCategory)
	App.Put("/category/:id", controller.IsAdmin, controller.UpdateCategory)
	App.Delete("/category/:id", controller.IsAdmin, controller.DeleteCategory)

	App.Get("/user/alamat", controller.IsLoggedIn, controller.GetMyAlamat)
	App.Get("/user/alamat/:id", controller.IsLoggedIn, controller.GetAlamatById)
	App.Post("/user/alamat", controller.IsLoggedIn, controller.CreateAlamat)
	App.Put("/user/alamat/:id", controller.IsLoggedIn, controller.UpdateAlamat)
	App.Delete("/user/alamat/:id", controller.IsLoggedIn, controller.DeleteAlamat)
	
	App.Get("/provcity/listcities/:prov_id", controller.GetALlCitiesInAProvince)
	App.Get("/provcity/detailcity/:city_id", controller.GetCityById)

	App.Get("/product", controller.GetAllProduct)
	App.Get("/product/:id", controller.GetProductById)
	App.Post("/product", controller.IsLoggedIn, controller.CreateProduct)
	App.Put("/product/:id", controller.IsLoggedIn, controller.UpdateProduct)
	App.Delete("/product/:id", controller.IsLoggedIn, controller.DeleteProduct)

	App.Get("/provcity/listprovincies", controller.GetAllProvince)
	App.Get("/provcity/detailprovince/:prov_id", controller.GetProvinceById)

	App.Get("/toko/my", controller.IsLoggedIn, controller.GetMyToko)
	App.Put("/toko/:id_toko", controller.IsLoggedIn, controller.UpdateToko)
	App.Get("/toko/:id_toko", controller.IsLoggedIn, controller.GetTokoById)
	App.Get("/toko", controller.IsLoggedIn, controller.GetAllToko)

	App.Get("/trx", controller.IsLoggedIn, controller.GetAllTrx)
	App.Get("/trx/:id", controller.IsLoggedIn, controller.GetTrxById)
	App.Post("/trx", controller.IsLoggedIn, controller.CreateTrx)

	App.Get("/user", controller.IsLoggedIn, controller.GetUser)
	App.Put("/user", controller.IsLoggedIn, controller.UpdateUser)

	App.Listen(":3000")
}