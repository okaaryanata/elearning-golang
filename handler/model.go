package handler

import (
	"time"

	"github.com/jinzhu/gorm"
)

// User - tabel menyimpan data user
type User struct {
	Id       uint   `gorm:"primary_key"`
	Nama     string `gorm:"varchar()"`
	Email    string `gorm:"varchar()"`
	Password string `gorm:"varchar()"`
}

// Buku - tabel menyimpan data buku
type Buku struct {
	Id          uint   `gorm:"primary_key"`
	Judul       string `gorm:"varchar()"`
	Tahunterbit int
	Pengarang   string `gorm:"varchar()"`
}

// Peminjaman - tabel menyimpan data peminjaman
type Peminjaman struct {
	Id             uint `gorm:"primary_key"`
	Users          User `gorm:"foreignkey:IDUser;association_foreignkey:ID"`
	IDUser         uint
	Bukus          Buku `gorm:"foreignkey:IDBuku;association_foreignkey:ID"`
	IDBuku         uint
	Tanggalpinjam  time.Time
	Tanggalkembali time.Time
	Ispinjam       bool
	Iskembali      bool
}

// CreateTableUser buat table User
func (t User) CreateTableUser(db *gorm.DB) (err error) {
	// auto migrate buat kalo tabelnya gaada, auto create. kalo structnya berubah, yaudah tablenya ngapdet
	err = db.AutoMigrate(&User{}).Error

	return
}

// CreateTableBuku buat table Buku
func (t Buku) CreateTableBuku(db *gorm.DB) (err error) {
	// auto migrate buat kalo tabelnya gaada, auto create. kalo structnya berubah, yaudah tablenya ngapdet
	err = db.AutoMigrate(&Buku{}).Error

	return
}

// CreateTablePeminjaman buat table Siswa
func (t Peminjaman) CreateTablePeminjaman(db *gorm.DB) (err error) {
	// auto migrate buat kalo tabelnya gaada, auto create. kalo structnya berubah, yaudah tablenya ngapdet
	err = db.AutoMigrate(&Peminjaman{}).Error

	return
}

// GetDataUserByEmail mendapatkan data user berdasarkan email
func (t User) GetDataUserByEmail(db *gorm.DB) (res User, err error) {

	err = db.Where("email = ?", t.Email).First(&res).Error

	return
}

// InsertDataBuku masukan data ke table Buku
func (t Buku) InsertDataBuku(db *gorm.DB) (err error) {
	err = db.Create(&t).Error

	return
}

// GetAllDataBuku menampilkan semua data Buku
func (t Buku) GetAllDataBuku(db *gorm.DB) (res []Buku, err error) {
	err = db.Find(&res).Error

	return
}

// GetDataBukuByJudul mendapatkan data uuku berdasarkan ID Buku
func (t Buku) GetDataBukuByJudul(db *gorm.DB) (res Buku, err error) {

	err = db.Where("judul = ?", t.Judul).First(&res).Error

	return
}

// InsertDataPeminjaman masukan data ke table peminjaman
func (t Peminjaman) InsertDataPeminjaman(db *gorm.DB) (err error) {
	err = db.Create(&t).Error
	if err != nil {
		return err
	} else {
		err = db.Model(&t).Updates(Peminjaman{Tanggalpinjam: time.Now(), Ispinjam: true}).Error
	}

	return
}

// UpdateDataPeminjaman update data ke tabel peminjaman
func (t Peminjaman) UpdateDataPeminjaman(db *gorm.DB) (err error) {
	err = db.Model(&t).Updates(Peminjaman{Tanggalkembali: time.Now(), Iskembali: true}).Error

	return
}

// CekEmailPass cek email dan pass saat login
func (t User) CekEmailPass(db *gorm.DB) (res User, err error) {
	err = db.Where("email = ? AND password = ?", t.Email, t.Password).First(&res).Error

	return
}
