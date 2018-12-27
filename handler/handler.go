package handler

import (
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo"
)

type buku struct {
	Judul       string `json:"judul"`
	Tahunterbit int    `json:"tahunterbit"`
	Pengarang   string `json:"pengarang"`
}

type credential struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type pinjam struct {
	Email string `json:"email"`
	Judul string `json:"judul"`
}

type riwayatBuku struct {
	Judul       string
	Tahunterbit int
	Pengarang   string
	riwayat     []peminjaman
}

type peminjaman struct {
	Nama           string
	Tanggalpinjam  string
	Tanggalkembali string
}

// InsertDataBuku buat data buku
func InsertDataBuku(c echo.Context) error {
	inputData := new(buku)

	if err := c.Bind(inputData); err != nil {
		return echo.NewHTTPError(400)

	}

	newBuku := Buku{
		Judul:       inputData.Judul,
		Tahunterbit: inputData.Tahunterbit,
		Pengarang:   inputData.Pengarang,
	}

	db, err := getConnection()
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}
	defer db.Close()

	err = newBuku.InsertDataBuku(db)
	if err != nil {
		return echo.NewHTTPError(500, err.Error())
	}

	return c.JSON(201, inputData)
}

// GetDataBuku - ini untuk get data buku by judul buku
func GetDataBuku(c echo.Context) error {
	// get input user dari query param dengan key judul buku
	judul := c.Param("judul")

	data := Buku{
		Judul: judul,
	}

	// get koneksi ke DB untuk di passing saat ingin query
	db, err := getConnection()
	if err != nil {
		return echo.NewHTTPError(500, err.Error)
	}
	defer db.Close()

	response, err := data.GetDataBukuByJudul(db)
	if err != nil {
		return echo.NewHTTPError(500, err.Error)
	}

	// riwayatPinjam

	rBuku := riwayatBuku{
		Judul:       response.Judul,
		Tahunterbit: response.Tahunterbit,
		Pengarang:   response.Pengarang,
		// riwayat :
	}

	return c.JSON(200, rBuku)
}

// InsertDataPeminjaman tambah data peminjaman ke tabel peminjaman
func InsertDataPeminjaman(c echo.Context) error {
	inputData := new(pinjam)

	if err := c.Bind(inputData); err != nil {
		return echo.NewHTTPError(400)

	}

	judulBuku := Buku{
		Judul: inputData.Judul,
	}

	emailUser := User{
		Email: inputData.Email,
	}

	// get koneksi ke DB untuk di passing saat ingin query
	db, err := getConnection()
	if err != nil {
		return echo.NewHTTPError(500, err.Error)
	}
	defer db.Close()

	responseBuku, err := judulBuku.GetDataBukuByJudul(db)
	if err != nil {
		return echo.NewHTTPError(500, err.Error)
	}

	responseUser, err := emailUser.GetDataUserByEmail(db)
	if err != nil {
		return echo.NewHTTPError(500, err.Error)
	}

	pinjam := Peminjaman{
		IDUser: responseUser.Id,
		IDBuku: responseBuku.Id,
	}

	err = db.Find(&pinjam, "id_user = ? AND ispinjam = ? AND id_buku = ?", responseUser.Id, true, responseBuku.Id).Error
	if err != nil {
		err = pinjam.InsertDataPeminjaman(db)
		if err != nil {
			return echo.NewHTTPError(500, err.Error())
		}
	} else if err = db.Find(&pinjam, "id_user = ? AND ispinjam = ? AND id_buku = ? AND iskembali = ?", responseUser.Id, true, responseBuku.Id, true).Error; err == nil {
		return c.JSON(201, "tidak dapat pinjam kembali")
	} else {
		err = pinjam.UpdateDataPeminjaman(db)
		if err != nil {
			return echo.NewHTTPError(500, err.Error())
		}
	}

	return c.JSON(201, pinjam)
}

// GetAllDataBuku - ini untuk get all data Buku
func GetAllDataBuku(c echo.Context) error {
	var d Buku

	// get koneksi ke DB untuk di passing saat ingin query
	db, err := getConnection()
	if err != nil {
		return echo.NewHTTPError(500, err.Error)
	}
	defer db.Close()

	response, err := d.GetAllDataBuku(db)
	if err != nil {
		return echo.NewHTTPError(500, err.Error)
	}

	return c.JSON(200, response)
}

// Login - fungsi untuk login
func Login(c *gin.Context) {
	var user credential
	err := c.Bind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "can't bind struct",
		})
	}

	loginData := User{
		Email:    user.Email,
		Password: user.Password,
	}

	db, err := getConnection()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": err,
		})
	}

	defer db.Close()

	res, err := loginData.CekEmailPass(db)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  http.StatusUnauthorized,
			"message": err,
		})
	}

	sign := jwt.New(jwt.GetSigningMethod("HS256"))
	sign.Claims = jwt.MapClaims{
		"id": res.Id,
	}
	token, err := sign.SignedString([]byte("secret"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		c.Abort()
	}
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// Loginecho - fungsi untuk login
func Loginecho(c echo.Context) error {
	var user credential
	err := c.Bind(&user)
	if err != nil {
		return echo.NewHTTPError(500, err.Error)
	}

	loginData := User{
		Email:    user.Email,
		Password: user.Password,
	}

	db, err := getConnection()
	if err != nil {
		return echo.NewHTTPError(500, err.Error)
	}

	defer db.Close()

	res, err := loginData.CekEmailPass(db)
	if err != nil {
		return echo.NewHTTPError(500, err.Error)
	} else {
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["id"] = res.Id
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}

	return echo.ErrUnauthorized
}

// Restricted func
func Restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id_user := claims["id"].(string)
	return c.String(http.StatusOK, "Welcome "+id_user+"!")
}

// Auth - fungsi untuk auth
func Auth(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")
	// return c.JSON(200, tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secret"), nil
	})

	if token != nil && err == nil {
		fmt.Println("token verified")
	} else {
		result := gin.H{
			"message": "not authorized",
			"error":   err.Error(),
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
	}
}

// GetDataUser - ini untuk get data user by email
func GetDataUser(c echo.Context) error {
	// get input user dari query param dengan key name
	email := c.Param("emailuser")

	// inisisasi data dengan type TableAdindaOka
	data := User{
		Email: email,
	}

	// get koneksi ke DB untuk di passing saat ingin query
	db, err := getConnection()
	if err != nil {
		return echo.NewHTTPError(500, err.Error)
	}
	defer db.Close()

	response, err := data.GetDataUserByEmail(db)
	if err != nil {
		return echo.NewHTTPError(500, err.Error)
	}

	return c.JSON(200, response)
}

// CreateTableBuku buat table siswa
func CreateTableBuku() (err error) {
	db, err := getConnection()
	if err != nil {
		return err
	}

	table := new(Buku)

	err = table.CreateTableBuku(db)
	if err != nil {
		panic(err)
	}
	return err
}

// CreateTableUser buat table siswa
func CreateTableUser() (err error) {
	db, err := getConnection()
	if err != nil {
		return err
	}

	table := new(User)

	err = table.CreateTableUser(db)
	if err != nil {
		panic(err)
	}
	return err
}

// CreateTablePeminjaman buat table siswa
func CreateTablePeminjaman() (err error) {
	db, err := getConnection()
	if err != nil {
		return err
	}

	table := new(Peminjaman)

	err = table.CreateTablePeminjaman(db)
	if err != nil {
		panic(err)
	}
	return err
}
