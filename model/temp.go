package model

type TempCategory struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	NamaCategory string `gorm:"type:varchar(255)" json:"nama_category"`
}

type TempAlamat struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	JudulAlamat  string `gorm:"type:varchar(255)" json:"judul_alamat"`
	NamaPenerima string `gorm:"type:varchar(255)" json:"nama_penerima"`
	NoTelp       string `gorm:"type:varchar(255)" json:"no_telp" gorm:"unique:not null"`
	DetailAlamat string `gorm:"type:varchar(255)" json:"detail_alamat"`
}

type TempToko struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	NamaToko string `gorm:"type:varchar(255)" json:"nama_toko"`
	UrlFoto  string `gorm:"type:varchar(255)" json:"url_foto"`
}

type Temp_Foto_Produk struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Url      string `gorm:"type:varchar(255)" json:"url"`
	ProdukID uint   `json:"produk_id"`
}

type TempProduk struct {
	ID            uint               `gorm:"primaryKey" json:"id"`
	NamaProduk    string             `gorm:"type:varchar(255)" json:"nama_produk"`
	Slug          string             `gorm:"type:varchar(255)" json:"slug"`
	HargaReseller int                `json:"harga_reseller"`
	HargaKonsumen int                `json:"harga_konsumen"`
	Deskripsi     string             `gorm:"type:text" json:"deskripsi"`
	Toko          TempToko           `json:"toko"`
	Category      TempCategory       `json:"category"`
	Photos        []Temp_Foto_Produk `gorm:"foreignKey:ProdukID" json:"photos"`
}

type TempDetailTrx struct {
	Product TempProduk `json:"product"`
}

type TempTrx struct {
	ID          uint            `gorm:"primaryKey" json:"id"`
	HargaTotal  int             `json:"harga_total"`
	KodeInvoice string          `gorm:"type:varchar(255)" json:"kode_invoice"`
	MethodBayar string          `gorm:"type:varchar(255)" json:"method_bayar"`
	Alamat      TempAlamat      `json:"alamat"`
	Detail_Trx  []TempDetailTrx `json:"detail_trx"`
}