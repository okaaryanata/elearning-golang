package handler

import (
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //ni utk db
)

//connect ke DB
// db, err := gorm.Open("postgres", "host=myhost port=myport user=gorm dbname=gorm password=mypassword")
func getConnection() (db *gorm.DB, err error) {

	db, err = gorm.Open("postgres", os.Getenv("DBCONN"))
	if err != nil {

		return
	}

	return
}
