package main

import (
	"os"

	"github.com/okaaryanata/elearningGolang/handler"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	envFile := os.Getenv("ENV")
	if envFile == "" {
		envFile = ".env"
	}

	err := godotenv.Load(envFile)
	if err != nil {
		panic(err)
	}

	handler.CreateTableBuku()
	handler.CreateTablePeminjaman()
	handler.CreateTableUser()
	e.POST("/insert-data-buku", handler.InsertDataBuku)
	e.GET("/buku/:judul", handler.GetDataBuku)
	e.GET("/get-all-data-buku", handler.GetAllDataBuku)
	e.GET("/siswa/:email", handler.GetDataUser)
	e.POST("/insert-data-peminjaman", handler.InsertDataPeminjaman)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
