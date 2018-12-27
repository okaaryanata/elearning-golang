package main

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/okaaryanata/elearningGolang/handler"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	g := gin.New()

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
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.POST("/insert-data-buku", handler.InsertDataBuku)
	e.GET("/buku/:judul", handler.GetDataBuku)
	e.GET("/get-all-data-buku", handler.GetAllDataBuku)
	e.GET("/user/:emailuser", handler.GetDataUser)
	e.POST("/insert-data-peminjaman", handler.InsertDataPeminjaman)
	g.GET("/logingin", handler.Login)
	e.GET("loginecho", handler.Loginecho)

	r := e.Group("/restricted")
	r.Use(middleware.JWT([]byte("secret")))
	r.GET("", handler.Restricted)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
