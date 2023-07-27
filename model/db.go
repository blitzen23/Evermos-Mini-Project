package model

import (
	"time"
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
